package extractor

import (
	"context"
	"io"
)

// TextExtractor defines the interface for extracting text from various file formats
type TextExtractor interface {
	// Extract extracts text from the provided reader
	Extract(ctx context.Context, r io.Reader, opts ...Option) (string, error)

	// ExtractFile extracts text from a file at the given path
	ExtractFile(ctx context.Context, filePath string, opts ...Option) (string, error)
}

// Manager handles registration and retrieval of format extractors
type Manager struct {
	extractors map[string]TextExtractor
}

// NewManager creates a new manager for text extractors
func NewManager() *Manager {
	return &Manager{
		extractors: make(map[string]TextExtractor),
	}
}

// Register adds a text extractor for a specific format
func (m *Manager) Register(format string, extractor TextExtractor) {
	m.extractors[format] = extractor
}

// Get returns the extractor for the specified format
func (m *Manager) Get(format string) (TextExtractor, bool) {
	extractor, ok := m.extractors[format]
	return extractor, ok
}

// RegisteredFormats returns a map of all registered format extractors
func (m *Manager) RegisteredFormats() map[string]TextExtractor {
	return m.extractors
}
