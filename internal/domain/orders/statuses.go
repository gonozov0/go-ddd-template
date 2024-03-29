package orders

type OrderStatus string

const (
	OrderStatusCreated           OrderStatus = "created"
	OrderStatusPaymentProcessing OrderStatus = "payment_processing"
	OrderStatusPaid              OrderStatus = "paid"
	OrderStatusShipped           OrderStatus = "shipped"
	OrderStatusDelivered         OrderStatus = "delivered"
)
