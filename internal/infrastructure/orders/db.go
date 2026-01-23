package orders

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lib/pq"

	gofrsuuid "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicwriter"
	"go.opentelemetry.io/otel/trace"

	domain "go-ddd-template/internal/domain/orders"
	"go-ddd-template/internal/domain/shared/valueobjects"
	"go-ddd-template/pkg/db/postgres"
	"go-ddd-template/pkg/parseutils"
	"go-ddd-template/pkg/traces"
)

type orderDB struct {
	ID         uuid.UUID          `db:"id"`
	UserID     uuid.UUID          `db:"user_id"`
	Status     domain.OrderStatus `db:"status"`
	ProductIDs []gofrsuuid.UUID   `db:"product_ids,array"`
}

type productDB struct {
	ID     uuid.UUID `db:"id"`
	Name   string    `db:"name"`
	Price  float64   `db:"price"`
	Status string    `db:"status"`
}

type orderCreatedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
}

type QueueWriters struct {
	OrderCreated *topicwriter.Writer
}

type DBRepo struct {
	cluster      *postgres.Cluster
	queueWriters QueueWriters
}

func NewDBRepo(cluster *postgres.Cluster, queueWriters QueueWriters) *DBRepo {
	return &DBRepo{
		cluster:      cluster,
		queueWriters: queueWriters,
	}
}

func (r *DBRepo) CreateOrder(
	ctx context.Context,
	productIDs valueobjects.ProductIDs,
	createFn func([]domain.Product) (*domain.Order, error),
) (*domain.Order, error) {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary dbx: %w", err)
	}

	var order *domain.Order

	err = postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		products, err := r.getProductsForUpdate(ctx, tx, productIDs)
		if err != nil {
			return fmt.Errorf("failed to get products for update: %w", err)
		}

		order, err = createFn(products)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		if err := r.createOrderDB(ctx, tx, order); err != nil {
			return fmt.Errorf("failed to create order in db: %w", err)
		}

		if err := r.updateProducts(ctx, tx, products); err != nil {
			return fmt.Errorf("failed to update products status in db: %w", err)
		}

		return nil
	})

	return order, err
}

func (r *DBRepo) getProductsForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	ids valueobjects.ProductIDs,
) ([]domain.Product, error) {
	productDBs, err := r.getProductsForUpdateDB(ctx, tx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get products for update from db: %w", err)
	}

	products := make([]domain.Product, 0, len(productDBs))

	for _, productDB := range productDBs {
		product, err := parseProductDBToDomain(productDB)
		if err != nil {
			return nil, fmt.Errorf("failed to init product entity: %w", err)
		}

		products = append(products, *product)
	}

	return products, nil
}

