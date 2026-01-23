package traces

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go-ddd-template/pkg/grpcutils"
)

type requestIDCtxKey struct{}

func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey{}, id)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDCtxKey{}).(string); ok {
		return id
	}

	return ""
}

type correlcationIDCtxKey struct{}

func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlcationIDCtxKey{}, id)
}

func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlcationIDCtxKey{}).(string); ok {
		return id
	}

	return ""
}

// NewTraceMiddleware creates a new middleware that adds request ID and correlation ID to the context.
func NewTraceMiddleware(cfg Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata in request")
		}

		if isTraceRequired(info.FullMethod) {
			requestID, err := grpcutils.GetSingleHeader(md, cfg.RequestIDHeader)
			if err != nil {
				return nil, fmt.Errorf("failed to get request ID: %w", err)
			}

			if requestID != "" {
				ctx = withRequestID(ctx, requestID)
			}

			corelationID := uuid.New().String()

			ctx = WithCorrelationID(ctx, corelationID)

			ctx, span := CreateSpan(
				ctx,
				ScopeTraceMiddleware,
				fmt.Sprintf("%s %s", ScopeTraceMiddleware, info.FullMethod),
				trace.SpanKindServer,
			)
			defer span.End()

			span.SetAttributes(
				attribute.String("rpc.method", info.FullMethod),
				attribute.String("request_id", requestID),
				attribute.String("correlation_id", corelationID),
			)

			resp, err = handler(ctx, req)
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}

			//nolint:descriptiveerrors
			return resp, err
		} else {
			//nolint:descriptiveerrors
			return handler(ctx, req)
		}
	}
}

func isTraceRequired(fullMethod string) bool {
	return !strings.HasPrefix(fullMethod, "/grpc.health.v1.Health/")
}
