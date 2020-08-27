GO := $(go env GOBIN)
ifeq ($(GOBIN),)
    GO := go
endif

GOOS 		= linux
CGO_ENABLED = 0
ROOT_DIR 	= $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/

.PHONY: clean

BUILD_TARGET := bin/

build: sapi

.PHONY: release
release: sapi;

.PHONY: release_darwin
release_darwin: darwin sapi;

darwin:
	$(eval GOOS := darwin)

sapi: ; $(info ======== compiled sapi:)
	env CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) go build -mod vendor -a -installsuffix cgo -o $(BUILD_TARGET)sapi $(LD_FLAGS) $(ROOT_DIR)cmd/sapi/*.go

clean:
	rm -rf $(BUILD_TARGET)/*

UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Darwin)
	.DEFAULT_GOAL := release_darwin
else
	.DEFAULT_GOAL := release
endif