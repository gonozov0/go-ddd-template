package orders

import (
	"context"
	"encoding/json"
	"fmt"

	hasql "golang.yandex/hasql/sqlx"

	domain "go-echo-template/internal/domain/orders"

	"github.com/google/uuid"
)

type itemDB struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}

type orderDB struct {
	ID     uuid.UUID          `db:"id"`
	UserID uuid.UUID          `db:"user_id"`
	Status domain.OrderStatus `db:"status"`
	Items  json.RawMessage    `db:"items"`
}

type productDB struct {
	ID        uuid.UUID `db:"id"`
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

func (r *PostgresRepo) SaveOrder(ctx context.Context, entity *domain.Order) error {
	db := r.cluster.Primary().DBx()

	items := make([]itemDB, 0, len(entity.Items()))
	for _, item := range entity.Items() {
		items = append(items, itemDB{
			ID:    item.ID(),
			Name:  item.Name(),
			Price: item.Price(),
		})
	}
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}
	order := orderDB{
		ID:     entity.ID(),
		UserID: entity.UserID(),
		Status: entity.Status(),
		Items:  itemsJSON,
	}

	query := `
		INSERT INTO orders (id, user_id, status, items)
		VALUES (:id, :user_id, :status, :items)
		ON CONFLICT (id) DO UPDATE SET
		user_id = excluded.user_id,
		status = excluded.status,
		items = excluded.items
	`
	_, err = db.NamedExecContext(ctx, query, order)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	db := r.cluster.StandbyPreferred().DBx()

	var order orderDB
	query := `SELECT id, user_id, status, items FROM orders WHERE id = $1`
	if err := db.GetContext(ctx, &order, query, id); err != nil {
		return nil, fmt.Errorf("failed to select order: %w", err)
	}

	var items []domain.Item
	if err := json.Unmarshal(order.Items, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	entity, err := domain.NewOrder(order.ID, order.UserID, order.Status, items)
	if err != nil {
		return nil, fmt.Errorf("failed to init order entity: %w", err)
	}

	return entity, nil
}

func (r *PostgresRepo) ReserveProducts(ctx context.Context, ids []uuid.UUID) error {
	db := r.cluster.Primary().DBx()

	query := `SELECT id, available FROM products WHERE id = ANY($1) FOR UPDATE`
	var products []productDB
	if err := db.SelectContext(ctx, &products, query, ids); err != nil {
		return fmt.Errorf("failed to select products: %w", err)
	}
	if len(products) != len(ids) {
		return fmt.Errorf("%w: selected %d products, expected %d", domain.ErrProductNotFound, len(products), len(ids))
	}

	for _, product := range products {
		if !product.Available {
			return fmt.Errorf("%w: id %s", domain.ErrProductAlreadyReserved, product.ID)
		}
	}

	query = `UPDATE products SET available = false WHERE id = ANY($1)`
	if _, err := db.ExecContext(ctx, query, ids); err != nil {
		return fmt.Errorf("failed to update products' availability: %w", err)
	}

	return nil
}
