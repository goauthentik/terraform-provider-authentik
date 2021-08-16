.SHELLFLAGS += -x -e
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)

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

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	rm -f examples_local/.terraform.lock.hcl
	cd examples_local && terraform init

gen:
	golint ./...
	go generate

gen-api:
	docker run \
		--rm -v ${PWD}:/local \
		--user ${UID}:${GID} \
		openapitools/openapi-generator-cli generate \
		--git-host github.com \
		--git-repo-id goauthentik \
		--git-user-id terraform-provider-authentik \
		-i /local/schema.yml \
		-g go \
		-o /local/api \
		--additional-properties=packageName=api,enumClassPrefix=true,useOneOfDiscriminatorLookup=true
	rm -f api/go.mod api/go.sum
