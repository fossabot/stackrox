#!/usr/bin/env bash

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
source "$ROOT/scripts/ci/lib.sh"

set -euo pipefail

shopt -s nullglob
for cred in /tmp/secret/**/[A-Z]*; do
    export "$(basename "$cred")"="$(cat "$cred")"
done

openshift_ci_mods

function hold() {
    while [[ -e /tmp/hold ]]; do
        info "Holding this job for debug"
        sleep 60
    done
}
trap hold EXIT

if [[ "$#" -lt 1 ]]; then
    die "usage: dispatch <ci-job> [<...other parameters...>]"
fi

ci_job="$1"
shift

gate_job "$ci_job"

case "$ci_job" in
    style-checks)
        make style
        ;;
    push-images)
        "$ROOT/scripts/ci/jobs/push-images.sh"
        ;;
    gke-upgrade-tests)
        "$ROOT/.openshift-ci/gke_upgrade_test.py"
        ;;
    *)
        # For ease of initial integration this function does not fail when the
        # job is unknown.
        info "nothing to see here: ${ci_job}"
        exit 0
esac
