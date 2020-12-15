SHELL := /bin/bash

VERBOSE ?= 0
ifeq ($(VERBOSE), 0)
# Make is verbose in Linux: let's make it silent.
MAKEFLAGS += --silent
endif

VERSION := $(shell git describe --tags 2> /dev/null || echo "0.0.0")
BUILD := $(shell git rev-parse --short HEAD 2> /dev/null)
COMMIT := $(if $(shell git status --porcelain --untracked-files=no),${BUILD}-dirty,${BUILD})
PROJECTNAME := maintainers-generator
BASE := $(shell pwd)
OUTPUT := $(BASE)/bin

# Go variables
GOFILES := $(wildcard *.go)
LDFLAGS=-ldflags "-X=github.com/leodido/$(PROJECTNAME)/pkg/version.version=$(VERSION) -X=github.com/leodido/$(PROJECTNAME)/pkg/version.commit=$(COMMIT)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := "/tmp/.$(PROJECTNAME)-stderr.txt"

.PHONY: build
build: $(GOFILES)
	@echo "  >  Building binary..."
	touch $(STDERR)
	rm $(STDERR)
	mkdir -p $(OUTPUT)
	go build $(LDFLAGS) -o "$(OUTPUT)/$(PROJECTNAME)" $(GOFILES) 2> $(STDERR)
	cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

test:
	@echo "  >  Executing tests..."
	env GOTRACEBACK=all go test ./...

clean:
	@echo "  >  Cleaning build cache"
	rm $(OUTPUT)/$(PROJECTNAME) 2> /dev/null
	go clean
