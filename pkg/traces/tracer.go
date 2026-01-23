package traces

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"fmt"

	loggerutils "go-ddd-template/pkg/logger/utils"
)

// InitDefaultTracer - inits default global tracer.
// Global tracer is accessed by platform's middleware.
// If you want to add custom spans in your application,
// create your own tracer by calling:
// otel.Tracer({your_package_name}, {options})
// and create/inherit spans with it.
func InitDefaultTracer(
	ctx context.Context,
	cfg Config,
) (func(context.Context) error, error) {
	provider, closeFunc, err := NewTracerProvider(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new tracer provider: %w", err)
	}

	otel.SetTracerProvider(provider)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		slog.WarnContext(ctx, "[OPEN-TELEMETRY ERROR]:", loggerutils.ErrAttr(err))
	}))
	// Setup propagator for support of cross-service traces.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return closeFunc, nil
}
