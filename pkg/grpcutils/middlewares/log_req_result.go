package middlewares

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	loggerutils "go-ddd-template/pkg/logger/utils"
)

func LogRequestResult() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)

		if !isLoggerRequired(info.FullMethod) {
			//nolint:descriptiveerrors
			return resp, err
		}

		if err == nil {
			slog.InfoContext(ctx, info.FullMethod, loggerutils.Attr("code", codes.OK))
			//nolint:descriptiveerrors
			return resp, err
		}

		status, ok := status.FromError(err)
		if !ok {
			slog.WarnContext(ctx, info.FullMethod, loggerutils.ErrAttr(err))
			//nolint:descriptiveerrors
			return resp, err
		}

		switch status.Code() {
		case codes.Internal,
			codes.Unknown,
			codes.DeadlineExceeded,
			codes.Unimplemented,
			codes.DataLoss,
			codes.Unavailable:
			//nolint:descriptiveerrors
			return resp, err
		default:
			slog.InfoContext(
				ctx,
				info.FullMethod,
				loggerutils.Attr("code", status.Code()),
				loggerutils.Attr("message", status.Message()),
				loggerutils.ErrAttr(err),
			)

			//nolint:descriptiveerrors
			return resp, err
		}
	}
}

func isLoggerRequired(fullMethod string) bool {
	return !strings.HasPrefix(fullMethod, "/grpc.health.v1.Health/")
}
