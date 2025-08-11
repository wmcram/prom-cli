BINARY_NAME := promcli
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
    GOBIN := $(shell go env GOPATH)/bin
endif

.PHONY: build install clean

build:
	go build -o bin/promcli cmd/promcli.go

install: 
	go build -o "$(GOBIN)/$(BINARY_NAME)" cmd/promcli.go

clean:
	rm -rf bin