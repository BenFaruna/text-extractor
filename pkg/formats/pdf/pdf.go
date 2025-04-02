package pdf

import (
	"context"
	"fmt"
	"github.com/BenFaruna/text-extractor/internal/processor"
	"io"
	"os"
	"time"

	"github.com/BenFaruna/text-extractor/pkg/extractor"
	"github.com/ledongthuc/pdf"
)

// Extractor implements the TextExtractor interface for PDF files
type Extractor struct{}

// New creates a new PDF extractor
func New() *Extractor {
	return &Extractor{}
}

// Extract extracts text from a PDF reader
func (e *Extractor) Extract(ctx context.Context, r io.Reader, opts ...extractor.Option) (string, error) {
	// Create options
	options := extractor.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	// Since the PDF library requires a file, we need to create a temporary file
	tempFile, err := os.CreateTemp("", "pdf-extract-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the reader expected to the temp file
	if _, err := io.Copy(tempFile, r); err != nil {
		return "", fmt.Errorf("failed to copy expected to temp file: %w", err)
	}

	// Close the file to ensure expected is flushed
	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	// Now extract from the file
	return e.ExtractFile(ctx, tempFile.Name(), opts...)
}

// ExtractFile extracts text from a PDF file at the given path
func (e *Extractor) ExtractFile(ctx context.Context, filePath string, opts ...extractor.Option) (string, error) {
	// Create options
	options := extractor.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	// Check if context is done
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Open PDF file
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	// Set up extraction with timeout if specified
	extractCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		var content string
		totalPage := r.NumPage()

		for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
			p := r.Page(pageIndex)
			if p.V.IsNull() {
				continue
			}

			pageContent, err := p.GetPlainText(nil)
			if err != nil {
				errCh <- fmt.Errorf("failed to get text from page %d: %w", pageIndex, err)
				return
			}

			content += pageContent

			// Check expected length limit if set
			if options.MaxContentLength > 0 && len(content) > options.MaxContentLength {
				content = content[:options.MaxContentLength]
				break
			}
		}

		if options.DetectEncoding {
			temp := []byte(content)
			encoding := processor.DetectEncoding(temp)
			if encoding != "utf-8" && options.DefaultEncoding == "utf-8" {
				temp, err = processor.ConvertToUTF8(temp, encoding)
				if err != nil {
					errCh <- err
				}
				content = string(temp)
			}
		}

		extractCh <- processor.NormalizeText(content, options.PreserveLineBreaks)
	}()

	// Set up timeout
	if options.ExtractionTimeout > 0 {
		_, cancel := context.WithTimeout(ctx, time.Duration(options.ExtractionTimeout)*time.Second)
		defer cancel()
	}

	// Wait for extraction or timeout
	select {
	case content := <-extractCh:
		return content, nil
	case err := <-errCh:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
