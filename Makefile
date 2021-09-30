GIT_TAG_NAME ?= $(shell git describe --abbrev=1 --tags)
GIT_HASH     ?= $(shell git rev-parse --verify HEAD)
BRANCH_NAME  ?= $(shell git rev-parse --abbrev-ref HEAD)
DATE         ?= $(shell date +%Y-%m-%d:%H%M:%s)

run:
	go run -mod vendor -ldflags "-s -w -X main.Version=$(GIT_TAG_NAME)" ./cmd/server

build:
	go build -mod vendor -ldflags "-s -w -X main.Version=$(GIT_TAG_NAME)" -o build/server ./cmd/server

