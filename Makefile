BINARY_NAME := kubectl-eso
MODULE := github.com/josegonzalez/kubectl-eso
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w \
	-X $(MODULE)/pkg/cmd.Version=$(VERSION) \
	-X $(MODULE)/pkg/cmd.Commit=$(COMMIT) \
	-X $(MODULE)/pkg/cmd.Date=$(DATE)
GOPATH := $(shell go env GOPATH)

.PHONY: build install test clean lint fmt vet

build:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/kubectl-eso

install: build
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	cp scripts/kubectl_complete-eso $(GOPATH)/bin/kubectl_complete-eso
	chmod +x $(GOPATH)/bin/kubectl_complete-eso

test:
	go test ./... -race -count=1

clean:
	rm -rf bin/ dist/

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet
