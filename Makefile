BINARY = allowed-signatures
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GITHUB_TOKEN?= "TBD"
SOURCE ?= $(shell pwd)

VERSION?="0.0.1"

COMMIT=$(shell git rev-parse HEAD)
GO=go
GOLINT=golangci-lint
GORELEASER=goreleaser

ifeq ($(DOCKER), true)
	GO=docker run --rm --volume $(shell pwd):/app --workdir /app golang:latest go
	GOLINT=docker run --rm --volume $(shell pwd):/app --workdir /app golangci/golangci-lint:latest golangci-lint
	GORELEASER=docker run --rm --privileged \
		--volume $(shell pwd):/go/src/github.com/frankywahl/allowed-signers \
		--volume /var/run/docker.sock:/var/run/docker.sock \
		--workdir /go/src/github.com/frankywahl/allowed-signers \
		--env GITHUB_TOKEN=${GITHUB_TOKEN} \
		--env SOURCE=${SOURCE} \
		goreleaser/goreleaser
endif

.PHONY: default
default: help

%-docker: ## Run a make command using docker
	@DOCKER=true $(MAKE) $*

.PHONY: clean
clean: ## Clean all dependencies
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}

.PHONY: coverage
coverage: ## Run coverage analysis with go cover
	$(GO) test -race -v -count=1 -coverprofile=cover.out ./...
	$(GO) tool cover -html=cover.out

.PHONY: fmt
fmt: ## Run fmt on go files
	$(GO) fmt $$($(GO) list ./... | grep -v /vendor/) ;

.PHONY: gorelease
gorelease: ## Create a fake release. Use Github Actions for a real release
		$(GORELEASER) release --rm-dist --snapshot

.PHONY: help
help: ## Show this help
	@echo "Usage: make <target>"
	@echo
	@echo "Targets:"
	@grep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: install
install: ## Install the binary on the machine
	$(GO) build ${LDFLAGS} -o ${GOPATH}/bin/${BINARY} *.go

.PHONY: lint
lint: ## Lint the files
	$(GOLINT) run --config .github/golangci.yml --verbose

.PHONY: test
test: ## Run the test suite
	$(GO) test -v --race -count=1 ./... 2>&1 | tee ${TEST_REPORT} ;

.PHONY: vendor
vendor: ## Vendor dependencies locally
	$(GO) mod vendor

.PHONY: vet
vet: ## Run vet on go files
	$(GO) vet ./... > ${VET_REPORT} 2>&1 ;

