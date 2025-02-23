#!/usr/bin/env bash

# A secure store for CI artifacts

SCRIPTS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd)"
source "$SCRIPTS_ROOT/scripts/ci/gcp.sh"

set -euo pipefail

store_artifacts() {
    info "Storing artifacts"

    if [[ "$#" -lt 1 ]]; then
        die "missing args. usage: store_artifacts <path> [<destination>]"
    fi

    local path="$1"
    local destination="${2:-$(basename "$path")}"

    if [ -z "${path:-}" ]; then
        echo ERROR: Missing required path parameter
        exit 1
    fi

    # Circle CI does a poor job with ~ expansion
    if [[ "$path" =~ ^~ ]]; then
        path="$HOME$(cut -c2- -<<< "$path")"
    fi

    if [[ ! -e "$path" ]]; then
        echo INFO: "$path" is missing, nothing to upload
        exit 0
    fi

    if [[ -d "$path" ]] && [[ -z "$(ls -A "$path")" ]]; then
        # skip empty dirs because gsutil considers this an error
        echo INFO: "$path" is empty, nothing to upload
        exit 0
    fi

    _artifacts_preamble

    local gs_destination
    gs_destination=$(get_unique_gs_destination "${destination}")

    info "Writing to $gs_destination"
    gsutil -m cp -r "$path" "$gs_destination"
}

_artifacts_preamble() {
    ensure_CI
    require_executable "gsutil"

    setup_gcp
    gsutil version -l

    set_gs_path_vars
}

get_unique_gs_destination() {
    local desired_destination="$1"
    local index=1
    local destination="$GS_JOB_URL/${desired_destination}"
    while gsutil ls "$destination" > /dev/null 2>&1; do
        (( index++ ))
        destination="$GS_JOB_URL/${desired_destination}-$index"
        if [[ $index -gt 50 ]]; then
            echo ERROR: too many attempts to find a unique destination suffix
            exit 1
        fi
    done
    echo "${destination}"
}

set_gs_path_vars() {
    GS_URL="gs://roxci-artifacts"

    if is_OPENSHIFT_CI; then
        require_environment "REPO_NAME"
        require_environment "BUILD_ID"
        require_environment "JOB_NAME"
        if [ -z "${PULL_PULL_SHA:-}" ] && [ -z "${PULL_BASE_SHA:-}" ]; then
            die "There is no ID suitable to separate artifacts for this commit"
        fi
        local workflow_id="${PULL_PULL_SHA:-${PULL_BASE_SHA}}"
        WORKFLOW_SUBDIR="${REPO_NAME}/${workflow_id}"
        JOB_SUBDIR="${BUILD_ID}-${JOB_NAME}"
        GS_JOB_URL="${GS_URL}/${WORKFLOW_SUBDIR}/${JOB_SUBDIR}"
    elif is_CIRCLECI; then
        require_environment "CIRCLE_PROJECT_REPONAME"
        require_environment "CIRCLE_WORKFLOW_ID"
        require_environment "CIRCLE_BUILD_NUM"
        require_environment "CIRCLE_JOB"

        WORKFLOW_SUBDIR="${CIRCLE_PROJECT_REPONAME}/${CIRCLE_WORKFLOW_ID}"
        JOB_SUBDIR="${CIRCLE_BUILD_NUM}-${CIRCLE_JOB}"
        GS_JOB_URL="${GS_URL}/${WORKFLOW_SUBDIR}/${JOB_SUBDIR}"
    else
        die "Support is missing for this CI environment"
    fi
}

fixup_artifacts_content_type() {

    _artifacts_preamble

    local fixups=(
        "*.log:text/plain"
    )

    for fixup in "${fixups[@]}"; do
        IFS=':' read -ra parts <<< "$fixup"
        local file_match="${parts[0]}"
        local content_type="${parts[1]}"

        gsutil -m setmeta -h "Content-Type:$content_type" "${GS_JOB_URL}/**/$file_match" || true
    done
}

make_artifacts_help() {

    _artifacts_preamble
    
    local gs_workflow_url="$GS_URL/$WORKFLOW_SUBDIR"
    local gs_job_url="$gs_workflow_url/$JOB_SUBDIR"
    local browser_url="https://console.cloud.google.com/storage/browser/roxci-artifacts"
    local browser_job_url="$browser_url/$WORKFLOW_SUBDIR/$JOB_SUBDIR"

    local help_file
    if is_OPENSHIFT_CI; then
        require_environment "ARTIFACT_DIR"
        help_file="$ARTIFACT_DIR/howto-locate-artifacts.html"
    elif is_CIRCLECI; then
        help_file="/tmp/howto-locate-artifacts.html"
    else
        die "This is an unsupported environment"
    fi

    cat > "$help_file" <<- EOH
        Artifacts are stored in a GCS bucket ($GS_URL). There are at least two options for access:

        <h3>gsutil cp</h3>

        Copy all artifacts for the build/job:
        <pre>gsutil -m cp -r $gs_job_url .</pre>

        Copy all artifacts for the entire workflow:
        <pre>gsutil -m cp -r $gs_workflow_url .</pre>

        <h3>Browse using the google cloud UI</h3>

        <p>The URL you use will depend on the <i>authuser</i> value you use for your @stackrox.com account.</p>

        <a href="$browser_job_url?authuser=0">authuser=0</a><br>
        <a href="$browser_job_url?authuser=1">authuser=1</a><br>
        <a href="$browser_job_url?authuser=2">authuser=2</a><br>
EOH

    info "Artifacts are stored in a GCS bucket ($GS_URL)"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    if [[ "$#" -lt 1 ]]; then
        die "When invoked at the command line a method is required."
    fi
    fn="$1"
    shift
    "$fn" "$@"
fi
