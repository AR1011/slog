.PHONY: all stdout file http channel multiple advanced clean

# Run all examples
all: stdout file http channel multiple advanced

# Run individual examples
stdout:
	@echo "Running stdout example..."
	@cd examples/stdout && go run main.go

file:
	@echo "Running file example..."
	@cd examples/file && go run main.go

http:
	@echo "Running http example..."
	@cd examples/http && go run main.go

channel:
	@echo "Running channel example..."
	@cd examples/channel && go run main.go

multiple:
	@echo "Running multiple example..."
	@cd examples/multiple && go run main.go

advanced:
	@echo "Running advanced example..."
	@cd examples/advanced && go run main.go

# Clean up generated log files
clean:
	@echo "Cleaning up log files..."
	@find . -name "*.log" -type f -delete
	@echo "Clean up complete"

# Help target
help:
	@echo "Available targets:"
	@echo "  all       - Run all examples"
	@echo "  stdout    - Run stdout example"
	@echo "  file      - Run file example"
	@echo "  http      - Run http example"
	@echo "  channel   - Run channel example"
	@echo "  multiple  - Run multiple example"
	@echo "  advanced  - Run advanced example"
	@echo "  clean     - Remove generated log files"
	@echo "  help      - Show this help message" 