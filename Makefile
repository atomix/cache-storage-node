export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ATOMIX_LOCAL_REPLICA_VERSION := latest

all: build

build: # @HELP build the source code
build:
	GOOS=linux GOARCH=amd64 go build -o build/_output/local-replica ./cmd/local-replica

test: # @HELP run the unit tests and source code validation
test: build license_check linters
	go test github.com/atomix/local-replica/...

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

image: # @HELP build local-replica Docker image
image: build
	docker build . -f build/docker/Dockerfile -t atomix/local-replica:${ATOMIX_LOCAL_REPLICA_VERSION}

push: # @HELP push local-replica Docker image
	docker push atomix/local-replica:${ATOMIX_LOCAL_REPLICA_VERSION}
