VERSION=0.0.6
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
all: mackerel-plugin-dnsdist

.PHONY: mackerel-plugin-dnsdist linux check lint

mackerel-plugin-dnsdist: cmd/mackerel-plugin-dnsdist/*.go
	go build $(LDFLAGS) -o mackerel-plugin-dnsdist ./cmd/mackerel-plugin-dnsdist/

linux: cmd/mackerel-plugin-dnsdist/*.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-dnsdist ./cmd/mackerel-plugin-dnsdist/

check:
	go test -v ./...

lint:
	golangci-lint run --timeout 5m ./...