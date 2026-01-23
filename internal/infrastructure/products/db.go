package products

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	domain "go-ddd-template/internal/domain/products"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/db/postgres"
	"go-ddd-template/pkg/sqs"
	"go-ddd-template/pkg/traces"
)

type (
	productDB struct {
		ID            uuid.UUID `db:"id"`
		Name          string    `db:"name"`
		Price         float64   `db:"price"`
		Status        string    `db:"status"`
		ImageFilename string    `db:"image_filename"`
		ImageURL      string    `db:"image_url"`
	}

	productInitedEvent struct {
		ID uuid.UUID `json:"id"`
	}
)

type QueueWriters struct {
	ProductInited *sqs.Writer
}

type PostgresRepo struct {
	cluster      *postgres.Cluster
	queueWriters QueueWriters
}

func NewDBRepo(cluster *postgres.Cluster, queueWriters QueueWriters) *PostgresRepo {
	return &PostgresRepo{
		cluster:      cluster,
		queueWriters: queueWriters,
	}
}

func (r *PostgresRepo) CreateProducts(ctx context.Context, createFn func() (domain.Products, error)) error {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get primary dbx: %w", err)
	}

	return postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		products, err := createFn()
		if err != nil {
			return fmt.Errorf("failed to create products: %w", err)
		}

		if err := r.createProductsDB(ctx, tx, products); err != nil {
			return fmt.Errorf("failed to create products in db: %w", err)
		}

		for _, product := range products {
			if err = r.sendProductInitedEvent(ctx, &product); err != nil {
				return fmt.Errorf("failed to send product inited event: %w", err)
			}
		}

		return nil
	})
}

func (r *PostgresRepo) CreateProductsInTests(ctx context.Context, products domain.Products) error {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get primary dbx: %w", err)
	}

	return postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		if err := r.createProductsDB(ctx, tx, products); err != nil {
			return fmt.Errorf("failed to save products in db: %w", err)
		}

		return nil
	})
}

func (r *PostgresRepo) sendProductInitedEvent(ctx context.Context, product *domain.Product) error {
	event := productInitedEvent{
		ID: product.GetID().UUID(),
	}

	message, err := json.Marshal(&event)
	if err != nil {
		return fmt.Errorf("failed to marshal product inited event: %w", err)
	}

	if err = r.queueWriters.ProductInited.SendMessage(ctx, string(message)); err != nil {
		return fmt.Errorf("failed to send product inited event: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetProducts(ctx context.Context, ids valueobjects.ProductIDs) ([]domain.Product, error) {
	db, err := r.cluster.StandbyPreferredDBx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get standby preferred dbx: %w", err)
	}

	productsDB, err := r.getProductsDB(ctx, db, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get products from db: %w", err)
	}

	products := make([]domain.Product, 0, len(productsDB))

	for _, productDB := range productsDB {
		product, err := parseProductDBToDomain(productDB)
		if err != nil {
			return nil, fmt.Errorf("failed to init product entity: %w", err)
		}

		products = append(products, *product)
	}

	return products, nil
}

func (r *PostgresRepo) DeleteProducts(
	ctx context.Context,
	ids valueobjects.ProductIDs,
	deleteFn func([]domain.Product) error,
) error {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get primary dbx: %w", err)
	}

	if err := postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		products, err := r.getProductsForUpdate(ctx, tx, ids)
		if err != nil {
			return fmt.Errorf("failed to get products for update: %w", err)
		}

		if err = deleteFn(products); err != nil {
			return fmt.Errorf("failed to delete products: %w", err)
		}

		if err := r.deleteProductsDB(ctx, tx, ids); err != nil {
			return fmt.Errorf("failed to delete products from db: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run delete products transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepo) UpdateProduct(
	ctx context.Context,
	id valueobjects.ProductID,
	updateFn func(*domain.Product) error,
) error {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get primary dbx: %w", err)
	}

	return postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		product, err := r.getProductForUpdate(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get product for update: %w", err)
		}

		if err := updateFn(product); err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}

		if err := r.updateProduct(ctx, tx, product); err != nil {
			return fmt.Errorf("failed to update product in db: %w", err)
		}

		return nil
	})
}

func parseProductDBToDomain(productDB productDB) (*domain.Product, error) {
	name, err := valueobjects.NewProductName(productDB.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to init product name: %w", err)
	}

	price, err := valueobjects.NewProductPrice(productDB.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to init product price: %w", err)
	}

	status, err := valueobjects.NewProductStatus(productDB.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to init product status: %w", err)
	}

	imageFilename, err := valueobjects.NewImageFilename(productDB.ImageFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to init product image filename: %w", err)
	}

	imageURL := valueobjects.NewImageURL(productDB.ImageURL)

	return domain.NewProduct(valueobjects.ProductID(productDB.ID), name, price, status, imageFilename, imageURL), nil
}

