package processor

import (
	"bytes"
	"regexp"
	"strings"
	//"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	// Common regex patterns
	spaceRegex        = regexp.MustCompile(`\s+`)
	nonPrintableRegex = regexp.MustCompile(`[\x00-\x09\x0B\x0C\x0E-\x1F\x7F]`)
)

// NormalizeText standardizes newlines, spacing, and other text features
func NormalizeText(text string, preserveLineBreaks bool) string {
	// Replace non-printable characters
	text = nonPrintableRegex.ReplaceAllString(text, "")

	if preserveLineBreaks {
		// Standardize line breaks
		text = strings.ReplaceAll(text, "\r\n", "\n")
		text = strings.ReplaceAll(text, "\r", "\n")

		// Handle multiple newlines (preserve but don't allow more than 2)
		text = regexp.MustCompile(`\s{2,}`).ReplaceAllString(text, "\n\n")
		text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

		// Remove spaces at the beginning and end of lines
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimSpace(line)
		}
		text = strings.Join(lines, "\n")
	} else {
		// Replace all whitespace with a single space
		text = spaceRegex.ReplaceAllString(text, " ")
	}

	// Trim space from the start and end
	text = strings.TrimSpace(text)

	return text
}

// DetectEncoding attempts to determine text encoding from a byte slice
// This is a simple implementation - in a real-world scenario you might want
// to use a more sophisticated library like chardet
func DetectEncoding(data []byte) string {
	// Check for BOM (Byte Order Mark)
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return "utf-8"
	}
	if len(data) >= 2 {
		if data[0] == 0xFE && data[1] == 0xFF {
			return "utf-16be"
		}
		if data[0] == 0xFF && data[1] == 0xFE {
			return "utf-16le"
		}
	}

	// Check if it's valid UTF-8
	isASCII := true
	for _, b := range data {
		if b >= 128 {
			isASCII = false
			break
		}
	}

	if isASCII {
		return "ascii"
	}

	if utf8.Valid(data) {
		return "utf-8"
	}

	// This is a simplified approach - in a real implementation,
	// you would use more sophisticated heuristics or libraries
	return "unknown"
}

// ConvertToUTF8 converts text from the detected encoding to UTF-8
func ConvertToUTF8(data []byte, detectedEncoding string) ([]byte, error) {
	var decoder transform.Transformer

	switch strings.ToLower(detectedEncoding) {
	case "ascii", "utf-8":
		// Already in UTF-8 compatible format
		return data, nil
	case "utf-16be":
		decoder = unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder()
	case "utf-16le":
		decoder = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	case "iso-8859-1", "latin1":
		decoder = charmap.ISO8859_1.NewDecoder()
	// Add more encodings as needed
	default:
		// Default to UTF-8
		return data, nil
	}

	output := bytes.Buffer{}
	r := transform.NewReader(bytes.NewReader(data), decoder)
	_, err := output.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
