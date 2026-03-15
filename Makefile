.PHONY: all build run test clean

BINARY=anixops
VERSION?=1.0.0
GO?=go

all: build

build:
	$(GO) build -ldflags="-X 'main.version=$(VERSION)'" -o $(BINARY) cmd/anixops/main.go

build-tui:
	$(GO) build -o anixops-tui cmd/anixops-tui/main.go

build-all: build build-tui

run:
	$(GO) run cmd/anixops/main.go

run-tui:
	$(GO) run cmd/anixops-tui/main.go

test:
	$(GO) test -v ./...

test-coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

clean:
	rm -f $(BINARY) anixops-tui
	rm -rf data/

deps:
	$(GO) mod download
	$(GO) mod tidy

lint:
	golangci-lint run

docker-build:
	docker build -t anixops:$(VERSION) .

docker-run:
	docker run -p 8080:8080 -p 50052:50052 anixops:$(VERSION)

install: build
	cp $(BINARY) /usr/local/bin/

uninstall:
	rm -f /usr/local/bin/$(BINARY)