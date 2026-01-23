package middlewares

import (
	"context"
	"log/slog"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func PanicHandler() grpcrecovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, err any) error {
		slog.ErrorContext(ctx, "runtime panic", loggerutils.ErrAttr(err))

		return grpcutils.ErrInternalError
	}
}
