package traces

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"fmt"
)

// NewDefaultResource - creates default open-telemetry resource.
// Resource represents unique entity which provides traces.
// Service name is required, other labels such as project name,
// datacenter, environment, etc. are optional.
func NewDefaultResource(serviceName string, additionalAttributes map[string]string) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
	}

	for attrName, attrVal := range additionalAttributes {
		attrs = append(attrs, attribute.String(attrName, attrVal))
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			attrs...,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource merge: %w", err)
	}

	return res, nil
}
