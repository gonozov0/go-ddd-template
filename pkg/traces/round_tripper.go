package traces

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RoundTripper struct {
	Proxied http.RoundTripper
	scope   string
}

func (rt RoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	ctx := req.Context()

	ctx, span := CreateSpan(
		ctx,
		rt.scope,
		fmt.Sprintf("%s %s", rt.scope, req.URL.Path),
		trace.SpanKindClient,
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", req.Method),
		attribute.String("http.url", req.URL.String()),
		attribute.String("http.host", req.Host),
	)

	req = req.Clone(ctx)

	resp, err := rt.Proxied.RoundTrip(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	span.SetStatus(codes.Ok, "")

	return resp, nil
}
