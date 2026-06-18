# Hatch — Framework TUI declarativo para Go
# Targets: build, test, lint, run, clean

GO ?= go
BINARY ?= hatch
CMD_DIR ?= ./cmd/hatch

.PHONY: build test lint run clean coverage

build:
	$(GO) build -o bin/$(BINARY) $(CMD_DIR)

test:
	$(GO) test ./... -v -cover

lint:
	@which golangci-lint > /dev/null 2>&1 || (echo "Instalando golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

run:
	$(GO) run $(CMD_DIR) run canva/demo.hml

clean:
	rm -rf bin/
	$(GO) clean

coverage:
	$(GO) test ./... -coverprofile=coverage.out && $(GO) tool cover -html=coverage.out
