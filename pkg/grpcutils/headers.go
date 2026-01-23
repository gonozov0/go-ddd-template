package grpcutils

import (
	"google.golang.org/grpc/metadata"

	"fmt"
)

func GetSingleHeader(md metadata.MD, headerName string) (string, error) {
	headerValues := md.Get(headerName)
	if len(headerValues) == 0 {
		return "", nil
	}

	if len(headerValues) > 1 {
		return "", fmt.Errorf("multiple values found (%d) for header '%s'", len(headerValues), headerName)
	}

	return headerValues[0], nil
}
