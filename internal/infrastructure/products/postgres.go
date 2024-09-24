package products

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	hasql "golang.yandex/hasql/sqlx"

	domain "go-echo-template/internal/domain/products"
	"go-echo-template/pkg/postgres"

	"github.com/google/uuid"
)

type productsDB struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Price float64   `db:"price"`
}

type PostgresRepo struct {
	cluster *hasql.Cluster
}

func NewPostgresRepo(cluster *hasql.Cluster) *PostgresRepo {
	return &PostgresRepo{
		cluster: cluster,
	}
}

func (r *PostgresRepo) CreateProducts(ctx context.Context, createFn func() ([]domain.Product, error)) error {
	db := r.cluster.Primary().DBx()
	return postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		ps, err := createFn()
		if err != nil {
			return fmt.Errorf("failed to create products: %w", err)
		}

		//nolint:sqlclosecheck // Linter bug because of RunInTx, it's closing the statement in defer
		stmt, err := tx.PreparexContext(ctx, `
			INSERT INTO products (id, name, price)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE SET
			name = excluded.name,
			price = excluded.price
		`)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, p := range ps {
			_, err = stmt.ExecContext(ctx, p.ID(), p.Name(), p.Price())
			if err != nil {
				return fmt.Errorf("failed to execute statement: %w", err)
			}
		}

		return nil
	})
}

func (r *PostgresRepo) GetProducts(ctx context.Context, id []uuid.UUID) ([]domain.Product, error) {
	db := r.cluster.StandbyPreferred().DBx()
	var productsDB []productsDB
	query := `SELECT id, name, price FROM products WHERE id = ANY($1) FOR UPDATE`
	if err := db.SelectContext(ctx, &productsDB, query, id); err != nil {
		return nil, fmt.Errorf("failed to select products: %w", err)
	}

	products := make([]domain.Product, 0, len(productsDB))
	for _, p := range productsDB {
		product, err := domain.NewProduct(p.ID, p.Name, p.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to init product entity: %w", err)
		}
		products = append(products, *product)
	}

	return products, nil
}
