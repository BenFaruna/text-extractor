# Text Extractor

A Go package for extracting text from various file formats including PDF, DOCX, HTML, and plain text files.

## Features

- Extract text from multiple file formats
- Configurable options for text extraction
- Command-line interface
- Simple and extensible API
- Format detection
- Timeout support

## Installation

```bash
go get github.com/BenFaruna/text-extractor
```

## Usage

### As a library

```go
package main

import (
	"context"
	"fmt"

	"github.com/BenFaruna/text-extractor/pkg/extractor"
	"github.com/BenFaruna/text-extractor/pkg/formats/pdf"
)

func main() {
	// Create an extractor for PDF files
	pdfExtractor := pdf.New()
	
	// Extract text from a file
	ctx := context.Background()
	text, err := pdfExtractor.ExtractFile(ctx, "document.pdf",
		extractor.WithPreserveLineBreaks(true),
		extractor.WithTimeout(30),
	)
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println(text)
}
```

### Using the command-line tool

```bash
# Extract text from a PDF file
./extractor --file=document.pdf --preserve-linebreaks=true

# Extract text and save to output file
./extractor --file=document.pdf --output=extracted.txt
```

## Supported Formats

- PDF (using github.com/ledongthuc/pdf)
- DOCX (Microsoft Word)
- HTML
- Plain text

## Adding New Format Support

To add support for a new format, create a new package under `pkg/formats/` and implement the `extractor.TextExtractor` interface.

## License

[MIT License](LICENSE)