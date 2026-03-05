package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

// CompressionLevel defines gzip compression levels
type CompressionLevel int

const (
	// CompressionDefault uses default compression
	CompressionDefault CompressionLevel = gzip.DefaultCompression
	// CompressionBest uses best compression (slower but smaller)
	CompressionBest CompressionLevel = gzip.BestCompression
	// CompressionFast uses fast compression (faster but larger)
	CompressionFast CompressionLevel = gzip.BestSpeed
)

// gzipResponseWriter wraps http.ResponseWriter to compress responses
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	wroteHeader bool
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	if !w.wroteHeader {
		// Set compression headers only when we actually write the header
		w.ResponseWriter.Header().Set("Content-Encoding", "gzip")
		w.ResponseWriter.Header().Del("Content-Length")
		w.ResponseWriter.Header().Set("Vary", "Accept-Encoding")
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.Writer.Write(b)
}

// Header returns the header map for the writer
func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Flush implements http.Flusher for SSE streaming
func (w *gzipResponseWriter) Flush() {
	// Flush gzip writer
	if gw, ok := w.Writer.(*gzip.Writer); ok {
		gw.Flush()
	}
	// Flush underlying response writer
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// gzipWriterPool reuses gzip writers for performance
var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		gz, _ := gzip.NewWriterLevel(nil, gzip.DefaultCompression)
		return gz
	},
}

// Compression middleware compresses responses using gzip
func Compression(level CompressionLevel) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if client accepts gzip
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			// Skip compression for certain paths
			if shouldSkipCompression(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Create gzip writer with appropriate level
			var gz *gzip.Writer
			var err error
			if level == CompressionDefault {
				gz = gzipWriterPool.Get().(*gzip.Writer)
				gz.Reset(w)
				defer gzipWriterPool.Put(gz)
			} else {
				gz, err = gzip.NewWriterLevel(w, int(level))
				if err != nil {
					// If we can't create gzip writer, serve uncompressed
					next.ServeHTTP(w, r)
					return
				}
			}
			defer gz.Close()

			// Wrap response writer (headers will be set when WriteHeader is called)
			gzw := &gzipResponseWriter{
				Writer:         gz,
				ResponseWriter: w,
			}

			next.ServeHTTP(gzw, r)
		})
	}
}

// shouldSkipCompression returns true for paths that shouldn't be compressed
func shouldSkipCompression(path string) bool {
	// Skip compression for:
	// 1. Already compressed files
	// 2. Streaming endpoints (they handle their own flushing)
	// 3. Small responses (compression overhead not worth it)
	skipPrefixes := []string{
		"/debug/logs/stream", // SSE stream
		"/metrics/",          // Already efficient format
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Skip for file extensions that are already compressed
	skipSuffixes := []string{
		".jpg", ".jpeg", ".png", ".gif", ".webp", // Images
		".mp4", ".webm", ".ogg", // Videos
		".zip", ".gz", ".bz2", // Archives
		".woff", ".woff2", ".ttf", // Fonts
	}

	for _, suffix := range skipSuffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}

	return false
}

// CompressionStats tracks compression metrics
type CompressionStats struct {
	OriginalBytes   int64
	CompressedBytes int64
	Requests        int64
}

// CompressionRatio returns the compression ratio (0.0 to 1.0)
func (cs *CompressionStats) CompressionRatio() float64 {
	if cs.OriginalBytes == 0 {
		return 0
	}
	return 1.0 - (float64(cs.CompressedBytes) / float64(cs.OriginalBytes))
}

// SavingsPercent returns savings percentage (e.g., 75.5 for 75.5% savings)
func (cs *CompressionStats) SavingsPercent() float64 {
	return cs.CompressionRatio() * 100
}
