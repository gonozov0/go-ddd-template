package products

import (
	"context"
	"fmt"

	hasql "golang.yandex/hasql/sqlx"

	domain "go-echo-template/internal/domain/products"

	"github.com/google/uuid"
)

type productsDB struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Price     float64   `db:"price"`
	Available bool      `db:"available"`
}

type PostgresRepo struct {
	cluster *hasql.Cluster
}

func NewPostgresRepo(cluster *hasql.Cluster) *PostgresRepo {
	return &PostgresRepo{
		cluster: cluster,
	}
}

func (r *PostgresRepo) SaveProducts(ctx context.Context, ps []domain.Product) error {
	db := r.cluster.Primary().DBx()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // Handling rollback error is not necessary

	stmt, err := tx.PreparexContext(ctx, `
        INSERT INTO products (id, name, price, available)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE SET
        name = excluded.name,
        price = excluded.price,
        available = excluded.available
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, p := range ps {
		_, err = stmt.ExecContext(ctx, p.ID(), p.Name(), p.Price(), p.Available())
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetProductsForUpdate(ctx context.Context, id []uuid.UUID) ([]domain.Product, error) {
	db := r.cluster.StandbyPreferred().DBx()
	var productsDB []productsDB
	query := `SELECT id, name, price, available FROM products WHERE id = ANY($1) FOR UPDATE`
	if err := db.SelectContext(ctx, &productsDB, query, id); err != nil {
		return nil, fmt.Errorf("failed to select products: %w", err)
	}

	products := make([]domain.Product, 0, len(productsDB))
	for _, p := range productsDB {
		product, err := domain.NewProduct(p.ID, p.Name, p.Price, p.Available)
		if err != nil {
			return nil, fmt.Errorf("failed to init product entity: %w", err)
		}
		products = append(products, *product)
	}

	return products, nil
}
