BINARY_NAME := kubectl-eso
MODULE := github.com/josegonzalez/kubectl-eso
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w \
	-X $(MODULE)/pkg/cmd.Version=$(VERSION) \
	-X $(MODULE)/pkg/cmd.Commit=$(COMMIT) \
	-X $(MODULE)/pkg/cmd.Date=$(DATE)

.PHONY: build install test clean lint fmt vet go-lint yaml-lint markdown-lint krew-install

build:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/kubectl-eso

install: build
	@mkdir -p $(HOME)/.krew/bin
	cp bin/$(BINARY_NAME) $(HOME)/.krew/bin/$(BINARY_NAME)
	cp scripts/kubectl_complete-eso $(HOME)/.krew/bin/kubectl_complete-eso
	chmod +x $(HOME)/.krew/bin/kubectl_complete-eso

test:
	go test ./... -race -count=1

clean:
	rm -rf bin/ dist/

fmt:
	go fmt ./...

vet:
	go vet ./...

go-lint:
	golangci-lint run --timeout 5m

yaml-lint:
	uvx yamllint -c .yamllint.yml .

markdown-lint:
	uvx pymarkdownlnt fix .

lint: fmt vet go-lint yaml-lint markdown-lint

krew-install: build
	@rm -rf dist/krew-archive
	@mkdir -p dist/krew-archive
	cp bin/$(BINARY_NAME) dist/krew-archive/
	cp LICENSE dist/krew-archive/
	cp README.md dist/krew-archive/
	cp -r scripts dist/krew-archive/
	tar -czf dist/kubectl-eso-local.tar.gz -C dist/krew-archive .
	@printf '%s\n' \
		'apiVersion: krew.googlecontainertools.github.com/v1alpha2' \
		'kind: Plugin' \
		'metadata:' \
		'  name: eso' \
		'spec:' \
		'  version: $(VERSION)' \
		'  shortDescription: Manage Kubernetes Secrets with External Secrets Operator' \
		'  platforms:' \
		'    - selector:' \
		'        matchLabels:' \
		'          os: $(shell go env GOOS)' \
		'          arch: $(shell go env GOARCH)' \
		'      uri: https://example.com/local' \
		'      sha256: "0000000000000000000000000000000000000000000000000000000000000000"' \
		'      bin: $(BINARY_NAME)' \
		'      files:' \
		'        - from: $(BINARY_NAME)' \
		'          to: .' \
		'        - from: LICENSE' \
		'          to: .' \
		> dist/krew-local.yaml
	-kubectl krew uninstall eso 2>/dev/null
	kubectl krew install --manifest=dist/krew-local.yaml --archive=dist/kubectl-eso-local.tar.gz
	@rm -rf dist/krew-archive
