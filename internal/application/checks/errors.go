package checks

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

	"errors"
	"fmt"

	"go-ddd-template/pkg/grpcutils"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

func (h CheckHandlers) CheckErrors(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := errors.New("test error")

	slog.ErrorContext(
		ctx,
		"test error",
		loggerutils.ErrAttr(fmt.Errorf("wrapped test error: %w", err)),
	)

	return nil, grpcutils.ErrInternalError
}
