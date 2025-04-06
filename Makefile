.PHONY: build clean test test-pdf test-text test-doc

build:
	@echo "Building extractor..."
	@go build -o extractor ./cmd/extractor

clean:
	@echo "Cleaning up..."
	@rm -f cmd/extractor/extractor extractor
	@rm -f *.txt
	@find pkg/formats -name "*.log" -type f -delete 2>/dev/null || true

test:
	@echo "Running all tests..."
	@go test ./...

test-pdf:
	@echo "Running PDF tests..."
	@go test ./pkg/formats/pdf -v

test-text:
	@echo "Running text tests..."
	@go test ./pkg/formats/text -v

test-doc:
	@echo "Running document tests..."
	@go test ./pkg/formats/docx -v

test-html:
	@echo "Running html tests..."
	@go test ./pkg/formats/html -v