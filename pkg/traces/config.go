package traces

import (
	"fmt"

	"go-ddd-template/pkg/envutils"
)

type Config struct {
	RequestIDHeader string
	ServiceName     string
	// TraceRatio float in [0;1] representing spans sampling ratio.
	// 0 means no spans will be exported,
	// 1 means all spans will be exported.
	TraceRatio        float64
	ExporterType      ExporterType
	CollectorEndpoint string
	// TraceProviderLabels - special labels for trace provider.
	// Can be set in order to mark traces with additional information.
	// By default, all traces are distinguished by service-name.
	TraceProviderLabels map[string]string
}

func (c *Config) Validate() error {
	if c.RequestIDHeader == "" {
		return fmt.Errorf("traces requires REQUEST_ID_HEADER to be set")
	}

	if c.TraceRatio < 0 || c.TraceRatio > 1 {
		return fmt.Errorf("traces requires TRACE_RATIO to be in range [0;1]")
	}

	if c.ExporterType == "" {
		return fmt.Errorf("traces requires EXPORTER_TYPE to be set")
	}

	if c.ExporterType == OtelCollector && c.CollectorEndpoint == "" {
		return fmt.Errorf("traces requires TRACING_OTLP_GRPC to be set")
	}

	return nil
}

func LoadConfig() (Config, error) {
	cfg := Config{
		RequestIDHeader:   envutils.GetEnv("REQUEST_ID_HEADER", "x-req-id"),
		ServiceName:       envutils.GetEnv("DEPLOY_UNIT_ID", ""),
		TraceRatio:        envutils.GetEnvFloat("TRACE_RATIO", 1),
		ExporterType:      ExporterType(envutils.GetEnv("EXPORTER_TYPE", string(Disable))),
		CollectorEndpoint: envutils.GetEnv("TRACING_OTLP_GRPC", ""),
		TraceProviderLabels: map[string]string{
			"environment": envutils.GetEnv("ENV_TYPE", "local"),
		},
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}
