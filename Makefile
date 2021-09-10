.SHELLFLAGS += -x -e

HOSTNAME=registry.terraform.io
NAMESPACE=goauthentik
NAME=authentik
BINARY=terraform-provider-${NAME}
VERSION=99999
OS_ARCH=darwin_amd64

default: gen

# Run acceptance tests
.PHONY: test
test:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build: gen-api
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	rm -f examples_local/.terraform.lock.hcl
	cd examples_local && terraform init

gen:
	golint ./...
	go generate
