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
	GOOS=linux GOARCH=amd64 go build -o build/cache-storage-node/_output/cache-storage-node ./cmd/cache-storage-node

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

test: # @HELP run the unit tests and source code validation
test: build license_check linters
	go test github.com/atomix/cache-storage/...

linters: # @HELP examines Go source code and reports coding problems
	GOGC=75  golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

images: # @HELP build cache-storage-node Docker image
images: build
	docker build . -f build/cache-storage-node/Dockerfile -t atomix/cache-storage-node:${STORAGE_VERSION}

push: # @HELP push cache-storage-node Docker image
	docker push atomix/cache-storage-node:${STORAGE_VERSION}
