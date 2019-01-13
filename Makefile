# Phony Targets
.PHONY: all build test coverage coveralls clean

all: build

build: hashi

test:
	go test -covermode=count -coverprofile=cover.out ./pkg/...

coverage: cover.out
	go tool cover -func=cover.out

coverage-html: cover.out
	go tool cover -html=cover.out

coveralls: cover.out
	goveralls -coverprofile=cover.out -service=travis-ci

clean:
	rm -f bin/hashi
	rm -f cover.out

# File Targets
hashi:
	go build -o bin/hashi ./cmd/hashi

cover.out: test
