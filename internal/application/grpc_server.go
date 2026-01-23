package application

import (
	"fmt"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"go-ddd-template/pkg/auth"

	pb "go-ddd-template/generated/server"
	"go-ddd-template/internal/application/checks"
	deliveries "go-ddd-template/internal/application/deliveries/server"
	orders "go-ddd-template/internal/application/orders/server"
	products "go-ddd-template/internal/application/products/server"
	users "go-ddd-template/internal/application/users/server"
	productsservice "go-ddd-template/internal/service/products"
	"go-ddd-template/pkg/grpcutils/middlewares"
	"go-ddd-template/pkg/traces"
)

type gRPCServer struct {
	users.UserHandlers
	orders.OrderHandlers
	products.ProductHandlers
	checks.CheckHandlers
	deliveries.DeliveryHandlers
}

func SetupGRPCServer(
	authCfg auth.Config,
	tracesCfg traces.Config,
	userRepo UserRepository,
	orderRepo OrderRepository,
	productRepo ProductRepository,
	deliveryRepo DeliveryRepository,
	imageStorage productsservice.ImageStorage,
) (*grpc.Server, error) {
	interceptors, err := setupMiddlewares(authCfg, tracesCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to setup middlewares: %w", err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))

	server := gRPCServer{
		UserHandlers:     users.SetupHandlers(userRepo),
		OrderHandlers:    orders.SetupHandlers(orderRepo, userRepo),
		ProductHandlers:  products.SetupHandlers(productRepo, imageStorage),
		CheckHandlers:    checks.SetupHandlers(),
		DeliveryHandlers: deliveries.SetupHandlers(deliveryRepo),
	}

	pb.RegisterOrderServiceServer(s, server)
	pb.RegisterUserServiceServer(s, server)
	pb.RegisterProductServiceServer(s, server)
	pb.RegisterCheckServiceServer(s, server)
	pb.RegisterDeliveryServiceServer(s, server)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(s)

	return s, nil
}

func setupMiddlewares(
	authCfg auth.Config,
	tracesCfg traces.Config,
) ([]grpc.UnaryServerInterceptor, error) {
	interceptors := []grpc.UnaryServerInterceptor{
		middlewares.AddMethodAndRequestToContext,
		traces.NewTraceMiddleware(tracesCfg),
		grpcrecovery.UnaryServerInterceptor(
			grpcrecovery.WithRecoveryHandlerContext(
				middlewares.PanicHandler(),
			),
		),
	}

	authMw, err := auth.NewAuthMiddleware(authCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init auth middleware: %w", err)
	}

	interceptors = append(interceptors, authMw)

	interceptors = append(interceptors, middlewares.LogRequestResult())

	return interceptors, nil
}
