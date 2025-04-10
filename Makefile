TARGET := vagabond
VERSION := 0.1.0
GO := go
GOFMT := gofmt
LINTER := golangci-lint

.PHONY: build
build:
	$(GO) build -o $(TARGET) cmd/${TARGET}/main.go

.PHONY: install
install:
	$(GO) install

.PHONY: run
run: build
	./$(TARGET)

.PHONY: test
test:
	$(GO) test ./...

.PHONY: fmt
fmt:
	$(GOFMT) -w .

.PHONY: lint
lint:
	$(LINTER) run

.PHONY: docs
docs:
	@echo "Generating docs (if any)..."

.PHONY: version
version:
	@echo "${TARGET} version $(VERSION)"

.PHONY: clean
clean:
	rm -f $(TARGET)
	rm -rf migrations

.PHONY: release
release: clean
	$(GO) build -o $(TARGET)
	tar -czvf $(TARGET)-v$(VERSION).tar.gz $(TARGET)

.PHONY: help
help:
	@echo "Makefile for Vagabond"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build       - Build the project"
	@echo "  install     - Install the project (go install)"
	@echo "  run         - Build and run the project"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format Go files"
	@echo "  lint        - Run linters"
	@echo "  docs        - Generate documentation"
	@echo "  version     - Display the current version"
	@echo "  clean       - Clean up build artifacts"
	@echo "  release     - Create a release archive"
	@echo "  help        - Show this help message"
