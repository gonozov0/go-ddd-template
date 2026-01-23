package orders

import (
	"context"

	"fmt"

	domain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

type InMemoryRepo struct {
	domain   map[valueobjects.OrderID]domain.Order
	products map[valueobjects.ProductID]domain.Product
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		domain:   make(map[valueobjects.OrderID]domain.Order),
		products: make(map[valueobjects.ProductID]domain.Product),
	}
}

func (r *InMemoryRepo) CreateOrder(
	_ context.Context,
	productIDs valueobjects.ProductIDs,
	createFn func([]domain.Product) (*domain.Order, error),
) (*domain.Order, error) {
	products, err := r.getProducts(productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	order, err := createFn(products)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	r.domain[order.GetID()] = *order
	r.updateProductsStatus(products)

	return order, nil
}

func (r *InMemoryRepo) getProducts(ids valueobjects.ProductIDs) ([]domain.Product, error) {
	products := make([]domain.Product, 0, len(ids))
	for _, id := range ids {
		p, ok := r.products[id]
		if !ok {
			return nil, domain.ErrProductNotFound
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *InMemoryRepo) updateProductsStatus(products []domain.Product) {
	for _, product := range products {
		r.products[product.GetID()] = product
	}
}

func (r *InMemoryRepo) GetOrder(_ context.Context, id valueobjects.OrderID) (*domain.Order, error) {
	order, ok := r.domain[id]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}

	return &order, nil
}

func (r *InMemoryRepo) DeleteOrder(_ context.Context, orderID valueobjects.OrderID) error {
	_, ok := r.domain[orderID]
	if !ok {
		return domain.ErrOrderNotFound
	}

	delete(r.domain, orderID)

	return nil
}

// CreateProducts used only for unit tests, another aggregate is responsible for creating products.
func (r *InMemoryRepo) CreateProducts(_ context.Context, createFn func() (domain.Products, error)) error {
	products, err := createFn()
	if err != nil {
		return fmt.Errorf("failed to create products: %w", err)
	}

	for _, p := range products {
		r.products[p.GetID()] = p
	}

	return nil
}

func (r *InMemoryRepo) ProcessOrders(
	ctx context.Context,
	updateOrders func([]*domain.Order) error,
) error {
	orders := make([]*domain.Order, 0, len(r.domain))

	for _, order := range r.domain {
		if order.GetStatus() == domain.OrderStatusCreated {
			orders = append(orders, &order)
		}
	}

	if err := updateOrders(orders); err != nil {
		return fmt.Errorf("failed to update orders: %w", err)
	}

	for _, order := range orders {
		r.domain[order.GetID()] = *order
	}

	return nil
}
