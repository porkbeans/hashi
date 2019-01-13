.PHONY: all build test cover clean

all: build

build: hashi

test:
	go test -cover -coverprofile=cover.out ./pkg/...
	go tool cover -func=cover.out
	rm -f cover.out

cover:
	go test -cover -coverprofile=cover.out ./pkg/...
	go tool cover -html=cover.out
	rm -f cover.out

clean:
	rm -f bin/hashi

hashi:
	go build -o bin/hashi ./cmd/hashi
