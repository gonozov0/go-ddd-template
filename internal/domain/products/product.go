package products

import (
	"github.com/google/uuid"

	"go-ddd-template/pkg/slices"

	"errors"

	"go-ddd-template/internal/domain/shared/valueobjects"
)

var (
	ErrProductNotFound         = errors.New("product not found")
	ErrInvalidProduct          = errors.New("invalid product")
	ErrProductAlreadyPublished = errors.New("product already published")
)

type Product struct {
	id            valueobjects.ProductID
	name          valueobjects.ProductName
	price         valueobjects.ProductPrice
	status        valueobjects.ProductStatus
	imageFilename valueobjects.ImageFilename
	imageURL      valueobjects.ImageURL
}

func NewProduct(
	id valueobjects.ProductID,
	name valueobjects.ProductName,
	price valueobjects.ProductPrice,
	status valueobjects.ProductStatus,
	imageFilename valueobjects.ImageFilename,
	imageURL valueobjects.ImageURL,
) *Product {
	return &Product{
		id:            id,
		name:          name,
		price:         price,
		status:        status,
		imageFilename: imageFilename,
		imageURL:      imageURL,
	}
}

func CreateProduct(name valueobjects.ProductName, price valueobjects.ProductPrice) *Product {
	return NewProduct(
		valueobjects.ProductID(uuid.New()),
		name,
		price,
		valueobjects.ProductStatusInit,
		valueobjects.EmptyImageFilename,
		valueobjects.EmptyImageURL,
	)
}

func (p *Product) GetID() valueobjects.ProductID {
	return p.id
}

func (p *Product) GetName() valueobjects.ProductName {
	return p.name
}

func (p *Product) GetPrice() valueobjects.ProductPrice {
	return p.price
}

func (p *Product) GetStatus() valueobjects.ProductStatus {
	return p.status
}

func (p *Product) Publish() error {
	if p.status == valueobjects.ProductStatusPublished {
		return ErrProductAlreadyPublished
	}

	p.status = valueobjects.ProductStatusPublished

	return nil
}

func (p *Product) GetImageFilename() valueobjects.ImageFilename {
	return p.imageFilename
}

func (p *Product) SetImageFilename(imageFilename valueobjects.ImageFilename) {
	p.imageFilename = imageFilename
	p.imageURL = valueobjects.EmptyImageURL
}

func (p *Product) GetImageURL() valueobjects.ImageURL {
	return p.imageURL
}

func (p *Product) SetImageURL(imageURL valueobjects.ImageURL) {
	p.imageURL = imageURL
}

type Products []Product

func (ps Products) IDs() valueobjects.ProductIDs {
	return slices.Map(ps, func(p Product) valueobjects.ProductID { return p.GetID() })
}
