package auth

import "strings"

const HealthCheckServicePrefix = "/grpc.health.v1.Health/"

func IsAuthRequired(fullMethod string) bool {
	return !strings.HasPrefix(fullMethod, HealthCheckServicePrefix)
}
