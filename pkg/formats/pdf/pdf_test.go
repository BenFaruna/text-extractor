package pdf

import (
	"context"
	"testing"
)

type testCase struct {
	filePath string
	content  string
}

var testCases = []testCase{
	{"../../../testdata/test_file_1.pdf", "This is a pdf that will be used for test cases."},
	{"../../../testdata/test_file_2.pdf", "This is a pdf that will be used for test cases.  Trying new line to preserve linebreaks. "},
}

func TestExtractor_ExtractFile(t *testing.T) {
	extractor := &Extractor{}
	for _, tc := range testCases {
		t.Run(tc.filePath, func(t *testing.T) {
			output, err := extractor.ExtractFile(context.Background(), tc.filePath)
			if err != nil {
				t.Errorf("ExtractFile() error = %v", err)
			}

			if output != tc.content {
				t.Errorf("ExtractFile() output = %q, want %q", output, tc.content)
			}
		})
	}
}
