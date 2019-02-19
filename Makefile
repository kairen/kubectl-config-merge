VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

ORG := github.com
OWNER := kairen
REPOPATH ?= $(ORG)/$(OWNER)/kubectl-config-merge

GOOS ?= $(shell go env GOOS)

$(shell mkdir -p ./out)

.PHONY: build
build: out/config-merge

.PHONY: out/config-merge
out/config-merge: 
	CGO_ENABLED=0 GOOS=$(GOOS) go build \
	  -ldflags="-s -w -X $(REPOPATH)/pkg/version.version=$(VERSION)" \
	  -a -o $@ cmd/kubectl-config-merge.go

.PHONY: dep 
dep:
	dep ensure

.PHONY: clean
clean:
	rm -rf out/

