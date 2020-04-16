# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GODOG=godog
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=clarity
BINARY_DIR=build
BINARY_WIN=$(BINARY_NAME).exe
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_DARWIN=$(BINARY_NAME)_osx

all: test build
build: build-darwin build-win build-linux 
test: gen unit
unit: 
	$(GOTEST) ./pkg/... -v
integration: 
	$(GOTEST) ./test/integration/... -v
e2e: 
	$(GOTEST) ./test/e2e/... -v
clean: 
	$(GOCLEAN)
	find . -name "*.test" | xargs rm 
	rm -fr $(BINARY_DIR)
gen:
	go generate ./...
dep:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
build-darwin: 
	cd ./cmd/clarity && \
	mkdir -p $(BINARY_DIR) && \
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GODOG) -output $(BINARY_DIR)/$(BINARY_DARWIN)
build-win:
	cd ./cmd/clarity && \
	mkdir -p $(BINARY_DIR) && \
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GODOG) -output $(BINARY_DIR)/$(BINARY_WIN)
build-linux:
	cd ./cmd/clarity && \
	mkdir -p $(BINARY_DIR) && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GODOG) -output $(BINARY_DIR)/$(BINARY_UNIX)
release:
	./bin/create_new_release.sh

.PHONY: all test clean build
