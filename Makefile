BINARY = allowed-signatures
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GITHUB_TOKEN?= "TBD"
SOURCE ?= $(shell pwd)

VERSION?="0.0.1"

COMMIT=$(shell git rev-parse HEAD)


GO=go

all: clean test vet

install:
	$(GO) build ${LDFLAGS} -o ${GOPATH}/bin/${BINARY} *.go

test:
	$(GO) test -v --race -count=1 ./... 2>&1 | tee ${TEST_REPORT} ;

vet:
	$(GO) vet ./... > ${VET_REPORT} 2>&1 ;

fmt:
	$(GO) fmt $$($(GO) list ./... | grep -v /vendor/) ;

.PHONY: coverage
coverage:
	$(GO) test -race -v -count=1 -coverprofile=cover.out ./...
	$(GO) tool cover -html=cover.out

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}


gorelease:
	docker run --rm --privileged \
		-v $(shell pwd):/go/src/github.com/frankywahl/allowed-signers \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-w /go/src/github.com/frankywahl/allowed-signers \
		-e GITHUB_TOKEN=${GITHUB_TOKEN} \
		-e SOURCE=${SOURCE} \
		goreleaser/goreleaser release \
			--rm-dist \
			--snapshot

GOLINT ?= docker run --rm --tty --volume $(shell pwd):/app --workdir /app golangci/golangci-lint:latest golangci-lint
.PHONY: lint
lint:
	$(GOLINT) run --config .github/golangci.yml --verbose
