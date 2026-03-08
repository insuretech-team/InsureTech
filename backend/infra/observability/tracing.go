package observability

import (
	"context"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig configures distributed tracing
type TracingConfig struct {
	ServiceName    string
	JaegerEndpoint string // e.g., "http://jaeger:14268/api/traces"
	SamplingRate   float64
	Enabled        bool
}

// DefaultTracingConfig returns default tracing configuration
func DefaultTracingConfig() *TracingConfig {
	return &TracingConfig{
		ServiceName:    getEnv("SERVICE_NAME", "insuretech-gateway"),
		JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
		SamplingRate:   1.0, // 100% sampling in development
		Enabled:        getEnv("ENABLE_TRACING", "true") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitTracer initializes OpenTelemetry tracing
func InitTracer(cfg *TracingConfig) (func(context.Context) error, error) {
	if cfg == nil {
		cfg = DefaultTracingConfig()
	}

	if !cfg.Enabled {
		return func(ctx context.Context) error { return nil }, nil
	}

	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)))
	if err != nil {
		return nil, err
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", getEnv("ENVIRONMENT", "development")),
		)),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SamplingRate)),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator (for context propagation)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

// TracingMiddleware adds distributed tracing to HTTP requests
func TracingMiddleware(serviceName string) func(http.Handler) http.Handler {
	tracer := otel.Tracer(serviceName)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from incoming request
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Start span
			ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPMethodKey.String(r.Method),
					semconv.HTTPTargetKey.String(r.URL.Path),
					semconv.HTTPURLKey.String(r.URL.String()),
					semconv.HTTPUserAgentKey.String(r.Header.Get("User-Agent")),
					semconv.HTTPClientIPKey.String(r.RemoteAddr),
				),
			)
			defer span.End()

			// Wrap response writer to capture status code
			wrapped := &statusRecorder{ResponseWriter: w, statusCode: 200}

			// Pass traced context to next handler
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// Record response status
			span.SetAttributes(semconv.HTTPStatusCodeKey.Int(wrapped.statusCode))

			if wrapped.statusCode >= 400 {
				span.RecordError(nil)
			}
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Flush implements http.Flusher to support SSE streaming
func (r *statusRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// StartSpan starts a new span (for manual instrumentation)
func StartSpan(ctx context.Context, operationName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("gateway")
	return tracer.Start(ctx, operationName)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, event string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(event, trace.WithAttributes(attrs...))
}

// SetSpanError marks the span as errored
func SetSpanError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}
