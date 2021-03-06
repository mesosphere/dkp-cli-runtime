# Copyright 2022 D2iQ, Inc. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

# The GOPRIVATE environment variable controls which modules the go command considers
# to be private (not available publicly) and should therefore not use the proxy or checksum database
export GOPRIVATE ?= github.com/mesosphere

ALL_GO_SUBMODULES := $(shell PATH=$(PATH) find -mindepth 2 -maxdepth 2 -name go.mod -printf '%P\n' | sort)
GO_SUBMODULES_NO_TOOLS := $(filter-out tools/go.mod,$(ALL_GO_SUBMODULES))

define go_test
	gotestsum \
		--junitfile junit-report.xml \
		--junitfile-testsuite-name=relative \
		--junitfile-testcase-classname=short \
		-- \
		-covermode=atomic \
		-coverprofile=coverage.out \
		-race \
		-short \
		-v \
		$(if $(GOTEST_RUN),-run "$(GOTEST_RUN)") \
		./... && \
	go tool cover \
		-html=coverage.out \
		-o coverage.html
endef

.PHONY: test
test: ## Runs go tests for all modules in repository
test: install-tool.go.gotestsum
ifneq ($(wildcard $(REPO_ROOT)/go.mod),)
	$(info $(M) running tests$(if $(GOTEST_RUN), matching "$(GOTEST_RUN)") for root module)
	$(call go_test)
endif
ifneq ($(words $(GO_SUBMODULES_NO_TOOLS)),0)
	$(MAKE) $(addprefix test.,$(GO_SUBMODULES_NO_TOOLS:/go.mod=))
endif

.PHONY: test.%
test.%: ## Runs go tests for a specific module
test.%: ; $(info $(M) running tests$(if $(GOTEST_RUN), matching "$(GOTEST_RUN)") for module $*)
	cd $* && $(call go_test)

.PHONY: integration-test
integration-test: ## Runs integration tests for all modules in repository
	$(MAKE) GOTEST_RUN=Integration test

.PHONY: integration-test.%
integration-test.%: ## Runs integration tests for a specific module
	$(MAKE) GOTEST_RUN=Integration test.$*

.PHONY: bench
bench: ## Runs go benchmarks for all modules in repository
ifneq ($(wildcard $(REPO_ROOT)/go.mod),)
	$(info $(M) running benchmarks$(if $(GOTEST_RUN), matching "$(GOTEST_RUN)") for root module)
	go test $(if $(GOTEST_RUN),-run "$(GOTEST_RUN)") -race -cover -bench=. -v ./...
endif
ifneq ($(words $(GO_SUBMODULES_NO_TOOLS)),0)
	$(MAKE) $(addprefix bench.,$(GO_SUBMODULES_NO_TOOLS:/go.mod=))
endif

.PHONY: bench.%
bench.%: ## Runs go benchmarks for a specific module
bench.%: ; $(info $(M) running benchmarks$(if $(GOTEST_RUN), matching "$(GOTEST_RUN)") for module $*)
	cd $* && go test $(if $(GOTEST_RUN),-run "$(GOTEST_RUN)") -race -cover -v ./...

GOLANGCI_CONFIG_FILE ?= $(wildcard $(REPO_ROOT)/.golangci.y*ml)

.PHONY: lint
lint: ## Runs golangci-lint for all modules in repository
lint: install-tool.golangci-lint
ifneq ($(wildcard $(REPO_ROOT)/go.mod),)
lint: lint.root
endif
ifneq ($(words $(GO_SUBMODULES_NO_TOOLS)),0)
lint: $(addprefix lint.,$(GO_SUBMODULES_NO_TOOLS:/go.mod=))
endif

.PHONY: lint.%
lint.%: ## Runs golangci-lint for a specific module
lint.%: install-tool.golangci-lint; $(info $(M) running golangci-lint for $* module)
	$(if $(filter-out root,$*),cd $* && )golangci-lint run --fix --config=$(GOLANGCI_CONFIG_FILE)
	$(if $(filter-out root,$*),cd $* && )go fmt ./...
	$(if $(filter-out root,$*),cd $* && )go fix ./...

.PHONY: mod-tidy
mod-tidy:  ## Run go mod tidy for all modules
ifneq ($(wildcard $(REPO_ROOT)/go.mod),)
	$(info $(M) running go mod tidy for root module)
	go mod tidy -v -compat=1.17
	go mod verify
endif
ifneq ($(words $(ALL_GO_SUBMODULES)),0)
	$(MAKE) $(addprefix mod-tidy.,$(ALL_GO_SUBMODULES:/go.mod=))
endif

.PHONY: mod-tidy.%
mod-tidy.%: ## Runs go mod tidy for a specific module
mod-tidy.%: ; $(info $(M) running go mod tidy for module $*)
	cd $* && go mod tidy -v -compat=1.17
	cd $* && go mod verify

.PHONY: go-clean
go-clean: ## Cleans go build, test and modules caches for all modules
ifneq ($(wildcard $(REPO_ROOT)/go.mod),)
	$(info $(M) running go clean for root module)
	go clean -r -i -cache -testcache -modcache
endif
ifneq ($(words $(ALL_GO_SUBMODULES)),0)
	$(MAKE) $(addprefix go-clean.,$(ALL_GO_SUBMODULES:/go.mod=))
endif

.PHONY: go-clean.%
go-clean.%: ## Cleans go build, test and modules caches for a specific module
go-clean.%: ; $(info $(M) running go clean for module $*)
	cd $* && go clean -r -i -cache -testcache -modcache
