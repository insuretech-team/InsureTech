package processor

import (
	"context"
	"fmt"
)

// OCRProcessor handles optical character recognition processing.
// This is a STUB implementation. To enable real OCR, install and configure Tesseract:
//
// 1. Install Tesseract:
//    - Ubuntu/Debian: sudo apt-get install tesseract-ocr libtesseract-dev
//    - macOS: brew install tesseract
//    - Windows: https://github.com/UB-Mannheim/tesseract/wiki
//
// 2. Enable CGO and use gosseract library:
//    - go get github.com/otiai10/gosseract/v2
//    - Set CGO_ENABLED=1 in build environment
//
// 3. Implement actual OCR extraction using gosseract:
//    import "github.com/otiai10/gosseract/v2"
//    client := gosseract.NewClient()
//    defer client.Close()
//    client.SetImageFromBytes(data)
//    text, err := client.Text()
//
type OCRProcessor struct {
	enabled bool
}

// NewOCRProcessor creates a new OCR processor.
// enabled: set to true to enable OCR (requires Tesseract + gosseract)
func NewOCRProcessor(enabled bool) *OCRProcessor {
	return &OCRProcessor{
		enabled: enabled,
	}
}

// ExtractText extracts text from an image using OCR.
// If enabled=false, returns empty string as a graceful no-op.
// If enabled=true, returns error indicating OCR is not configured.
func (op *OCRProcessor) ExtractText(ctx context.Context, data []byte, mimeType string) (string, error) {
	if !op.enabled {
		return "", nil
	}

	return "", fmt.Errorf("OCR not configured: install tesseract and enable CGO")
}

// IsEnabled returns whether OCR processing is enabled.
func (op *OCRProcessor) IsEnabled() bool {
	return op.enabled
}
