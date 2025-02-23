# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
# The version reported by the root Makefile is converted to be compatible with SemVer. Specifically, inner-zero is
# dropped (e.g. 3.0.61.1 -> 3.61.1) and development version ".x" is changed to ".0" (e.g. 3.0.61.x-123 -> 3.0.61.0-123).
VERSION ?= $(shell $(MAKE) --quiet --no-print-directory -C .. tag | sed -E 's@^(([[:digit:]]+\.)+)x(-)?@\10\3@g' | sed -E 's@^3.0.([[:digit:]]+\.[[:digit:]]+)(-)?@3.\1\2@g')

GO_VERSION ?= $(shell cat ../EXPECTED_GO_VERSION | sed 's/go//')


# CHANNELS define the bundle channels used in the bundle.
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "preview,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=preview,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="preview,fast,stable")
CHANNELS ?= latest
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle.
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
DEFAULT_CHANNEL ?= latest
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# IMAGE_REPO is the repository (server and namespace) into which the operator images will get pushed.
IMAGE_REPO ?= docker.io/stackrox
# IMAGE_TAG_BASE defines the docker.io namespace and part of the image name for remote images.
# This variable is used to construct full image tags for bundle and catalog images.
IMAGE_TAG_BASE ?= $(IMAGE_REPO)/stackrox-operator

# BUNDLE_IMG defines the image:tag used for the bundle.
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)
BUNDLE_IMG ?= $(IMAGE_TAG_BASE)-bundle:v$(VERSION)

# INDEX_IMG defines the image:tag used for the index (a.k.a. catalog) image, containing the single operator bundle image.
INDEX_IMG_BASE = $(IMAGE_TAG_BASE)-index
INDEX_IMG ?= $(INDEX_IMG_BASE):v$(VERSION)

# Image URL to use all building/pushing image targets
IMG ?= $(IMAGE_TAG_BASE):$(VERSION)

# Set up the so-called "storage image" for executing scorecard test suite.
# By default a busybox image is pulled from docker.io. In order to rule out timing issues related to rate limiting
# we use our own copy of a busybox image from quay.io.
#
# We expect this dependency on busybox to not be very volatile, but in case it is necessary to use a newer busybox
# image for the scorecard tests, proceed as follows:
#
# 1. Push a new busybox image to the quay.io repository.
# 2. Update the image image tag below.
SCORECARD_STORAGE_IMAGE_BASE ?= quay.io/rhacs-eng/mirror-busybox
SCORECARD_STORAGE_IMAGE_TAG ?= @sha256:fe9664ceda9c13b7fd65f5bc3c0bca5060de3a66fedc39d33926b573599da3b7
SCORECARD_STORAGE_IMAGE ?= $(SCORECARD_STORAGE_IMAGE_BASE):$(SCORECARD_STORAGE_IMAGE_TAG)

# Increase the scorecard wait time from (default) 30s to 5m, since we experienced some time out issues in the past.
SCORECARD_WAIT_TIME_DURATION ?= 5m

# SCORECARD_ARGS should be passed to any invocation of "$(OPERATOR_SDK) scorecard".
SCORECARD_ARGS ?= --storage-image="$(SCORECARD_STORAGE_IMAGE)" --wait-time="$(SCORECARD_WAIT_TIME_DURATION)" --verbose

# TEST_NAMESPACE is where the operator is installed for e2e tests by CI.
TEST_NAMESPACE ?= stackrox-operator

# ENABLE_WEBHOOKS determines if webhooks should be set up in manager when running locally.
# Disabled by default since it requires setting up of custom routing from your k8s API server
# to your manager process.
ENABLE_WEBHOOKS ?= false

# ROX_IMAGE_FLAVOR is an ARG used in Dockerfiles that defines the default registries for main, scanner, and collector images.
# ROX_IMAGE_FLAVOR valid values are: development_build, stackrox.io, rhacs.
# Use the value figured in the parent directory Makefile, unless provided explicitly in the environment.
ROX_IMAGE_FLAVOR ?= $(shell $(MAKE) --quiet --no-print-directory -C .. image-flavor)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Lowercase Operating System name, needed for downloading GitHub releases.
OS=$(shell uname | tr A-Z a-z)

PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash
.SHELLFLAGS = -o pipefail -euc

