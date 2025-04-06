package html

import (
	"context"
	"io"
	"os"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"

	"github.com/BenFaruna/text-extractor/internal/processor"
	"github.com/BenFaruna/text-extractor/pkg/extractor"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
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
		htmlOptions := converter.WithContext(ctx)

		output, err := htmltomarkdown.ConvertReader(r, htmlOptions)
		if err != nil {
			errCh <- err
		}

		extractCh <- processor.NormalizeText(string(output), options.PreserveLineBreaks)
	}()

	if options.ExtractionTimeout > 0 {
		_, cancel := context.WithTimeout(ctx, time.Duration(options.ExtractionTimeout)*time.Second)
		defer cancel()
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case output := <-extractCh:
		return output, nil
	case err := <-errCh:
		return "", err
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

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return e.Extract(ctx, file, opts...)
}
