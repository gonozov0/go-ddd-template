package traces

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
	"go.opentelemetry.io/otel/trace/noop"

	"fmt"
)

const (
	traceIDLogKey    = "trace_id"
	requestIDSpanKey = "request_id"
)

// NewTracerProvider - creates tracer provider from config.
// Should be used directly if it is important to create specific tracer
// and default global tracer is not suitable.
// Otherwise, it is recommended to use InitDefaultTracer.
func NewTracerProvider(
	ctx context.Context,
	cfg Config,
) (trace.TracerProvider, func(context.Context) error, error) {
	var (
		exporter tracesdk.SpanExporter
		err      error
	)

	switch cfg.ExporterType {
	case Disable:
		return noop.NewTracerProvider(), func(context.Context) error { return nil }, nil
	case Stdout:
		exporter, err = NewStdoutExporter()
	case OtelCollector:
		exporter, err = NewOpenTelemetryExporter(ctx, cfg)
	default:
		return nil, nil, fmt.Errorf("invalid exporter type: %v", cfg.ExporterType)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("exporter: %w", err)
	}

	exporter = &exporterWithMetrics{
		exporter: exporter,
	}

	resource, err := NewDefaultResource(cfg.ServiceName, cfg.TraceProviderLabels)
	if err != nil {
		return nil, nil, fmt.Errorf("new default resource: %w", err)
	}

	ratioSampler := tracesdk.TraceIDRatioBased(cfg.TraceRatio)
	sampler := tracesdk.ParentBased(ratioSampler, tracesdk.WithRemoteParentNotSampled(ratioSampler))

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource),
		tracesdk.WithSampler(sampler),
	)

	//nolint:exhaustivestruct
	return &tracerProviderWrapper{
		provider: tracerProvider,
	}, tracerProvider.Shutdown, nil
}

type tracerProviderWrapper struct {
	embedded.TracerProvider
	provider trace.TracerProvider
}

func (p *tracerProviderWrapper) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	tracer := p.provider.Tracer(name, options...)
	//nolint:exhaustivestruct
	return &tracerWrapper{tracer: tracer}
}

type tracerWrapper struct {
	embedded.Tracer
	tracer trace.Tracer
}

func (t *tracerWrapper) Start(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(ctx, spanName, opts...)

	traceID := span.SpanContext().TraceID()
	if traceID.IsValid() {
		ctx = context.WithValue(ctx, traceIDLogKey, traceID.String())
	}

	requestID := GetRequestID(ctx)
	if requestID != "" {
		span.SetAttributes(attribute.Key(requestIDSpanKey).String(requestID))
	}

	return ctx, span
}
