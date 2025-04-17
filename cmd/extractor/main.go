package main

import (
	"context"
	"flag"
	"fmt"
	. "github.com/BenFaruna/text-extractor/internal/logging"
	_ "github.com/BenFaruna/text-extractor/internal/logging"
	"github.com/BenFaruna/text-extractor/pkg/extractor"
	"github.com/BenFaruna/text-extractor/pkg/formats/docx"
	"github.com/BenFaruna/text-extractor/pkg/formats/html"
	"github.com/BenFaruna/text-extractor/pkg/formats/pdf"
	"github.com/BenFaruna/text-extractor/pkg/formats/text"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	filePath := flag.String("file", "", "Path to the file to extract text from")
	preserveLineBreaks := flag.Bool("preserve-linebreaks", true, "Preserve line breaks in the extracted text")
	preserveFormatting := flag.Bool("preserve-formatting", false, "Preserve text formatting")
	timeoutSeconds := flag.Int("timeout", 60, "Timeout in seconds")
	outputFile := flag.String("output", "", "Output file (if not provided, prints to stdout)")
	fileFormat := flag.String("file-format", "", "Input file format (e.g txt, docx, pdf, html)")

	flag.Parse()

	// Check if file path is provided
	if *filePath == "" {
		ErrorLogger.Println("File path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create the extractor manager
	manager := extractor.NewManager()

	// Register format extractors
	manager.Register("pdf", pdf.New())
	manager.Register("docx", docx.New())
	manager.Register("html", html.New())
	manager.Register("txt", text.New())

	var ext string

	if *fileFormat != "" {
		ext = strings.ToLower(*fileFormat)
	} else {
		// Determine format from file extension
		ext = strings.ToLower(filepath.Ext(*filePath))
		ext = strings.TrimPrefix(ext, ".")
	}

	// Get the appropriate extractor
	formatExtractor, ok := manager.Get(ext)
	if !ok {
		ErrorLogger.Printf("Unsupported file format: %s\n", ext)
		os.Exit(1)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeoutSeconds)*time.Second)
	defer cancel()

	// Extract text with options
	t, err := formatExtractor.ExtractFile(ctx, *filePath,
		extractor.WithPreserveLineBreaks(*preserveLineBreaks),
		extractor.WithPreserveFormatting(*preserveFormatting),
	)

	if err != nil {
		ErrorLogger.Printf("Text extraction failed: %v\n", err)
		os.Exit(1)
	}

	// Output the extracted text
	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(t), 0644)
		if err != nil {
			ErrorLogger.Printf("Writing to output file failed: %v\n", err)
			os.Exit(1)
		}
		InfoLogger.Printf("Text extracted successfully to: %s\n", *outputFile)
	} else {
		fmt.Println(t)
	}
	//fmt.Println(t)

	CloseLogger()
}
