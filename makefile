PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)

.PHONY: build build-server docker release check

build:
	CGO_ENABLED=0 go build -o ./bin/server github.com/moov-io/fincen/cmd/server

build-webui:
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js ./cmd/webui/assets/wasm_exec.js
	GOOS=js GOARCH=wasm go build -o ./cmd/webui/assets/fincen.wasm github.com/moov-io/fincen/cmd/webui/fincen/
	CGO_ENABLED=0 go build -o ./bin/webui ./cmd/webui

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	GOOS=js GOARCH=wasm GOCYCLO_LIMIT=115 COVER_THRESHOLD=80.0 time ./lint-project.sh
endif

.PHONY: client
client:
# Versions from https://github.com/OpenAPITools/openapi-generator/releases
	@chmod +x ./openapi-generator
	@rm -rf ./client
	OPENAPI_GENERATOR_VERSION=4.2.0 ./openapi-generator generate -i openapi.yaml -g go -o ./client
	rm -f client/go.mod client/go.sum ./client/.travis.yml
	go fmt ./...
	go build github.com/moov-io/fincen/client
	go test ./client

.PHONY: clean
clean:
	@rm -rf ./bin/ ./tmp/ coverage.txt misspell* staticcheck lint-project.sh

dist: clean build
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=1 GOOS=windows go build -o bin/fincen.exe github.com/moov-io/fincen/cmd/server
else
	CGO_ENABLED=0 GOOS=$(PLATFORM) go build -o bin/fincen-$(PLATFORM)-amd64 github.com/moov-io/fincen/cmd/server
endif

docker: clean docker-hub docker-webui

docker-hub:
	docker build --pull -t moov/fincen:$(VERSION) -f Dockerfile .
	docker tag moov/fincen:$(VERSION) moov/fincen:latest

docker-openshift:
	docker build --pull -t quay.io/moov/fincen:$(VERSION) -f Dockerfile.openshift --build-arg VERSION=$(VERSION) .
	docker tag quay.io/moov/fincen:$(VERSION) quay.io/moov/fincen:latest

docker-webui:
	docker build --pull -t moov/fincen-webui:$(VERSION) -f Dockerfile.webui .
	docker tag moov/fincen-webui:$(VERSION) moov/fincen-webui:latest

.PHONY: clean-integration test-integration

clean-integration:
	docker-compose kill
	docker-compose rm -v -f

test-integration: clean-integration
	docker-compose up -d
	sleep 5
	curl -v http://localhost:8088/files

release: docker AUTHORS
	go vet ./...
	go test -coverprofile=cover-$(VERSION).out ./...
	git tag -f $(VERSION)

release-push:
	docker push moov/fincen:$(VERSION)
	docker push moov/fincen:latest
	docker push moov/fincen-webui:$(VERSION)

quay-push:
	docker push quay.io/moov/fincen:$(VERSION)
	docker push quay.io/moov/fincen:latest

.PHONY: cover-test cover-web
cover-test:
	go test -coverprofile=cover.out ./...
cover-web:
	go tool cover -html=cover.out

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@

.PHONY: tagged-release
tagged-release:
	@./tagged-release.sh $(VERSION)
