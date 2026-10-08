.DEFAULT_GOAL := help

SHELL := /bin/bash

GO111MODULE := on

GOPKG := github.com/veraison/eat

TOPDIR := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))

GOLINT_ARGS ?= run --timeout=3m -E dupl -E gocritic -E prealloc

GOLINT_VERSION = v2.13.2
GOLINT = $(TOPDIR)/tools-bin/golangci-lint
GOLINT_STAMP = $(TOPDIR)/tools-bin/golangci-lint-$(GOLINT_VERSION).stamp

$(GOLINT): $(GOLINT_STAMP)

$(GOLINT_STAMP):
	mkdir -p $(dir $(GOLINT))
	touch $(GOLINT_STAMP)
	GOBIN=$(dir $(GOLINT)) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLINT_VERSION)

.PHONY: lint
lint: $(GOLINT)
	$(GOLINT) $(GOLINT_ARGS)

ifeq ($(MAKECMDGOALS),lint)
GOLINT_ARGS ?= run --timeout=3m
else
  ifeq ($(MAKECMDGOALS),lint-extra)
  GOLINT_ARGS ?= run --timeout=3m --issues-exit-code=0 -E dupl -E gocritic -E lll -E prealloc
  endif
endif

ifeq ($(MAKECMDGOALS),test)
GOTEST_ARGS ?= -v -race $(GOPKG)
else
  ifeq ($(MAKECMDGOALS),test-cover)
  GOTEST_ARGS ?= -short -cover $(GOPKG)
  endif
endif

COVER_THRESHOLD := $(shell grep '^name: cover' .github/workflows/ci-go-cover.yml | cut -c13-)

.PHONY: test test-cover
test test-cover: ; go test $(GOTEST_ARGS)

presubmit:
	@echo
	@echo ">>> Check that the reported coverage figures are $(COVER_THRESHOLD)"
	@echo
	$(MAKE) test-cover
	@echo
	@echo ">>> Fix any lint error"
	@echo
	$(MAKE) lint-extra

.PHONY: licenses
licenses: ; @./scripts/licenses.sh

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  * test:       run unit tests for $(GOPKG)"
	@echo "  * test-cover: run unit tests and measure coverage for $(GOPKG)"
	@echo "  * lint:       lint sources using default configuration"
	@echo "  * presubmit:  check you are ready to push your local branch to remote"
	@echo "  * help:       print this menu"
	@echo "  * licenses:   check licenses of dependent packages"