# This Makefile is not created for and likely would not benefit from concurrent jobs execution.
.NOTPARALLEL:

all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

tag: ## Print the correct operator version (== image tag).
	@echo $(VERSION)

image-flavor: ## Print the current image flavor.
	@echo $(ROX_IMAGE_FLAVOR)

# `everything` and `all` are separate targets.
# `everything` pushes images. `all` is the default target that's executed when you simply run `make` without arguments.
# If `all` would push images, parameter-less `make` execution would have side-effect of pushing images to registry
# which might not be what the user intended to do (even though it could be safe in many cases).
# `everything` is a target that needs to be explicitly called for this reason.
.PHONY: everything
everything: build bundle ## Build everything (local binary, operator image, bundle image, index image) and push all images.
	$(MAKE) docker-build bundle-build
	$(MAKE)	docker-push docker-push-bundle
	$(MAKE) index-build
	$(MAKE) docker-push-index

##@ Dependencies download

CONTROLLER_GEN = $(PROJECT_DIR)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

KUSTOMIZE = $(PROJECT_DIR)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

OPERATOR_SDK_VERSION = 1.14.0
OPERATOR_SDK = $(PROJECT_DIR)/bin/operator-sdk-$(OPERATOR_SDK_VERSION)
.PHONY: operator-sdk
operator-sdk: ## Download operator-sdk necessary for scaffolding and bundling.
# See https://sdk.operatorframework.io/docs/installation/#install-from-github-release
	@ARCH=$$(case $$(uname -m) in (x86_64) echo -n amd64 ;; (aarch64) echo -n arm64 ;; (*) echo -n $$(uname -m) ;; esac) ;\
	OPERATOR_SDK_URL=https://github.com/operator-framework/operator-sdk/releases/download/v$(OPERATOR_SDK_VERSION)/operator-sdk_$(OS)_$${ARCH} ;\
	$(GET_GITHUB_RELEASE_FN); \
	get_github_release $(OPERATOR_SDK) $${OPERATOR_SDK_URL}

KUTTL_VERSION = 0.11.0
KUTTL_UPSTREAM = kudobuilder
KUTTL ?= $(PROJECT_DIR)/bin/kubectl-kuttl-$(KUTTL_VERSION)
.PHONY: kuttl
kuttl: ## Download kuttl.
	@ARCH=$$(case $$(uname -m) in (aarch64) echo -n armv6 ;; (*) echo -n $$(uname -m) ;; esac) ;\
	KUTTL_URL=https://github.com/$(KUTTL_UPSTREAM)/kuttl/releases/download/v$(KUTTL_VERSION)/kubectl-kuttl_$(KUTTL_VERSION)_$(OS)_$${ARCH} ;\
	$(GET_GITHUB_RELEASE_FN); \
	get_github_release $(KUTTL) $${KUTTL_URL}

ENVTEST = $(PROJECT_DIR)/bin/setup-envtest
envtest: ## Download envtest-setup locally if necessary.
	$(call go-get-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

GET_GITHUB_RELEASE_FN = get_github_release() { \
	[ -f $${1} ] || { \
		set -euxo pipefail ;\
		mkdir -p bin ;\
		curl --fail --location --output $${1} $${2} ;\
		chmod +x $${1} ;\
		[[ "$$(uname -s)" != "Darwin" ]] || xattr -c $${1} ;\
	} \
}

# go-get-tool will 'go get -d' any package $2 and install it to $(PROJECT_DIR)/bin unless $1 already exists.
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
go get -d $(2) ;\
GOBIN=$(PROJECT_DIR)/bin go install $(firstword $(subst @, ,$(2))) ;\
rm -rf $$TMP_DIR ;\
}
endef

##@ Development

parent-proto-generate: ## Make sure ../generated directory has up-to-date content that this operator (transitively) depends upon.
	$(MAKE) -C .. proto-generated-srcs

manifests: parent-proto-generate controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: parent-proto-generate controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	# The generated source files might not comply with the current go formatting, so format them explicitly.
	go fmt ./apis/...

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.22

test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./... -coverprofile cover.out

.PHONY: validate-crs
CRDS := centrals.platform.stackrox.io securedclusters.platform.stackrox.io
validate-crs:
	./tests/scripts/validate-crs.sh ./tests $(CRDS)

