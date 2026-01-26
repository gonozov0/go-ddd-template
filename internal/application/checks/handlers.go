package checks

import protobuf "go-ddd-template/generated/server"

type CheckHandlers struct {
	protobuf.UnimplementedCheckServiceServer
}

func SetupHandlers() CheckHandlers {
	//nolint:exhaustruct // partial initialization is intentional
	return CheckHandlers{}
}
