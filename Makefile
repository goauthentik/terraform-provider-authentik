.SHELLFLAGS += -x -e
PWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)
default: gen

# Run acceptance tests
.PHONY: test
test:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

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
		-i http://goauthentik.io/schema.yaml \
		-g go \
		-o /local/api \
		--additional-properties=packageName=api,enumClassPrefix=true,useOneOfDiscriminatorLookup=true
	rm -f api/go.mod api/go.sum