test-e2e: build install validate-crs kuttl ensure-rox-main-image-exists ## Run e2e tests with local manager.
	mkdir -p $(PROJECT_DIR)/build/kuttl-test-artifacts
	ENABLE_WEBHOOKS=$(ENABLE_WEBHOOKS) $(KUTTL) test

test-e2e-deployed: validate-crs kuttl ## Run e2e tests with manager deployed on cluster.
	mkdir -p $(PROJECT_DIR)/build/kuttl-test-artifacts
	SKIP_MANAGER_START=1 $(KUTTL) test

test-upgrade: kuttl bundle-post-process ## Run OLM-based operator upgrade tests.
	mkdir -p $(PROJECT_DIR)/build/kuttl-test-artifacts
	SKIP_MANAGER_START=1 \
	NEW_PRODUCT_VERSION=$$(make --quiet --no-print-directory -C .. tag) \
	KUTTL=$(KUTTL) $(KUTTL) test --config kuttl-test.yaml tests/upgrade

stackrox-image-pull-secret: ## Create default image pull secret for StackRox images on DockerHub. Used by Helm chart.
# Create stackrox namespace if not exists.
	echo '{ "apiVersion": "v1", "kind": "Namespace", "metadata": { "name": "stackrox" } }' | kubectl apply -f -
# Create stackrox image pull secret in stackrox namespace.
	$(PROJECT_DIR)/../deploy/common/pull-secret.sh stackrox docker.io | kubectl -n stackrox apply -f -

.PHONY: check-ci-setup
check-ci-setup: ## Make sure this target is not started in CI environment.
	@if [ -n "$$CI" ]; then echo "Setup error: operator should be installed and started by OLM." >&2; exit 1; fi

.PHONY: run
run: check-ci-setup manifests generate ensure-rox-main-image-exists fmt vet ## Run operator from your host without deploying it to a cluster.
	ENABLE_WEBHOOKS=$(ENABLE_WEBHOOKS) ../scripts/go-run.sh ./main.go

.PHONY: ensure-rox-main-image-exists
ensure-rox-main-image-exists:
	./hack/ensure-rox-image-exist.sh

##@ Build

build: generate fmt vet ## Build operator local binary.
	../scripts/go-build-file.sh ./main.go bin/manager

docker-build: test ## Build docker image with the operator.
	@( \
		`# status.sh file is used by scripts/go-build.sh which we try to run in the docker container.` \
		`# status.sh needs git repo, make and Makefile and who knows what else but its actual output is simple.` \
		`# Here we grab the output and create a new build/status.sh from it. That file will be copied into the` \
		`# container and used there without having to bring all dependencies.` \
		set -euxo pipefail ;\
		smuggled_status_sh="$(PROJECT_DIR)/build/status.sh" ;\
		mkdir -p "$(PROJECT_DIR)/build" ;\
		cd "$(PROJECT_DIR)/.." ;\
		echo "#!/bin/sh" > "$${smuggled_status_sh}" ;\
		./status.sh | awk '{print "echo \"" $$0 "\""}' >> "$${smuggled_status_sh}" ;\
		chmod +x "$${smuggled_status_sh}" ;\
		`# Verify that the resulting status.sh is actually runnable` \
		"$${smuggled_status_sh}" ;\
	)
	DOCKER_BUILDKIT=1 BUILDKIT_PROGRESS=plain docker build \
		--build-arg "GO_VERSION=${GO_VERSION}" \
		-t ${IMG} \
		--ssh default=${SSH_AUTH_SOCK} \
		--build-arg ROX_IMAGE_FLAVOR=$(ROX_IMAGE_FLAVOR) \
		-f Dockerfile \
		..

##@ Deployment

.PHONY: install
install: check-ci-setup manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: check-ci-setup manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete --ignore-not-found -f -

