package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
	"go-ddd-template/pkg/db/postgres"
	"go-ddd-template/pkg/traces"
)

type userDB struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Email string    `db:"email"`
}

type PostgresRepo struct {
	cluster *postgres.Cluster
}

func NewPostgresRepo(cluster *postgres.Cluster) *PostgresRepo {
	return &PostgresRepo{
		cluster: cluster,
	}
}

func (r *PostgresRepo) CreateUser(
	ctx context.Context,
	createFn func() (*domain.User, error),
) (*domain.User, error) {
	var user *domain.User

	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary dbx: %w", err)
	}

	err = postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		user, err = createFn()
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if err := r.checkUserEmailAlreadyExists(ctx, tx, user.GetEmail().String()); err != nil {
			return fmt.Errorf("failed to check user exist: %w", err)
		}

		userDB := userDB{
			ID:    user.GetID().UUID(),
			Name:  user.GetName().String(),
			Email: user.GetEmail().String(),
		}

		if err := r.createUserDB(ctx, tx, userDB); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})

	return user, nil
}

func (r *PostgresRepo) UpdateUser(
	ctx context.Context,
	id valueobjects.UserID,
	updateFn func(*domain.User) error,
) (*domain.User, error) {
	var user *domain.User

	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary dbx: %w", err)
	}

	err = postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		user, err = r.getUserForUpdate(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get user for update: %w", err)
		}

		oldEmail := user.GetEmail()

		err = updateFn(user)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		if oldEmail != user.GetEmail() {
			if err := r.checkUserEmailAlreadyExists(ctx, tx, user.GetEmail().String()); err != nil {
				return fmt.Errorf("failed to check user exist: %w", err)
			}
		}

		userDB := userDB{
			ID:    user.GetID().UUID(),
			Name:  user.GetName().String(),
			Email: user.GetEmail().String(),
		}

		if err := r.updateUserDB(ctx, tx, userDB); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepo) getUserForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	id valueobjects.UserID,
) (*domain.User, error) {
	user, err := r.getUserDBForUpdate(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}

	userName, err := domain.NewName(user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to init user name: %w", err)
	}

	userEmail, err := valueobjects.NewEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to init user email: %w", err)
	}

	return domain.NewUser(valueobjects.UserID(user.ID), userName, userEmail), nil
}

func (r *PostgresRepo) GetUser(ctx context.Context, id valueobjects.UserID) (*domain.User, error) {
	db, err := r.cluster.StandbyPreferredDBx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get standby preferred dbx: %w", err)
	}

	user, err := r.getUserDB(ctx, db, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}

	userName, err := domain.NewName(user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to init user name: %w", err)
	}

	userEmail, err := valueobjects.NewEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to init user email: %w", err)
	}

	return domain.NewUser(valueobjects.UserID(user.ID), userName, userEmail), nil
}

func (r *PostgresRepo) DeleteUser(
	ctx context.Context,
	id valueobjects.UserID,
	deleteFn func(*domain.User) error,
) error {
	db, err := r.cluster.PrimaryDBx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get primary dbx: %w", err)
	}

	return postgres.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		user, err := r.getUserForUpdate(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get user for update: %w", err)
		}

		err = deleteFn(user)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		if err := r.deleteUserDB(ctx, tx, id); err != nil {
			return fmt.Errorf("failed to delete user from db: %w", err)
		}

		return nil
	})
}

func (r *PostgresRepo) checkUserEmailAlreadyExists(ctx context.Context, tx *sqlx.Tx, email string) (err error) {
	exist, err := r.checkUserEmailExistsDB(ctx, tx, email)
	if err != nil {
		return fmt.Errorf("failed to check user email exists: %w", err)
	}

	if exist {
		return domain.ErrUserAlreadyExist
	}

	return nil
}

func (r *PostgresRepo) createUserDB(ctx context.Context, tx *sqlx.Tx, userDB userDB) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Create User", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `
		INSERT INTO users (id, name, email)
		VALUES (:id, :name, :email)
	`

	_, err = tx.NamedExecContext(ctx, query, userDB)
	if err != nil {
		return fmt.Errorf("failed to exec insert user query: %w", err)
	}

	return nil
}

func (r *PostgresRepo) updateUserDB(ctx context.Context, tx *sqlx.Tx, userDB userDB) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Update User", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `
		UPDATE users
		SET name = :name, email = :email
		WHERE id = :id
	`

	_, err = tx.NamedExecContext(ctx, query, userDB)
	if err != nil {
		return fmt.Errorf("failed to exec update user query: %w", err)
	}

	return nil
}

func (r *PostgresRepo) deleteUserDB(ctx context.Context, tx *sqlx.Tx, id valueobjects.UserID) (err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Delete User", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `DELETE FROM users WHERE id = $1`

	_, err = tx.ExecContext(ctx, query, id.UUID())
	if err != nil {
		return fmt.Errorf("failed to exec delete user query: %w", err)
	}

	return nil
}

func (r *PostgresRepo) getUserDB(
	ctx context.Context,
	db *sqlx.DB,
	id valueobjects.UserID,
) (user userDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get User", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT id, name, email FROM users WHERE id = $1`
	if err := db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userDB{}, domain.ErrUserNotFound
		}

		return userDB{}, fmt.Errorf("failed to select user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepo) getUserDBForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	id valueobjects.UserID,
) (user userDB, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get User For Update", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT id, name, email FROM users WHERE id = $1 FOR UPDATE`
	if err := tx.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userDB{}, domain.ErrUserNotFound
		}

		return userDB{}, fmt.Errorf("failed to select user for update: %w", err)
	}

	return user, nil
}

func (r *PostgresRepo) checkUserEmailExistsDB(ctx context.Context, tx *sqlx.Tx, email string) (exist bool, err error) {
	ctx, span := traces.CreateSpan(ctx, traces.ScopePostgres, "Postgres Get User For Update", trace.SpanKindClient)

	defer func() {
		traces.SetSpanStatus(span, err)
		span.End()
	}()

	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	err = tx.GetContext(ctx, &exist, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to select user exist: %w", err)
	}

	return exist, nil
}
