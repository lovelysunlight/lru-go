GO ?= go
GOFMT ?= gofmt "-s"
CURRENT_DIR := $(shell pwd)
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /examples/)
GOFILES := $(shell find . -name "*.go")
TESTFOLDER := $(shell $(GO) list ./... | grep -v examples | grep -v benchmark)
BENCHFOLDER := $(shell $(GO) list .)/benchmark
BENCHMARK_OUT := $(CURRENT_DIR)/benchmark.out
TESTTAGS ?= ""

.PHONY: test
# Run tests to verify code functionality.
test:
	echo "mode: count" > coverage.out
	for d in $(TESTFOLDER); do \
		$(GO) test $(TESTTAGS) -v -covermode=count -coverprofile=profile.out $$d > tmp.out; \
		cat tmp.out; \
		if grep -q "^--- FAIL" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "build failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "setup failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		fi; \
		if [ -f profile.out ]; then \
			cat profile.out | grep -v "mode:" >> coverage.out; \
			rm profile.out; \
		fi; \
	done

.PHONY: bench
# runs benchmark.
bench: SHELL:=/bin/bash
bench:
	pushd $(CURRENT_DIR)/benchmark > /dev/null; \
	$(GO) test -benchmem -bench . -v $(BENCHFOLDER) > $(BENCHMARK_OUT); \
	popd > /dev/null; \
	cat $(BENCHMARK_OUT); \
	if grep -q "^--- FAIL" $(BENCHMARK_OUT); then \
		rm $(BENCHMARK_OUT); \
		exit 1; \
	elif grep -q "build failed" $(BENCHMARK_OUT); then \
		rm $(BENCHMARK_OUT); \
		exit 1; \
	elif grep -q "setup failed" $(BENCHMARK_OUT); then \
		rm $(BENCHMARK_OUT); \
		exit 1; \
	fi; \

.PHONY: fmt
# Ensure consistent code formatting.
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
# format (check only).
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: vet
# Examine packages and report suspicious constructs if any.
vet:
	$(GO) vet $(VETPACKAGES)

.PHONY: lint
# Inspect source code for stylistic errors or potential bugs.
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: misspell
# Correct commonly misspelled English words in source code.
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

.PHONY: misspell-check
# misspell (check only).
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: tools
# Install tools (golint and misspell).
tools:
	@if [ $(GO_VERSION) -gt 15 ]; then \
		$(GO) install golang.org/x/lint/golint@latest; \
		$(GO) install github.com/client9/misspell/cmd/misspell@latest; \
	elif [ $(GO_VERSION) -lt 16 ]; then \
		$(GO) install golang.org/x/lint/golint; \
		$(GO) install github.com/client9/misspell/cmd/misspell; \
	fi

.PHONY: help
# Help.
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf " - \033[36m%-20s\033[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help