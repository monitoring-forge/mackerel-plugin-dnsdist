VERSION=0.0.6
GITCOMMIT?=$(shell git describe --dirty --always 2>/dev/null)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.commit=${GITCOMMIT}"
all: mackerel-plugin-dnsdist

.PHONY: mackerel-plugin-dnsdist

mackerel-plugin-dnsdist: cmd/mackerel-plugin-dnsdist/*.go
	go build $(LDFLAGS) -o mackerel-plugin-dnsdist ./cmd/mackerel-plugin-dnsdist/

linux: cmd/mackerel-plugin-dnsdist/*.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-dnsdist ./cmd/mackerel-plugin-dnsdist/

check:
	go test -v ./...

lint:
	golangci-lint run --timeout 5m ./...