package application

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	"go-ddd-template/generated"
	pb "go-ddd-template/generated/server"
	"go-ddd-template/pkg/auth"
	"go-ddd-template/pkg/traces"
)

const (
	readHeaderTimeout = 2 * time.Second
)

func SetupHTTPServer(
	conn *grpc.ClientConn,
	authCfg auth.Config,
	tracesCfg traces.Config,
) (*http.Server, error) {
	headers := []string{
		tracesCfg.RequestIDHeader,
		authCfg.UserHeader,
	}

	gatewayMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithMetadata(headerTransporter(headers...)),
	)
	if err := pb.RegisterOrderServiceHandler(context.Background(), gatewayMux, conn); err != nil {
		return nil, fmt.Errorf("failed to register order service handler: %w", err)
	}

	if err := pb.RegisterUserServiceHandler(context.Background(), gatewayMux, conn); err != nil {
		return nil, fmt.Errorf("failed to register user service handler: %w", err)
	}

	if err := pb.RegisterProductServiceHandler(context.Background(), gatewayMux, conn); err != nil {
		return nil, fmt.Errorf("failed to register product service handler: %w", err)
	}

	if err := pb.RegisterCheckServiceHandler(context.Background(), gatewayMux, conn); err != nil {
		return nil, fmt.Errorf("failed to register ping service handler: %w", err)
	}

	if err := pb.RegisterDeliveryServiceHandler(context.Background(), gatewayMux, conn); err != nil {
		return nil, fmt.Errorf("failed to register delivery service handler: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/ping", getHealthCheckHandler(conn))
	mux.Handle("/swagger/", getSwaggerHandler())

	mux.Handle("/", gatewayMux)

	return &http.Server{
		Handler:           removeSlashMiddleware(mux),
		ReadHeaderTimeout: readHeaderTimeout,
	}, nil
}

func getSwaggerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger.json" {
			w.Header().Set("Content-Type", "application/json")
			w.Write(generated.SwaggerJSON)

			return
		}

		httpSwagger.Handler(
			httpSwagger.URL("/swagger.json"),
		).ServeHTTP(w, r)
	})
}

func getHealthCheckHandler(conn *grpc.ClientConn) http.Handler {
	healthClient := grpc_health_v1.NewHealthClient(conn)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
			Service: "", // check all services
		})
		if err != nil {
			http.Error(w, "gRPC health check failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			http.Error(w, "Service not healthy", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// headerTransporter returns a function that transposes headers from the incoming HTTP request to gRPC metadata.
func headerTransporter(headers ...string) func(ctx context.Context, req *http.Request) metadata.MD {
	return func(ctx context.Context, req *http.Request) metadata.MD {
		existingMD, _ := metadata.FromIncomingContext(ctx)
		newMD := metadata.MD{}

		for _, header := range headers {
			value := req.Header.Get(header)
			if value != "" {
				newMD.Append(header, value)
			}
		}

		return metadata.Join(existingMD, newMD)
	}
}

func removeSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") && r.URL.Path != "/" && !strings.HasPrefix(r.URL.Path, "/swagger") {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}
