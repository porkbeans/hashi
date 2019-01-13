# Phony Targets
.PHONY: all lint build test coverage coveralls clean

all: build

lint:
	golint -set_exit_status cmd/... pkg/...

build: hashi

test:
	go test -covermode=count -coverprofile=cover.out ./pkg/...

coverage: test
	go tool cover -func=cover.out

coverage-html: test
	go tool cover -html=cover.out

coveralls: test
	goveralls -coverprofile=cover.out -service=travis-ci

clean:
	rm -f bin/hashi
	rm -f cover.out

# File Targets
hashi:
	go build -o bin/hashi ./cmd/hashi