.PHONY: deploy
deploy: check-ci-setup manifests kustomize ## Deploy operator image to the K8s cluster specified in ~/.kube/config.
	rm -rf config/local-deploy-versioned && \
		mkdir config/local-deploy-versioned && \
		cd config/local-deploy-versioned && \
		$(KUSTOMIZE) create --resources ../default && \
		$(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/local-deploy-versioned | kubectl create -f -

.PHONY: undeploy
undeploy: check-ci-setup kustomize ## Undeploy operator image from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found -f -

.PHONY: deploy-via-olm
deploy-via-olm: ## Deploy operator image to the cluster using OLM.
# This requires operator, bundle and index images to already be pushed, but we do not depend on the
# target(s) that do it here, because CI does it in a separate job, in parallel with cluster deployment.
	KUTTL=$(KUTTL) ./hack/olm-operator-install.sh $(TEST_NAMESPACE) $(VERSION)

.PHONY: deploy-dirty-tag-via-olm
deploy-dirty-tag-via-olm: ## Deploy operator dirty tagged image to the cluster using OLM.
# This target ignores filter for dirty tagged images.
	KUTTL=$(KUTTL) ./hack/olm-operator-install.sh --allow-dirty-tag $(TEST_NAMESPACE) $(VERSION)

.PHONY: deploy-previous-via-olm
deploy-previous-via-olm: kuttl bundle-post-process ## Deploy replaced version of operator image to the cluster using OLM.
# This requires operator, bundle and index images to already be pushed, but we do not depend on the
# target(s) that do it here, because CI does it in a separate job, in parallel with cluster deployment.
	@replaced_version_no_v=$$(sed -E -n 's/^[[:space:]]*replaces:[[:space:]]*[^.]+\.v(.*)$$/\1/p' build/bundle/manifests/rhacs-operator.clusterserviceversion.yaml) ;\
	set -x ;\
	KUTTL=$(KUTTL) ./hack/olm-operator-install.sh $(TEST_NAMESPACE) $(VERSION) $${replaced_version_no_v}

.PHONY: upgrade-via-olm
upgrade-via-olm: kuttl
	KUTTL=$(KUTTL) ./hack/olm-operator-upgrade.sh $(TEST_NAMESPACE) $(VERSION)

.PHONY: upgrade-dirty-tag-via-olm
upgrade-dirty-tag-via-olm: kuttl
	KUTTL=$(KUTTL) ./hack/olm-operator-upgrade.sh --allow-dirty-tag $(TEST_NAMESPACE) $(VERSION)

##@ Bundle and Index build

# Commands to enter local Python virtual environment and get needed dependencies there.
ACTIVATE_PYTHON = python3 -m venv bundle_helpers/.venv ;\
	. bundle_helpers/.venv/bin/activate ;\
	pip3 install --upgrade pip==21.3.1 setuptools==59.6.0 ;\
	pip3 install -r bundle_helpers/requirements.txt

.PHONY: bundle
bundle: manifests kustomize operator-sdk ## Generate bundle manifests and metadata, then validate generated files.
# Version is hardcoded to 0.0.1 here because otherwise git-versioned files are changed which we want to avoid.
# The correct version is updated later.
# Likewise, we hardcode the image reference to docker.io/stackrox/stackrox-operator. If this is overridden via
# the IMG_REPO environment variable, the final reference will be injected in the bundle-post-process step.
	$(OPERATOR_SDK) generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=docker.io/stackrox/stackrox-operator:0.0.1
	cd config/scorecard-versioned && $(KUSTOMIZE) edit set image scorecard-test=quay.io/operator-framework/scorecard-test:v$(OPERATOR_SDK_VERSION)
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle -q --overwrite --version 0.0.1 $(BUNDLE_METADATA_OPTS)
# Delete lines that copy the original bundle files in bundle.Dockerfile (we can't just use `rm` in scratch image).
# Post-processed files will be copied instead, as configured in bundle.Dockerfile.extra.
	sed -i'.bak' -e '/# Copy files to locations specified by labels./d' bundle.Dockerfile
	sed -i'.bak' -E -e '/^COPY .* \/(manifests|metadata|tests\/scorecard)\/$$/d' bundle.Dockerfile
	rm -f bundle.Dockerfile.bak
# We cannot add changes to bundle.Dockerfile directly since `operator-sdk generate bundle` would wipe them.
# Yet we want most of the contents autogenerated from the Makefile variables as a single source of truth.
# Therefore we append ".extra" file to the end of bundle's dockerfile.
	cat bundle.Dockerfile.extra >> bundle.Dockerfile
# Run a python script to fix the orders in the specDescriptors (children must not appear before their parents).
	set -euo pipefail ;\
	$(ACTIVATE_PYTHON) ;\
	bundle_helpers/fix-spec-descriptor-order.py \
	  <bundle/manifests/rhacs-operator.clusterserviceversion.yaml \
	  >bundle/manifests/rhacs-operator.clusterserviceversion.yaml.fixed
	mv bundle/manifests/rhacs-operator.clusterserviceversion.yaml.fixed \
       bundle/manifests/rhacs-operator.clusterserviceversion.yaml
	$(OPERATOR_SDK) bundle validate ./bundle --select-optional suite=operatorframework

.PHONY: bundle-post-process
bundle-post-process: operator-sdk ## Post-process CSV file to include correct operator versions, etc.
# Copy the original bundle to build/ directory.
	mkdir -p build/
	rm -rf build/bundle
	cp -a bundle build/
# Run Python script to post-process necessary files in the build/bundle/. Configure Python dependencies if not done already.
	set -euo pipefail ;\
	$(ACTIVATE_PYTHON) ;\
	./bundle_helpers/patch-csv.py \
		--use-version $(VERSION) \
		--first-version 3.62.0 `# 3.62.0 is the first operator version ever released` \
		--operator-image $(IMG) \
		--no-related-images \
		< bundle/manifests/rhacs-operator.clusterserviceversion.yaml \
	| sed 's,replaces: rhacs-operator.v3.66.0,replaces: rhacs-operator.v3.66.1,' > build/bundle/manifests/rhacs-operator.clusterserviceversion.yaml
# TODO(ROX-8618): make the above sed on `replaces:` line generic (quick hack to address ROX-8630)
# Check that the resulting bundle still passes validations.
	$(OPERATOR_SDK) bundle validate ./build/bundle --select-optional suite=operatorframework

.PHONY: bundle-build
bundle-build: bundle.Dockerfile bundle-post-process ## Build the bundle image.
	docker build -f $< -t $(BUNDLE_IMG) .

.PHONY: bundle-test
bundle-test: operator-sdk bundle-post-process ## Run scorecard tests against bundle files.
	$(OPERATOR_SDK) scorecard $(SCORECARD_ARGS) build/bundle/

.PHONY: bundle-test-image
bundle-test-image: operator-sdk ## Run scorecard tests against bundle image in the registry. Bundle image must be pushed beforehand.
	$(OPERATOR_SDK) scorecard $(SCORECARD_ARGS) $(BUNDLE_IMG)

.PHONY: index-build
index-build: bundle-post-process ## Build Index (a.k.a. Catalog) image with the bundle. Bundle image must be pushed beforehand.
# If the bundle has `replaces:` attribute in CSV, then the one it replaces has to already be in the index. Otherwise,
# validation of the index fails because more than one version is identified as the head of the channel. An example of an
# error: `multiple channel heads found in graph: rhacs-operator.v3.67.0, rhacs-operator.v3.69.0`
# We assume that the index for the replaced bundle is already in the registry with the same tag and so we use previous
# index image as a base when building the current index image. The sed expression below should extract everything after
# the first dot in `replaces:` value. E.g. `  replaces: rhacs-operator.v3.62.0` -> `v3.62.0`
	@set -euo pipefail ;\
	replaced_version=$$(sed -E -n 's/^[[:space:]]*replaces:[[:space:]]*[^.]+\.(.*)$$/\1/p' build/bundle/manifests/rhacs-operator.clusterserviceversion.yaml) ;\
	./hack/build-index-image.sh \
		--base-index-tag "$(INDEX_IMG_BASE):$${replaced_version}" \
		--index-tag "$(INDEX_IMG)" \
		--bundle-tag "$(BUNDLE_IMG)" \
		--replaced-version "$${replaced_version}" \
		--clean-output-dir

##@ Push images

.PHONY: docker-push
docker-push: ## Push docker image with the operator.
# Use push-as-manifest-list.sh here to make sure manifests for the operator image are created. The script also retries on error.
# We're not using push-as-manifest-list.sh for pushing other two images because they are independent from architecture.
	../scripts/ci/push-as-manifest-list.sh "$(IMG)"

.PHONY: docker-push-bundle
docker-push-bundle: ## Push docker image with the bundle.
	docker push "$(BUNDLE_IMG)"

.PHONY: docker-push-index
docker-push-index: ## Push docker image with the index.
	docker push "$(INDEX_IMG)"
