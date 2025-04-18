# Variables
TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=local
NAMESPACE=provider
NAME=psm
BINARY=terraform-provider-${NAME}
VERSION=0.5.14
OS_ARCH=linux_amd64
OS=linux

# Default target
default: install

# Build the provider binary
build:
	@echo "==> Building ${BINARY}..."
	go build -o ${BINARY}

# Build using go build -o psm
psm:
	@echo "==> Building using go build -o psm..."
	go build -o psm

# Release using goreleaser
release:
	@echo "==> Creating release..."
	goreleaser release --rm-dist --snapshot --skip-publish --skip-sign

# Install the provider locally for Terraform
install: build
	@echo "==> Installing ${BINARY}..."
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

# Run standard tests
test:
	@echo "==> Running tests..."
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

# Run acceptance tests
testacc:
	@echo "==> Running acceptance tests..."
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
