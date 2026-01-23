package products

import (
	"context"

	"github.com/google/uuid"

	"fmt"

	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
)

type productMDB struct {
	ID            uuid.UUID
	Name          string
	Price         float64
	Status        string
	ImageFilename string
	ImageURL      string
}

type InMemoryRepo struct {
	products map[valueobjects.ProductID]productMDB
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		products: make(map[valueobjects.ProductID]productMDB),
	}
}

func (r *InMemoryRepo) CreateProducts(_ context.Context, createFn func() (domain.Products, error)) error {
	ps, err := createFn()
	if err != nil {
		return fmt.Errorf("failed to create products: %w", err)
	}

	for _, p := range ps {
		r.products[p.GetID()] = parseProductDomainToMDB(&p)
	}

	return nil
}

func (r *InMemoryRepo) GetProducts(_ context.Context, ids valueobjects.ProductIDs) ([]domain.Product, error) {
	ps := make([]domain.Product, 0, len(ids))

	for _, id := range ids {
		p, ok := r.products[id]
		if !ok {
			return nil, fmt.Errorf("%w: id %s", domain.ErrProductNotFound, id)
		}

		product, err := parseProductMDBToDomain(p)
		if err != nil {
			return nil, fmt.Errorf("failed to parse product mdb: %w", err)
		}

		ps = append(ps, *product)
	}

	return ps, nil
}

func (r *InMemoryRepo) DeleteProducts(
	_ context.Context,
	ids valueobjects.ProductIDs,
	deleteFn func([]domain.Product) error,
) error {
	productsList := make([]domain.Product, 0, len(ids))

	for _, id := range ids {
		if _, ok := r.products[id]; !ok {
			return fmt.Errorf("%w: id %s", domain.ErrProductNotFound, id)
		}
	}

	if err := deleteFn(productsList); err != nil {
		return fmt.Errorf("failed to delete products: %w", err)
	}

	for _, id := range ids {
		delete(r.products, id)
	}

	return nil
}

func (r *InMemoryRepo) UpdateProduct(
	ctx context.Context,
	id valueobjects.ProductID,
	updateFn func(*domain.Product) error,
) error {
	productMDB, ok := r.products[id]
	if !ok {
		return fmt.Errorf("%w: id %s", domain.ErrProductNotFound, id)
	}

	product, err := parseProductMDBToDomain(productMDB)
	if err != nil {
		return fmt.Errorf("failed to parse product mdb: %w", err)
	}

	if err := updateFn(product); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	r.products[id] = parseProductDomainToMDB(product)

	return nil
}

func parseProductDomainToMDB(product *domain.Product) productMDB {
	return productMDB{
		ID:            product.GetID().UUID(),
		Name:          product.GetName().String(),
		Price:         product.GetPrice().Float64(),
		Status:        product.GetStatus().String(),
		ImageFilename: product.GetImageFilename().String(),
		ImageURL:      product.GetImageURL().String(),
	}
}

func parseProductMDBToDomain(productMDB productMDB) (*domain.Product, error) {
	name, err := valueobjects.NewProductName(productMDB.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to init product name: %w", err)
	}

	price, err := valueobjects.NewProductPrice(productMDB.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to init product price: %w", err)
	}

	status, err := valueobjects.NewProductStatus(productMDB.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to init product status: %w", err)
	}

	imageFilename, err := valueobjects.NewImageFilename(productMDB.ImageFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to init product image filename: %w", err)
	}

	imageURL := valueobjects.NewImageURL(productMDB.ImageURL)

	return domain.NewProduct(valueobjects.ProductID(productMDB.ID), name, price, status, imageFilename, imageURL), nil
}
