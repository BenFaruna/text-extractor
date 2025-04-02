package text

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	filePath string
	expected string
}

var testCases = []testCase{
	{
		filePath: "../../../testdata/test_file_3.txt",
		expected: "This is a text file.\n\nThis file contains multiple lines.\n\nWhen this is parsed, it should be normalised and whitespaces removed.",
	}, {
		filePath: "../../../testdata/test_file_4.txt",
		expected: "太陽能板張開，古文見\n每當看、航天器升空完成星箭分離進入預定軌道後不久就會聽到：太陽能板打開、的宣告。這宣告在古時曰",
	},
}

func TestExtractor_TextExtractFile(t *testing.T) {
	extractor := &Extractor{}
	for _, tc := range testCases {
		path := strings.Split(tc.filePath, "/")
		t.Run(path[4], func(t *testing.T) {
			output, err := extractor.ExtractFile(context.Background(), tc.filePath)
			if err != nil {
				t.Errorf("ExtractFile() error = %v", err)
			}

			if output != tc.expected {
				t.Errorf("ExtractFile() output = %s, want %s", output, tc.expected)
			}
		})
	}
}

func TestExtractor_TextExtract(t *testing.T) {
	extractor := &Extractor{}
	for _, tc := range testCases {
		path := strings.Split(tc.filePath, "/")
		t.Run(path[4], func(t *testing.T) {
			r := createReader(t, tc.filePath)
			output, err := extractor.Extract(context.Background(), r)
			if err != nil {
				t.Errorf("ExtractFile() error = %v", err)
			}

			if output != tc.expected {
				t.Errorf("ExtractFile() output = %s, want %s", output, tc.expected)
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