func (r *DBRepo) createOrderDB(ctx context.Context, tx *sqlx.Tx, order *domain.Order) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Create Order", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	orderDB := parseOrderDomainToDB(order)

	query := `
		INSERT INTO orders (id, user_id, status, product_ids)
		VALUES ($1, $2, $3, $4)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare context: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, orderDB.ID, orderDB.UserID, orderDB.Status, pq.Array(orderDB.ProductIDs))
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}

// updateProducts может обновить только те поля, которые регулируются доменом Order
func (r *DBRepo) updateProducts(ctx context.Context, tx *sqlx.Tx, products []domain.Product) (err error) {
	ctx, span := traces.CreateSpan(
		ctx,
		traces.ScopePostgres,
		"Postgres Update Products",
		trace.SpanKindClient,
	)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	stmt, err := tx.PrepareNamedContext(ctx, `UPDATE products SET status = :status WHERE id = :id`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, product := range products {
		productDB := parseProductDomainToDB(&product)

		_, err := stmt.ExecContext(ctx, productDB)
		if err != nil {
			return fmt.Errorf("failed to update product status: %w", err)
		}
	}

	return nil
}

func (r *DBRepo) GetOrder(ctx context.Context, id valueobjects.OrderID) (*domain.Order, error) {
	db, err := r.cluster.StandbyPreferredDBx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get standby preferred dbx: %w", err)
	}

	order, err := r.getOrder(ctx, db, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

func (r *DBRepo) getOrder(ctx context.Context, db *sqlx.DB, id valueobjects.OrderID) (*domain.Order, error) {
	orderDB, err := r.getOrderDB(ctx, db, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from db: %w", err)
	}

	productIDs := make(valueobjects.ProductIDs, 0, len(orderDB.ProductIDs))

	for _, rawProductID := range parseutils.GofrsUUIDsToGoogleUUIDs(orderDB.ProductIDs) {
		productID, err := valueobjects.NewProductID(rawProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product id: %w", err)
		}

		productIDs = append(productIDs, productID)
	}

	products, err := r.getProducts(ctx, db, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return parseOrderDBToDomain(orderDB, products), nil
}

func (r *DBRepo) getCreatedOrdersForUpdate(ctx context.Context, tx *sqlx.Tx) ([]*domain.Order, error) {
	orderDBs, err := r.getCreateOrderDBsForUpdate(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from db: %w", err)
	}

	var productIDs valueobjects.ProductIDs

	var orderIDToProductIDs = make(map[uuid.UUID]valueobjects.ProductIDs)

	for _, orderDB := range orderDBs {
		for _, rawProductID := range parseutils.GofrsUUIDsToGoogleUUIDs(orderDB.ProductIDs) {
			productID, err := valueobjects.NewProductID(rawProductID)
			if err != nil {
				return nil, fmt.Errorf("failed to get product id: %w", err)
			}

			orderIDToProductIDs[orderDB.ID] = append(orderIDToProductIDs[orderDB.ID], productID)
			productIDs = append(productIDs, productID)
		}
	}

	allProducts, err := r.getProductsForUpdate(ctx, tx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	productIDToProduct := make(map[valueobjects.ProductID]domain.Product)
	for _, product := range allProducts {
		productIDToProduct[product.GetID()] = product
	}

	var orders []*domain.Order

	for _, orderDB := range orderDBs {
		var products []domain.Product
		for _, productID := range orderIDToProductIDs[orderDB.ID] {
			products = append(products, productIDToProduct[productID])
		}

		orders = append(orders, parseOrderDBToDomain(orderDB, products))
	}

	return orders, nil
}

func (r *DBRepo) getCreateOrderDBsForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
) (orderDBs []orderDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get Orders", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT id, user_id, status, product_ids FROM orders WHERE status = $1 FOR UPDATE`

	rows, err := tx.QueryContext(ctx, query, domain.OrderStatusCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}

	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			if err != nil {
				err = fmt.Errorf("failed to close rows: %w; %w", err, closeErr)
			} else {
				err = closeErr
			}
		}
	}()

	for rows.Next() {
		var orderDB orderDB
		if err = rows.Scan(
			&orderDB.ID,
			&orderDB.UserID,
			&orderDB.Status,
			pq.Array(&orderDB.ProductIDs),
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		orderDBs = append(orderDBs, orderDB)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to check rows error: %w", err)
	}

	return orderDBs, nil
}

func (r *DBRepo) getProducts(ctx context.Context, db *sqlx.DB, ids valueobjects.ProductIDs) ([]domain.Product, error) {
	productDBs, err := r.getProductsDB(ctx, db, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get products from db: %w", err)
	}

	products := make([]domain.Product, 0, len(productDBs))

	for _, productDB := range productDBs {
		product, err := parseProductDBToDomain(productDB)
		if err != nil {
			return nil, fmt.Errorf("failed to init product entity: %w", err)
		}

		products = append(products, *product)
	}

	return products, nil
}

func (r *DBRepo) sendOrderProcessingEvent(ctx context.Context, orderID valueobjects.OrderID) error {
	event := orderCreatedEvent{
		OrderID: orderID.UUID(),
	}

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = r.queueWriters.OrderCreated.Write(ctx,
		topicwriter.Message{
			Data: bytes.NewBuffer(message),
			Metadata: map[string][]byte{
				"service": []byte("go-template"),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write message to topic: %w", err)
	}

	return nil
}

func (r *DBRepo) DeleteOrder(ctx context.Context, orderID valueobjects.OrderID) error {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get primary dbx: %w", err)
	}

	if err := r.deleteOrderDB(ctx, db, orderID); err != nil {
		return fmt.Errorf("failed to delete order from db: %w", err)
	}

	return nil
}

func (r *DBRepo) ProcessOrders(ctx context.Context, updateFn func([]*domain.Order) error) error {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get primary dbx: %w", err)
	}

	if err = postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		orders, err := r.getCreatedOrdersForUpdate(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get created order ids: %w", err)
		}

		if err := updateFn(orders); err != nil {
			return fmt.Errorf("failed to update orders: %w", err)
		}

		for _, order := range orders {
			if err := r.sendOrderProcessingEvent(ctx, order.GetID()); err != nil {
				return fmt.Errorf("failed to send order processing event: %w", err)
			}
		}

		if err = r.updateOrders(ctx, tx, orders); err != nil {
			return fmt.Errorf("failed to set processing status: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed transaction for sending order processing events: %w", err)
	}

	return nil
}

