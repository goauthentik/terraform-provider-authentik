.SHELLFLAGS += -x -e
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)
TESTARGS=-v -p 1 -race -coverprofile=coverage.txt -covermode=atomic

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
	AUTHENTIK_URL="" go generate

# see https://github.com/goauthentik/authentik/blob/main/Makefile#LL99-L113
gen-client-go:
	mkdir -p gen-go-api gen-go-api/templates
	wget https://raw.githubusercontent.com/goauthentik/authentik/master/schema.yml -O gen-go-api/schema.yml
	wget https://raw.githubusercontent.com/goauthentik/client-go/main/config.yaml -O gen-go-api/config.yaml
	wget https://raw.githubusercontent.com/goauthentik/client-go/main/templates/README.mustache -O gen-go-api/templates/README.mustache
	wget https://raw.githubusercontent.com/goauthentik/client-go/main/templates/go.mod.mustache -O gen-go-api/templates/go.mod.mustache
	docker run \
		--rm -v ${PWD}/gen-go-api:/local \
		--user ${UID}:${GID} \
		docker.io/openapitools/openapi-generator-cli:v6.0.0 generate \
		-i /local/schema.yml \
		-g go \
		-o /local/ \
		-c /local/config.yaml
	go mod edit -replace goauthentik.io/api/v3=./gen-go-api
	rm -rf ./gen-go-api/config.yaml ./gen-go-api/templates/ ./gen-go-api/schema.yml
