package deliveries

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	domain "go-ddd-template/internal/domain/deliveries"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/db/postgres"
	"go-ddd-template/pkg/traces"
)

type deliveryDB struct {
	ID        uuid.UUID `db:"id"`
	OrderID   uuid.UUID `db:"order_id"`
	CreatedAt time.Time `db:"created_at"`
}

type DBRepo struct {
	cluster *postgres.Cluster
}

func NewDBRepo(cluster *postgres.Cluster) *DBRepo {
	return &DBRepo{
		cluster: cluster,
	}
}

func (r *DBRepo) CreateDelivery(
	ctx context.Context,
	createFn func() (*domain.Delivery, error),
) (*domain.Delivery, error) {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary dbx: %w", err)
	}

	var delivery *domain.Delivery

	err = postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		delivery, err = createFn()
		if err != nil {
			return fmt.Errorf("failed to create delivery: %w", err)
		}

		if err := r.createDeliveryDB(ctx, tx, delivery); err != nil {
			return fmt.Errorf("failed to create delivery in db: %w", err)
		}

		return nil
	})

	return delivery, err
}

func (r *DBRepo) ListDeliveries(ctx context.Context) ([]domain.Delivery, error) {
	db, err := r.cluster.StandbyPreferredDBx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get standby preferred dbx: %w", err)
	}

	deliveries, err := r.getDeliveriesDB(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliveries from db: %w", err)
	}

	ds := make([]domain.Delivery, 0, len(deliveries))

	for _, item := range deliveries {
		i := parseDeliveryDBToDomain(item)
		ds = append(ds, *i)
	}

	return ds, nil
}

func (r *DBRepo) DeleteDelivery(ctx context.Context, deliveryID valueobjects.DeliveryID) error {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get primary dbx: %w", err)
	}

	if err := r.deleteDeliveryDB(ctx, db, deliveryID); err != nil {
		return fmt.Errorf("failed to delete delivery from db: %w", err)
	}

	return nil
}

func parseDeliveryDBToDomain(deliveryDB deliveryDB) *domain.Delivery {
	return domain.NewDelivery(
		valueobjects.DeliveryID(deliveryDB.ID),
		valueobjects.OrderID(deliveryDB.OrderID),
		valueobjects.NewTimestamp(deliveryDB.CreatedAt),
	)
}

func parseDeliveryDomainToDB(delivery *domain.Delivery) deliveryDB {
	return deliveryDB{
		ID:        delivery.GetID().UUID(),
		OrderID:   delivery.GetOrderID().UUID(),
		CreatedAt: delivery.GetCreatedAt().Time(),
	}
}

func (r *DBRepo) createDeliveryDB(ctx context.Context, tx *sqlx.Tx, delivery *domain.Delivery) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Create Delivery", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	deliveryDB := parseDeliveryDomainToDB(delivery)

	query := `
		INSERT INTO deliveries (id, order_id, created_at)
		VALUES (:id, :order_id, :created_at)
	`

	_, err = tx.NamedExecContext(ctx, query, deliveryDB)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *DBRepo) getDeliveriesDB(ctx context.Context, db *sqlx.DB) (deliveries []deliveryDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get Deliveries", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT id, order_id, created_at FROM deliveries`

	if err := db.SelectContext(ctx, &deliveries, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get deliveries: %w", err)
	}

	return deliveries, nil
}

func (r *DBRepo) deleteDeliveryDB(ctx context.Context, db *sqlx.DB, deliveryID valueobjects.DeliveryID) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Delete Delivery", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `DELETE FROM deliveries WHERE id = $1`

	if _, err := db.ExecContext(ctx, query, deliveryID.String()); err != nil {
		return fmt.Errorf("failed executing delete query: %w", err)
	}

	return nil
}
