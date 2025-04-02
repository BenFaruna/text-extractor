package pdf

import (
	"context"
	"strings"
	"testing"
)

type testCase struct {
	filePath string
	expected string
}

var testCases = []testCase{
	{"../../../testdata/test_file_1.pdf", "This is a pdf that will be used for test cases."},
	{"../../../testdata/test_file_2.pdf", "This is a pdf that will be used for test cases.\n\nTrying new line to preserve linebreaks."},
}

func TestExtractor_PDFExtractFile(t *testing.T) {
	extractor := &Extractor{}
	for _, tc := range testCases {
		path := strings.Split(tc.filePath, "/")
		t.Run(path[4], func(t *testing.T) {
			output, err := extractor.ExtractFile(context.Background(), tc.filePath)
			if err != nil {
				t.Errorf("ExtractFile() error = %v", err)
			}

			if output != tc.expected {
				t.Errorf("ExtractFile() output = %q, want %q", output, tc.expected)
			}
		})
	}
}
