# Copyright (c) 2024 Red Hat, Inc.

LOCALBIN = $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

BIN_GO ?= go
BIN_ENVTEST ?= $(LOCALBIN)/setup-envtest
BIN_CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
OPERATOR_RUN_ARGS ?=

VERSION_CONTROLLER_GEN = v0.15.0
ENVTEST_K8S_VERSION = 1.28.x

OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

kubeAssets = "KUBEBUILDER_ASSETS=$(shell $(BIN_ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)"

testCmd = "$(kubeAssets) $(BIN_GO) test -coverprofile=cov.out -v ./pkg/... -ginkgo.v"
ifdef TEST_NAME
testCmd += " -ginkgo.focus \"$(TEST_NAME)\""
endif

test: $(BIN_ENVTEST) generate_testdata
	@eval $(testCmd)

cov: test
	@$(BIN_GO) tool cover -func=cov.out
	@$(BIN_GO) tool cover -html=cov.out -o cov.html

.PHONY: generate_testdata
generate_testdata: $(BIN_CONTROLLER_GEN)
	@$(BIN_CONTROLLER_GEN) crd paths="./pkg/testdata/v1/..." output:dir="./pkg/testdata/crd"
	@rm -rf ./pkg/testdata/v1/zz_generated.deepcopy.go
	@$(BIN_CONTROLLER_GEN) object paths="./pkg/testdata/v1/..."

$(BIN_ENVTEST): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(BIN_GO) install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

$(BIN_CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(BIN_GO) install sigs.k8s.io/controller-tools/cmd/controller-gen@$(VERSION_CONTROLLER_GEN)
