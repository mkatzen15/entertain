VERSION = 1.0.0
BUILD_TIME = $(shell date +%Y-%m-%dT%H:%M:%S:%z)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X github.com/mkatzen15/entertain/core.Version=${VERSION} \
-X github.com/mkatzen15/entertain/core.BuildTime=${BUILD_TIME}"

all: lint build run swagger

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build ${LDFLAGS} -o bin/entertain main.go

.PHONY: clean
clean:
	rm -rf *.o

.PHONY: swagger
swagger:
	swagger generate spec -m -o ./swagger.yaml

.PHONY: lint
lint:
	go vet ./...
	gofumpt -l -w .