func (r *DBRepo) updateOrders(ctx context.Context, tx *sqlx.Tx, orders []*domain.Order) (err error) {
	ctx, span := traces.CreateSpan(
		ctx,
		traces.ScopePostgres,
		"Postgres Update Orders",
		trace.SpanKindClient,
	)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	stmt, err := tx.PrepareContext(ctx, `UPDATE orders SET user_id = $2 , status = $3, product_ids = $4  WHERE id = $1`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, order := range orders {
		_, err := stmt.ExecContext(
			ctx,
			order.GetID().String(),
			order.GetUserID().UUID(),
			order.GetStatus().String(),
			pq.Array(parseutils.GoogleUUIDsToGofrsUUIDs(order.GetProductIDs().UUIDs())),
		)
		if err != nil {
			return fmt.Errorf("failed to update order status for order %s: %w", order.GetID().String(), err)
		}
	}

	return nil
}

func parseProductDomainToDB(product *domain.Product) productDB {
	return productDB{
		ID:     product.GetID().UUID(),
		Name:   product.GetName().String(),
		Price:  product.GetPrice().Float64(),
		Status: product.GetStatus().String(),
	}
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

	return domain.NewProduct(
		valueobjects.ProductID(productDB.ID),
		name,
		price,
		status,
	), nil
}

func parseOrderDomainToDB(order *domain.Order) orderDB {
	return orderDB{
		ID:         order.GetID().UUID(),
		UserID:     order.GetUserID().UUID(),
		Status:     order.GetStatus(),
		ProductIDs: parseutils.GoogleUUIDsToGofrsUUIDs(order.GetProductIDs().UUIDs()),
	}
}

func parseOrderDBToDomain(orderDB orderDB, products []domain.Product) *domain.Order {
	return domain.NewOrder(
		valueobjects.OrderID(orderDB.ID),
		valueobjects.UserID(orderDB.UserID),
		orderDB.Status,
		products,
	)
}

func (r *DBRepo) getProductsForUpdateDB(
	ctx context.Context,
	tx *sqlx.Tx,
	ids valueobjects.ProductIDs,
) (productDBs []productDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get Products For Update", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT id, name, price, status FROM products WHERE id = ANY($1) FOR UPDATE`

	err = tx.SelectContext(ctx, &productDBs, query, pq.Array(ids.Strings()))
	if err != nil {
		return nil, fmt.Errorf("failed to select products: %w", err)
	}

	return productDBs, nil
}

func (r *DBRepo) getOrderDB(ctx context.Context, db *sqlx.DB, id valueobjects.OrderID) (orderDB orderDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get Order", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT id, user_id, status, product_ids FROM orders WHERE id = $1`

	if err = db.QueryRowContext(ctx, query, id.String()).Scan(
		&orderDB.ID,
		&orderDB.UserID,
		&orderDB.Status,
		pq.Array(&orderDB.ProductIDs),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orderDB, domain.ErrOrderNotFound
		}

		return orderDB, fmt.Errorf("failed to select order: %w", err)
	}

	return orderDB, nil
}

func (r *DBRepo) getProductsDB(
	ctx context.Context,
	db *sqlx.DB,
	ids valueobjects.ProductIDs,
) (productDBs []productDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get Products", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT id, name, price, status FROM products WHERE id = ANY($1)`

	err = db.SelectContext(ctx, &productDBs, query, ids.Strings())
	if err != nil {
		return nil, fmt.Errorf("failed to select products: %w", err)
	}

	return productDBs, nil
}

func (r *DBRepo) deleteOrderDB(ctx context.Context, db *sqlx.DB, orderID valueobjects.OrderID) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Delete Order", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `DELETE FROM orders WHERE id = $1`

	if _, err := db.ExecContext(ctx, query, orderID.String()); err != nil {
		return fmt.Errorf("failed executing delete query: %w", err)
	}

	return nil
}
