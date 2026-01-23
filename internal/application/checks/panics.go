package checks

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h CheckHandlers) CheckPanics(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	panic("test panic")
}
