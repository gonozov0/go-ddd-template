package traces

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"

	"fmt"
)

// ExporterType - span exporter type.
type ExporterType string

const (
	// Disable - no export, spans wil be dropped.
	Disable ExporterType = "disable"
	// Stdout - print spans to stdout for debug.
	Stdout ExporterType = "stdout"
	// OtelCollector - spans will be sent to the specified
	// open-telemetry collector endpoint.
	OtelCollector ExporterType = "otel"
)

// NewStdoutExporter - build stdout exporter.
func NewStdoutExporter() (trace.SpanExporter, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	return exporter, err
}

// NewOpenTelemetryExporter - build open-telemetry exporter.
func NewOpenTelemetryExporter(ctx context.Context, cfg Config) (trace.SpanExporter, error) {
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.CollectorEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("new exporter: %w", err)
	}

	return exporter, nil
}

type exporterWithMetrics struct {
	exporter trace.SpanExporter
}

func (e *exporterWithMetrics) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	return e.exporter.ExportSpans(ctx, spans)
}

func (e *exporterWithMetrics) Shutdown(ctx context.Context) error {
	return e.exporter.Shutdown(ctx)
}
