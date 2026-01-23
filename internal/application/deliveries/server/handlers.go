package server

import (
	pb "go-ddd-template/generated/server"
	service "go-ddd-template/internal/service/deliveries"
)

type DeliveryHandlers struct {
	pb.UnimplementedDeliveryServiceServer
	deliveryService service.DeliveryService
}

func SetupHandlers(deliveryRepo service.DeliveryRepository) DeliveryHandlers {
	deliveryService := service.NewDeliveryService(deliveryRepo)

	//nolint:exhaustivestruct
	return DeliveryHandlers{
		deliveryService: deliveryService,
	}
}
