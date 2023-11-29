GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
SOURCES := $(shell find . -type f  -name '*.go')

BUILD_ARCH ?= linux/$(GOARCH)
ifeq ($(BUILD_ARM),true)
ifneq ($(GOARCH),arm64)
	  BUILD_ARCH= linux/$(GOARCH),linux/arm64
endif
endif
ifeq ($(BUILD_X86),true)
ifneq ($(GOARCH),amd64)
	  BUILD_ARCH= linux/$(GOARCH),linux/amd64
endif
endif

REGISTRY_SERVER_ADDRESS?=release.daocloud.io
REGISTRY_REPO?=$(REGISTRY_SERVER_ADDRESS)/datatunerx
REGISTRY_USER_NAME ?=
REGISTRY_PASSWORD ?=

# Set your version by env or using latest tags from git
# The official tag version will be passed in from the environment variable. If the development version is not passed in, mark the dev tag
VERSION?=""
ifeq ($(VERSION), "")
    LATEST_DEV_TAG=$(shell ./hack/get-version.sh rayproxy)-dev-$(shell git rev-parse --short=8 HEAD)
    ifeq ($(LATEST_DEV_TAG),)
        # Forked repo may not sync tags from upstream, so give it a default tag to make CI happy.
        VERSION="unknown"
    else
        VERSION=$(LATEST_DEV_TAG)
    endif
endif

RAY_PROXY_VERSION := $(shell echo $(VERSION) | sed 's/-/+/1')

RAY_PROXY_IMAGE_VERSION := $(shell echo $(RAY_PROXY_VERSION) | sed 's/+/-/1')

.PHONY: docker-login
docker-login:docker-login
	@echo "push images to $(REGISTRY_REPO)"
	echo ${REGISTRY_PASSWORD} | docker login ${REGISTRY_SERVER_ADDRESS} -u ${REGISTRY_USER_NAME} --password-stdin

.PHONY: ray-proxy
ray-proxy: $(SOURCES) docker-login
	echo "Building ray-proxy for arch = $(BUILD_ARCH)"
	export DOCKER_CLI_EXPERIMENTAL=enabled ;\
	! ( docker buildx ls | grep ray-multi-platform-builder ) && docker buildx create --use --platform=$(BUILD_ARCH) --name kpanda-ingress-multi-platform-builder --driver-opt image=docker.m.daocloud.io/moby/buildkit:buildx-stable-1 ;\
	docker buildx build \
		  --build-arg kpanda_version=$(KPANDA_VERSION) \
			--build-arg UBUNTU_MIRROR=$(UBUNTU_MIRROR) \
			--builder ray-multi-platform-builder \
			--platform $(BUILD_ARCH) \
			--build-arg LDFLAGS=$(LDFLAGS) \
			--tag $(REGISTRY_REPO)/ray-proxy:$(RAY_PROXY_IMAGE_VERSION)  \
			--tag $(REGISTRY_REPO)/ray-proxy:latest  \
			-f Dockerfile \
			--push \
			.


