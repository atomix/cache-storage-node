export CGO_ENABLED=0
export GO111MODULE=on

ifdef VERSION
STORAGE_VERSION := $(VERSION)
else
STORAGE_VERSION := latest
endif

.PHONY: build

all: build

build: # @HELP build the source code
build: deps license_check linters
	GOOS=linux GOARCH=amd64 go build -o build/_output/atomix-memory-storage-node ./cmd/atomix-memory-storage-node
	GOOS=linux GOARCH=amd64 go build -o build/_output/atomix-memory-storage-driver ./cmd/atomix-memory-storage-driver

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

test: # @HELP run the unit tests and source code validation
test: build license_check linters
	go test github.com/atomix/cache-storage-node/...

linters: # @HELP examines Go source code and reports coding problems
	GOGC=75  golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	./build/bin/license-check

images: # @HELP build Docker images
images: build
	docker build . -f build/atomix-memory-storage-node/Dockerfile -t atomix/atomix-memory-storage-node:${STORAGE_VERSION}
	docker build . -f build/atomix-memory-storage-driver/Dockerfile -t atomix/atomix-memory-storage-driver:${STORAGE_VERSION}

kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image atomix/atomix-memory-storage-node:${STORAGE_VERSION}
	kind load docker-image atomix/atomix-memory-storage-driver:${STORAGE_VERSION}

clean: # @HELP clean build files
	@rm -rf vendor build/_output

push: # @HELP push atomix-memory-storage-node Docker image
	docker push atomix/atomix-memory-storage-node:${STORAGE_VERSION}
	docker push atomix/atomix-memory-storage-driver:${STORAGE_VERSION}
