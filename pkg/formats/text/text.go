package text

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BenFaruna/text-extractor/internal/processor"
	"github.com/BenFaruna/text-extractor/pkg/extractor"
	"io"
	"os"
	"time"
)

type Extractor struct{}

func New() *Extractor {
	return &Extractor{}
}

func (e *Extractor) Extract(ctx context.Context, r io.Reader, opts ...extractor.Option) (string, error) {
	options := extractor.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	extractCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		var content string
		buf := &bytes.Buffer{}
		_, err := buf.ReadFrom(r)
		if err != nil {
			errCh <- fmt.Errorf("failed to read from reader: %w", err)
			return
		}

		content = buf.String()
		if options.MaxContentLength > 0 && len(content) > options.MaxContentLength {
			content = content[:options.MaxContentLength]
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

	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("file open error: %w", err)
	}
	defer f.Close()

	extractCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		temp := &bytes.Buffer{}

		_, err = io.Copy(temp, f)
		if err != nil {
			errCh <- err
		}

		output := temp.Bytes()
		if options.DetectEncoding {
			encoding := processor.DetectEncoding(temp.Bytes())
			if encoding != "utf-8" {
				output, err = processor.ConvertToUTF8(temp.Bytes(), encoding)
				if err != nil {
					errCh <- err
				}
			}
		}

		extractCh <- processor.NormalizeText(string(output), options.PreserveLineBreaks)
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
