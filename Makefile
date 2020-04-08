export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ATOMIX_CACHE_STORAGE_VERSION := latest

all: build

build: # @HELP build the source code
build:
	GOOS=linux GOARCH=amd64 go build -o build/cache-storage/_output/cache-storage ./cmd/cache-storage
	GOOS=linux GOARCH=amd64 go build -o build/cache-controller/_output/cache-controller ./cmd/cache-controller

test: # @HELP run the unit tests and source code validation
test: build license_check linters
	go test github.com/atomix/cache-storage/...

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

images: # @HELP build cache-storage Docker image
images: build
	docker build . -f build/cache-storage/Dockerfile -t atomix/cache-storage:${ATOMIX_CACHE_STORAGE_VERSION}
	docker build . -f build/cache-controller/Dockerfile -t atomix/cache-controller:${ATOMIX_CACHE_STORAGE_VERSION}

push: # @HELP push cache-storage Docker image
	docker push atomix/cache-storage:${ATOMIX_CACHE_STORAGE_VERSION}
	docker push atomix/cache-controller:${ATOMIX_CACHE_STORAGE_VERSION}
