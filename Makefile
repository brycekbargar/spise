SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
WORKTREE_ROOT := $(shell git rev-parse --show-toplevel 2> /dev/null)

.DEFAULT_GOAL := build
.PHONY: build
build: spise

spise:
	@go build

.PHONY: test
pkgs?='./...'
test:
	@go test $(pkgs)
