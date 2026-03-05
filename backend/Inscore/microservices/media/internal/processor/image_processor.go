package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
)

// ImageProcessor handles image processing operations using pure Go.
// Uses github.com/disintegration/imaging for image manipulation.
type ImageProcessor struct{}

// NewImageProcessor creates a new image processor.
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// ResizeImage resizes an image to the specified dimensions.
// Supports JPEG and PNG formats only.
func (ip *ImageProcessor) ResizeImage(data []byte, width, height int, mimeType string) ([]byte, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid dimensions: width and height must be positive")
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	if format != "jpeg" && format != "png" {
		return nil, fmt.Errorf("unsupported image format: %s (only JPEG and PNG supported)", format)
	}

	resized := imaging.Resize(img, width, height, imaging.Lanczos)

	return encodeImage(resized, mimeType)
}

// GenerateThumbnail generates a thumbnail of the specified size.
// Supports JPEG and PNG formats only.
func (ip *ImageProcessor) GenerateThumbnail(data []byte, width, height int) ([]byte, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid thumbnail dimensions: width and height must be positive")
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	if format != "jpeg" && format != "png" {
		return nil, fmt.Errorf("unsupported image format: %s (only JPEG and PNG supported)", format)
	}

	// Thumbnail creates a cropped and resized image to fit dimensions
	thumbnail := imaging.Thumbnail(img, width, height, imaging.Lanczos)

	// Encode as JPEG by default for thumbnails
	return encodeImage(thumbnail, "image/jpeg")
}

// CompressImage compresses an image to the specified quality.
// Quality is 1-100 for JPEG, ignored for PNG.
// Supports JPEG and PNG formats only.
func (ip *ImageProcessor) CompressImage(data []byte, quality int, mimeType string) ([]byte, error) {
	if quality < 1 || quality > 100 {
		quality = 80 // default quality
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	if format != "jpeg" && format != "png" {
		return nil, fmt.Errorf("unsupported image format: %s (only JPEG and PNG supported)", format)
	}

	buf := new(bytes.Buffer)

	switch mimeType {
	case "image/jpeg":
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, fmt.Errorf("failed to encode JPEG: %w", err)
		}
	case "image/png":
		if err := png.Encode(buf, img); err != nil {
			return nil, fmt.Errorf("failed to encode PNG: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported MIME type for compression: %s (only image/jpeg and image/png supported)", mimeType)
	}

	return buf.Bytes(), nil
}

// GetImageDimensions returns the width and height of an image.
// Supports JPEG and PNG formats only.
func (ip *ImageProcessor) GetImageDimensions(data []byte) (width, height int, err error) {
	img, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode image config: %w", err)
	}

	if format != "jpeg" && format != "png" {
		return 0, 0, fmt.Errorf("unsupported image format: %s (only JPEG and PNG supported)", format)
	}

	return img.Width, img.Height, nil
}

// encodeImage encodes an image to the specified MIME type.
func encodeImage(img image.Image, mimeType string) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch mimeType {
	case "image/jpeg":
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 85}); err != nil {
			return nil, fmt.Errorf("failed to encode JPEG: %w", err)
		}
	case "image/png":
		if err := png.Encode(buf, img); err != nil {
			return nil, fmt.Errorf("failed to encode PNG: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported MIME type: %s (only image/jpeg and image/png supported)", mimeType)
	}

	return buf.Bytes(), nil
}
