package orders

import (
	"go-echo-ddd-template/internal/domain/orders"
	"go-echo-ddd-template/internal/domain/products"

	"github.com/google/uuid"
)

type Item struct {
	ID       uuid.UUID
	Quantity int
}

func (s *OrderService) CreateAndPay(userID uuid.UUID, items []Item) (*orders.Order, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}

	ps, err := s.productRepo.GetProductsForUpdate(getItemIDs(items))
	if err != nil {
		return nil, err
	}

	orderItems, err := makeOrderItems(items, ps)
	if err != nil {
		return nil, err
	}
	order, err := orders.CreateOrder(userID, orderItems)
	if err != nil {
		return nil, err
	}
	if err = s.orderRepo.SaveOrder(*order); err != nil {
		return nil, err
	}

	if err = s.productRepo.SaveProducts(ps); err != nil {
		return nil, err
	}

	if err = order.Pay(); err != nil {
		return nil, err
	}
	if err = s.orderRepo.SaveOrder(*order); err != nil {
		return nil, err
	}

	invoice, err := order.MakeInvoice()
	if err != nil {
		return nil, err
	}
	if err = user.SendToEmail(invoice); err != nil {
		return nil, err
	}

	return order, nil
}

func getItemIDs(items []Item) []uuid.UUID {
	itemIDs := make([]uuid.UUID, 0, len(items))
	for _, i := range items {
		itemIDs = append(itemIDs, i.ID)
	}
	return itemIDs
}

func makeOrderItems(items []Item, ps []products.Product) ([]orders.Item, error) {
	orderItems := make([]orders.Item, 0, len(items))
	for i, product := range ps {
		item := items[i]
		if err := product.ReduceQuantity(item.Quantity); err != nil {
			return nil, err
		}
		orderItem, err := orders.NewItem(item.ID, product.Name(), product.Price(), item.Quantity)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, *orderItem)
	}
	return orderItems, nil
}
