REGISTRY         ?= quay.io
ORG              ?= hegemone
TAG              ?= latest
IMAGE            ?= $(REGISTRY)/$(ORG)/kore
BUILD_DIR        = "${GOPATH}/src/github.com/hegemone/kore/build"
SOURCES          := $(shell find . -name '*.go' -not -path "*/vendor/*" -not -path "*/extensions/*")
TEST_SOURCES	 := $(shell go list ./... | grep -v extension)
PROJECT_ROOT := ""$(abspath $(lastword $(MAKEFILE_LIST)))/..""
.DEFAULT_GOAL    := build

vendor:
	@dep ensure

plugins: $(PLUGIN_SOURCES)
	@go build -buildmode=plugin -o ${BUILD_DIR}/bacon.plugins.kore.nsk.io.so -i -ldflags="-s -w" ./pkg/extension/plugin/bacon.go

adapters: $(ADAPTER_SOURCES)
	@go build -buildmode=plugin -o ${BUILD_DIR}/ex-discord.adapters.kore.nsk.io.so -i -ldflags="-s -w" ./pkg/extension/adapter/discord.go
	@go build -buildmode=plugin -o ${BUILD_DIR}/ex-irc.adapters.kore.nsk.io.so -i -ldflags="-s -w" ./pkg/extension/adapter/irc.go

kore: $(SOURCES) adapters plugins
	@go build -o ${BUILD_DIR}/kore -i -ldflags="-s -w" ./cmd/kore

build: kore test
	@echo > /dev/null

test:
	@go test -cover ${TEST_SOURCES}

clean:
	@rm -rf ${BUILD_DIR}

run: kore
	@KORE_PLUGIN_DIR=${PROJECT_ROOT}/build \
	KORE_ADAPTER_DIR=${PROJECT_ROOT}/build \
	./build/kore

image:
	docker build -t ${IMAGE}:${TAG} ${PROJECT_ROOT}

run-image:
	docker run -it ${IMAGE}:${TAG}

push:
	docker push ${IMAGE}:${TAG}

.PHONY: vendor image push clean run image run-image push