func parseProductsDBToDomain(productsDB []productDB) ([]domain.Product, error) {
	products := make([]domain.Product, 0, len(productsDB))

	for _, productDB := range productsDB {
		product, err := parseProductDBToDomain(productDB)
		if err != nil {
			return nil, fmt.Errorf("failed to init product entity: %w", err)
		}

		products = append(products, *product)
	}

	return products, nil
}

func parseProductDomainToDB(product *domain.Product) productDB {
	return productDB{
		ID:            product.GetID().UUID(),
		Name:          product.GetName().String(),
		Price:         product.GetPrice().Float64(),
		Status:        product.GetStatus().String(),
		ImageFilename: product.GetImageFilename().String(),
		ImageURL:      product.GetImageURL().String(),
	}
}

func (r *PostgresRepo) createProductsDB(ctx context.Context, tx *sqlx.Tx, products []domain.Product) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Create Products", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	stmt, err := tx.PreparexContext(ctx, `
		INSERT INTO products (id, name, price, status, image_filename, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, product := range products {
		productDB := parseProductDomainToDB(&product)

		_, err = stmt.ExecContext(
			ctx,
			productDB.ID,
			productDB.Name,
			productDB.Price,
			productDB.Status,
			productDB.ImageFilename,
			productDB.ImageURL,
		)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	return nil
}

func (r *PostgresRepo) getProductsDB(
	ctx context.Context,
	db *sqlx.DB,
	ids valueobjects.ProductIDs,
) (productsDB []productDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get Products", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT id, name, price, status, image_filename, image_url FROM products WHERE id = ANY($1)`
	if err := db.SelectContext(ctx, &productsDB, query, ids.Strings()); err != nil {
		return nil, fmt.Errorf("failed to select products: %w", err)
	}

	return productsDB, nil
}

func (r *PostgresRepo) deleteProductsDB(ctx context.Context, tx *sqlx.Tx, ids valueobjects.ProductIDs) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Delete Products", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := "DELETE FROM products WHERE id = ANY($1)"

	if _, err := tx.ExecContext(ctx, query, ids.Strings()); err != nil {
		return fmt.Errorf("failed executing delete query: %w", err)
	}

	return nil
}

func (r *PostgresRepo) getProductsForUpdateDB(
	ctx context.Context,
	tx *sqlx.Tx,
	ids valueobjects.ProductIDs,
) (productsDB []productDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get Products For Update", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT id, name, price, status, image_filename, image_url FROM products WHERE id = ANY($1) FOR UPDATE`
	if err := tx.SelectContext(ctx, &productsDB, query, ids.Strings()); err != nil {
		return nil, fmt.Errorf("failed to select products: %w", err)
	}

	return productsDB, nil
}

func (r *PostgresRepo) getProductsForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	ids valueobjects.ProductIDs,
) ([]domain.Product, error) {
	productsDB, err := r.getProductsForUpdateDB(ctx, tx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get products db: %w", err)
	}

	products, err := parseProductsDBToDomain(productsDB)
	if err != nil {
		return nil, fmt.Errorf("failed to parse products db: %w", err)
	}

	return products, nil
}

func (r *PostgresRepo) getProductDBForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	id valueobjects.ProductID,
) (product productDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get Product For Update", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := "SELECT id, name, price, status, image_filename, image_url FROM products WHERE id = $1 FOR UPDATE"

	if err := tx.GetContext(ctx, &product, query, id.UUID()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return productDB{}, domain.ErrProductNotFound
		}

		return productDB{}, fmt.Errorf("failed executing select query: %w", err)
	}

	return product, nil
}

func (r *PostgresRepo) getProductForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	id valueobjects.ProductID,
) (*domain.Product, error) {
	productDB, err := r.getProductDBForUpdate(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product db: %w", err)
	}

	product, err := parseProductDBToDomain(productDB)
	if err != nil {
		return nil, fmt.Errorf("failed to parse product db: %w", err)
	}

	return product, nil
}

func (r *PostgresRepo) updateProduct(ctx context.Context, tx *sqlx.Tx, product *domain.Product) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Update Product", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	productDB := parseProductDomainToDB(product)

	query := "UPDATE products SET name = :name, price = :price, status = :status, image_filename = :image_filename, image_url = :image_url WHERE id = :id"

	if _, err := tx.NamedExecContext(ctx, query, productDB); err != nil {
		return fmt.Errorf("failed to executing update query: %w", err)
	}

	return nil
}
