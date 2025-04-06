package html

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
)

type TestCase struct {
	filePath string
	expected string
}

var testCases = []TestCase{
	{filePath: "../../../testdata/test_file_8.html", expected: "This is the first paragraph\n\nAnd there is a follow-up paragraph\n\nThis is the second paragraph"},
	{filePath: "../../../testdata/test_file_9.html", expected: "1. List item 1\n2. List item 2\n3. List item 3\n\n<!--THE END-->\n\n- [List item 1](/)\n- List item 2\n- List item 3"},
}

func TestExtractor_HTMLExtractFile(t *testing.T) {
	extractor := New()
	for _, tc := range testCases {
		path := strings.Split(tc.filePath, "/")
		t.Run(path[4], func(t *testing.T) {
			output, err := extractor.ExtractFile(context.Background(), tc.filePath)
			if err != nil {
				t.Errorf("failed to extract html: %v", err)
			}

			if output != tc.expected {
				t.Errorf("expected: %s, got: %s", tc.expected, output)
			}
		})
	}
}

func TestExtractor_HTMLExtract(t *testing.T) {
	extractor := New()
	for _, tc := range testCases {
		path := strings.Split(tc.filePath, "/")
		t.Run(path[4], func(t *testing.T) {
			r := createReader(t, tc.filePath)
			output, err := extractor.Extract(context.Background(), r)
			if err != nil {
				t.Errorf("failed to extract html: %v", err)
			}

			if output != tc.expected {
				t.Errorf("expected: %s, got: %s", tc.expected, output)
			}
		})
	}
}

func createReader(t testing.TB, filePath string) io.Reader {
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return f
}
