.SHELLFLAGS += -x -e

HOSTNAME=registry.terraform.io
NAMESPACE=goauthentik
NAME=authentik
BINARY=terraform-provider-${NAME}
VERSION=99999
OS_ARCH=darwin_amd64
TESTARGS=-v -p 1 -race -coverprofile=coverage.txt -covermode=atomic

default: gen

# Run acceptance tests
.PHONY: test
test:
	TF_ACC=1 go test $(TESTARGS) ./...
	go tool cover -html coverage.txt -o coverage.html

build: gen-api
	go build -o ${BINARY}

gen:
	golangci-lint run -v
	go generate
