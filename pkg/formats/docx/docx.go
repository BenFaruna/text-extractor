package docx

import (
	"baliance.com/gooxml/document"
	"context"
	"fmt"
	"github.com/BenFaruna/text-extractor/internal/processor"
	"github.com/BenFaruna/text-extractor/pkg/extractor"
	"io"
	"time"
)

// Extractor implements the TextExtractor interface for Doc files
type Extractor struct{}

func New() *Extractor {
	return &Extractor{}
}

func (e *Extractor) Extract(ctx context.Context, r io.Reader, opts ...extractor.Option) (string, error) {
	options := extractor.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	//TODO implement me
	panic("implement me")
}

func (e *Extractor) ExtractFile(ctx context.Context, filePath string, opts ...extractor.Option) (string, error) {
	options := extractor.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	doc, err := document.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	extractCh := make(chan string, 1)
	errCh := make(chan error, 1)

	defer close(errCh)
	defer close(extractCh)

	go func() {
		var content string

		paragraphs := doc.Paragraphs()
		for _, p := range paragraphs {
			runs := p.Runs()
			for _, r := range runs {
				content = content + r.Text() + "\n"
			}

			if options.MaxContentLength > 0 && len(content) > options.MaxContentLength {
				content = content[:options.MaxContentLength]
				break
			}
		}

		if options.DetectEncoding {
			temp := []byte(content)
			encoding := processor.DetectEncoding(temp)

			if encoding != "utf-8" && options.DefaultEncoding == "utf-8" {
				temp, err := processor.ConvertToUTF8(temp, encoding)
				if err != nil {
					errCh <- fmt.Errorf("failed to convert to UTF-8: %w", err)
				}
				content = string(temp)
			}
		}
		extractCh <- processor.NormalizeText(content, options.PreserveLineBreaks)
	}()

	if options.ExtractionTimeout > 0 {
		_, cancel := context.WithTimeout(ctx, time.Duration(options.ExtractionTimeout)*time.Second)
		defer cancel()
	}

	select {
	case content := <-extractCh:
		return content, nil
	case err := <-errCh:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
