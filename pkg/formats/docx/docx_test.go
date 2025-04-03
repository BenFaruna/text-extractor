package docx

import (
	"context"
	"strings"
	"testing"

	e "github.com/BenFaruna/text-extractor/pkg/extractor"
)

type testCase struct {
	filePath string
	expected string
}

var testCases = []testCase{
	{"../../../testdata/test_file_5.docx", "Text extraction on steroids\nThis is a docx file. Let’s see if we can extract this without an external library.\n\nLet’s see how this goes. Looking forward to a great time."},
	{"../../../testdata/test_file_6.docx", "Heading!!!!\nThis is a docx file.\n\nThis is page two. Can we get the most from this?"},
	{"../../../testdata/test_file_7.docx", "List item 1\nList item 2\nList item 3\n\nList item numbering 1\nList item numbering 2\nList item numbering 3"},
}

func TestExtractor_DocxExtractFile(t *testing.T) {
	extractor := &Extractor{}
	for _, tc := range testCases {
		path := strings.Split(tc.filePath, "/")
		t.Run(path[4], func(t *testing.T) {
			output, err := extractor.ExtractFile(context.Background(), tc.filePath, e.WithPreserveFormatting(false))
			if err != nil {
				t.Errorf("ExtractFile() error = %v", err)
			}

			if output != tc.expected {
				t.Errorf("ExtractFile() output = %s, want %s", output, tc.expected)
			}
		})
	}
}
