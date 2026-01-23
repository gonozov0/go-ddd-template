package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	loggerutils "go-ddd-template/pkg/logger/utils"
)

type ctxKey string

const (
	requestStr ctxKey = "grpc_request_str"
	method     ctxKey = "grpc_method"
)

func AddMethodAndRequestToContext(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	reqStr, err := maskAndStringify(req)
	if err != nil {
		slog.WarnContext(ctx, "failed to mask and stringify request", loggerutils.ErrAttr(err))

		reqStr = "failed to mask and stringify request"
	}

	ctx = context.WithValue(ctx, requestStr, reqStr)
	ctx = context.WithValue(ctx, method, info.FullMethod)

	return handler(ctx, req)
}

func maskAndStringify(req any) (string, error) {
	protoMsg, ok := req.(proto.Message)
	if !ok {
		return fmt.Sprintf("%v", req), nil
	}

	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,  // использовать snake_case из proto файлов
		EmitUnpopulated: false, // не включать пустые поля
	}

	data, err := marshaler.Marshal(protoMsg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal protoMsg: %w", err)
	}

	var jsonMap map[string]any
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return "", fmt.Errorf("failed to unmarshal data: %w", err)
	}

	maskSensitiveFields(jsonMap)

	maskedData, err := json.Marshal(jsonMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal jsonMap: %w", err)
	}

	return string(maskedData), nil
}

func maskSensitiveFields(data map[string]any) {
	sensitiveFields := map[string]bool{
		"phone":     true,
		"full_name": true,
		"email":     true,
		"password":  true,
		"telegram":  true,
		"token":     true,
		"api_key":   true,
	}

	for key, value := range data {
		if sensitiveFields[strings.ToLower(key)] {
			if strVal, ok := value.(string); ok && strVal != "" {
				data[key] = maskString(key, strVal)
			}

			continue
		}

		switch v := value.(type) {
		case map[string]any:
			maskSensitiveFields(v)
		case []any:
			for _, item := range v {
				if subMap, ok := item.(map[string]any); ok {
					maskSensitiveFields(subMap)
				}
			}
		}
	}
}

func maskString(fieldName, value string) string {
	switch strings.ToLower(fieldName) {
	case "phone":
		if len(value) > 5 {
			return value[:3] + "***" + value[len(value)-2:]
		}

		return "***"

	case "email":
		if idx := strings.Index(value, "@"); idx > 0 {
			return value[:1] + "***@" + value[idx+1:]
		}

		return "***"

	default:
		return "***"
	}
}

func GetGRPCMethod(ctx context.Context) string {
	if method, ok := ctx.Value(method).(string); ok {
		return method
	}

	return ""
}

func GetGRPCRequestStr(ctx context.Context) string {
	if reqStr, ok := ctx.Value(requestStr).(string); ok {
		return reqStr
	}

	return ""
}
