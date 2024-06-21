package generated

//go:generate oapi-codegen -generate types -package openapi -o ./openapi/types.go ../api/openapi.yaml
//go:generate oapi-codegen -generate server -package openapi -o ./openapi/server.go ../api/openapi.yaml
//go:generate oapi-codegen -generate client -package openapi -o ./openapi/client.go ../api/openapi.yaml
//go:generate oapi-codegen -generate spec -package openapi -o ./openapi/spec.go ../api/openapi.yaml
