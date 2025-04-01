package extractor

// Options contains configurations for text extraction
type Options struct {
	PreserveLineBreaks bool
	PreserveFormatting bool
	MaxContentLength   int
	DetectEncoding     bool
	DefaultEncoding    string
	IncludeMetadata    bool
	ExtractionTimeout  int // in seconds
}

// DefaultOptions returns the default extraction options
func DefaultOptions() Options {
	return Options{
		PreserveLineBreaks: true,
		PreserveFormatting: false,
		MaxContentLength:   0, // unlimited
		DetectEncoding:     true,
		DefaultEncoding:    "utf-8",
		IncludeMetadata:    false,
		ExtractionTimeout:  60,
	}
}

// Option is a function that modifies Options
type Option func(*Options)

// WithPreserveLineBreaks sets whether to preserve line breaks
func WithPreserveLineBreaks(preserve bool) Option {
	return func(o *Options) {
		o.PreserveLineBreaks = preserve
	}
}

// WithPreserveFormatting sets whether to preserve text formatting
func WithPreserveFormatting(preserve bool) Option {
	return func(o *Options) {
		o.PreserveFormatting = preserve
	}
}

// WithMaxContentLength sets the maximum content length to extract
func WithMaxContentLength(maxLength int) Option {
	return func(o *Options) {
		o.MaxContentLength = maxLength
	}
}

// WithEncoding sets the encoding options
func WithEncoding(detect bool, defaultEncoding string) Option {
	return func(o *Options) {
		o.DetectEncoding = detect
		o.DefaultEncoding = defaultEncoding
	}
}

// WithMetadata sets whether to include metadata
func WithMetadata(include bool) Option {
	return func(o *Options) {
		o.IncludeMetadata = include
	}
}

// WithTimeout sets the extraction timeout in seconds
func WithTimeout(seconds int) Option {
	return func(o *Options) {
		o.ExtractionTimeout = seconds
	}
}
