.SHELLFLAGS += -x -e -o pipefail
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)
TESTARGS=-v -p 1 -race -coverprofile=coverage.txt -covermode=atomic

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

# see https://github.com/goauthentik/authentik/blob/main/Makefile#LL99-L113
gen-client-go:
	mkdir -p ${PWD}/${AUTHENTIK_MAIN}
ifeq ($(wildcard ${PWD}/${AUTHENTIK_MAIN}/.*),)
	git clone --depth 1 https://github.com/goauthentik/authentik.git ${PWD}/${AUTHENTIK_MAIN}
else
	cd ${PWD}/${AUTHENTIK_MAIN}
	git pull
endif
	mkdir -p ${PWD}/${GEN_API_GO}
ifeq ($(wildcard ${PWD}/${GEN_API_GO}/.*),)
	git clone --depth 1 https://github.com/goauthentik/client-go.git ${PWD}/${GEN_API_GO}
else
	cd ${PWD}/${GEN_API_GO}
	git pull
endif
	cp ${PWD}/${AUTHENTIK}/schema.yml ${PWD}/${GEN_API_GO}/schema.yml
	# cp ${PWD}/schema.yml ${PWD}/${GEN_API_GO}
	make -C ${PWD}/${GEN_API_GO} build
	go mod edit -replace goauthentik.io/api/v3=./${GEN_API_GO}
