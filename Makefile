.SHELLFLAGS += -x -e -o pipefail
SHELL = /bin/bash
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)
TESTARGS=-v -p 1 -race -cover -coverprofile=coverage.txt -covermode=atomic

GEN_API_GO = gen-go-api
AUTHENTIK_MAIN = authentik_main

default: gen

# Run acceptance tests
.PHONY: test
test:
	TF_ACC=1 go test $(TESTARGS) ./...
	go tool cover -html coverage.txt -o coverage.html

build:
	go build -o /dev/null -v ./...

gen:
	golangci-lint run -v
	AUTHENTIK_URL="" go tool github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	terraform fmt -recursive .

gen-client-go:
	mkdir -p ${PWD}/${AUTHENTIK_MAIN}
ifeq ($(wildcard ${PWD}/${AUTHENTIK_MAIN}/.*),)
	git clone --depth 1 https://github.com/goauthentik/authentik.git ${PWD}/${AUTHENTIK_MAIN}
else
	cd ${PWD}/${AUTHENTIK_MAIN}
	git pull
endif
	make -C ${PWD}/${AUTHENTIK_MAIN} gen-client-go
	go mod edit -replace goauthentik.io/api/v3=./${AUTHENTIK_MAIN}/${GEN_API_GO}
