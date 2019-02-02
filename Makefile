# Variables
SOURCES := $(shell find . -name '*.go' | grep -v vendor)

# Phony Targets
.PHONY: all lint build test coverage coveralls clean

all: build

lint:
	golint -set_exit_status internal/... pkg/...

build: hashi

test cover.out: $(SOURCES)
	go test -covermode=count -coverprofile=cover.out ./internal/... ./pkg/...

coverage: cover.out
	go tool cover -func=cover.out

coverage-html: cover.out
	go tool cover -html=cover.out

coveralls: cover.out
	goveralls -coverprofile=cover.out -service=travis-ci

clean:
	rm -f hashi
	rm -f cover.out

# File Targets
hashi: $(SOURCES)
	go build
