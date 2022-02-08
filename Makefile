BINARY := docker-credfile-gen
SHELL := /bin/bash -o pipefail
ifndef GOBIN
	export GOBIN ?= $(shell pwd)/bin
endif

TEST ?= ./...
TESTUNITARGS ?= -timeout 10s -race -cover
TEST_REPORT ?= test-report.xml
BUILD_OUTPUT ?= bin/$(BINARY)

include scripts/Makefile.help
include scripts/Makefile.deps

### Build Targets

## Builds docker-credfile-gen to $(BUILD_OUTPUT) runs clean first.
.PHONY: build
build: clean $(BUILD_OUTPUT)

## Builds docker-credfile-gen to $(BUILD_OUTPUT) inside docker, runs clean first.
.PHONY: docker-build
docker-build: clean
	docker container rm -f temp > /dev/null 2>&1 || true
	docker build --build-arg=goos=$(OS) --build-arg=goarch=$(ARCH) -t $(BINARY):latest .
	docker container create --name temp $(BINARY):latest
	docker container cp temp:$(BINARY)/bin/$(BINARY) $(GOBIN)
	docker container rm -f temp

$(BUILD_OUTPUT):
	@ go build -o $(BUILD_OUTPUT) .

## Removes docker-credfile-gen binariess from $(BUILD_OUTPUT)
clean:
	@ rm -f $(BUILD_OUTPUT)

### Dev Targets

## Runs all the project linters
lint: deps
	@ echo "-> Running golint..."
	@ $(GOBIN)/golangci-lint run
	@ echo "-> Checking source file license headers..."
	@ $(GOBIN)/go-licenser -d

## Formats all Go files to the desired format.
.PHONY: format
format: deps
	@ echo "-> Formatting Go files..."
	@ $(GOBIN)/go-licenser
	@ $(GOBIN)/golangci-lint run --fix --deadline=5m
	@ echo "-> Done."

## Runs unit tests. Use TESTARGS and TEST to control which flags and packages are used and tested.
.PHONY: unit
unit:
	@ echo "-> Running unit tests for $(BINARY)..."
	@ go test $(TESTARGS) $(TESTUNITARGS) $(TEST)
