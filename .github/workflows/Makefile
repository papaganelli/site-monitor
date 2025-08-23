.PHONY: build test clean run fmt lint help install stats history status demo

# Variables
BINARY_NAME=site-monitor

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BINARY_NAME) .

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and show coverage
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME) coverage.out coverage.html site-monitor.db

run: build ## Build and run the monitoring daemon
	@echo "Running $(BINARY_NAME) in monitoring mode..."
	@./$(BINARY_NAME) run

stats: build ## Show monitoring statistics
	@echo "Showing monitoring statistics..."
	@./$(BINARY_NAME) stats

history: build ## Show monitoring history
	@echo "Showing monitoring history..."
	@./$(BINARY_NAME) history --limit 20

status: build ## Show current monitoring status
	@echo "Showing current status..."
	@./$(BINARY_NAME) status

demo: build ## Run a quick demo of all CLI commands
	@echo "🚀 Site Monitor CLI Demo"
	@echo ""
	@echo "📋 Available commands:"
	@./$(BINARY_NAME) --help
	@echo ""
	@echo "📊 Statistics (if data available):"
	@./$(BINARY_NAME) stats --since 1h 2>/dev/null || echo "   No data yet - run 'make run' first"
	@echo ""
	@echo "📋 History (if data available):"
	@./$(BINARY_NAME) history --limit 5 2>/dev/null || echo "   No data yet - run 'make run' first"
	@echo ""
	@echo "🔍 Status (if data available):"
	@./$(BINARY_NAME) status 2>/dev/null || echo "   No data yet - run 'make run' first"

fmt: ## Format Go code
	@echo "Formatting code..."
	@gofmt -s -w .
	@go mod tidy

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

install: build ## Install binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ $(BINARY_NAME) installed successfully"
	@echo "You can now run:"
	@echo "  site-monitor run      # Start monitoring"
	@echo "  site-monitor stats    # View statistics"
	@echo "  site-monitor history  # View history"
	@echo "  site-monitor status   # View current status"

uninstall: ## Remove binary from /usr/local/bin
	@echo "Removing $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) uninstalled successfully"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"