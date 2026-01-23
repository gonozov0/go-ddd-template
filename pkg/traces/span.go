package traces

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	ScopeTraceMiddleware = "middleware"
	ScopePostgres        = "postgres"
	ScopeConsumer        = "consumer"
)

type (
	spanIDCtxKey  struct{}
	traceIDCtxKey struct{}
)

func withSpanID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, spanIDCtxKey{}, id)
}

func withTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDCtxKey{}, id)
}

func GetSpanID(ctx context.Context) string {
	if id, ok := ctx.Value(spanIDCtxKey{}).(string); ok {
		return id
	}

	return ""
}

func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDCtxKey{}).(string); ok {
		return id
	}

	return ""
}

func CreateSpan(
	ctx context.Context,
	scope string,
	spanName string,
	spanKind trace.SpanKind,
) (context.Context, trace.Span) {
	provider := otel.GetTracerProvider()
	tracer := provider.Tracer(scope)

	ctx, span := tracer.Start(
		ctx,
		spanName,
		trace.WithSpanKind(spanKind),
	)

	if span.SpanContext().HasSpanID() {
		ctx = withSpanID(ctx, span.SpanContext().SpanID().String())
	}

	if span.SpanContext().HasTraceID() {
		ctx = withTraceID(ctx, span.SpanContext().TraceID().String())
	}

	return ctx, span
}

func SetSpanStatus(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
}
