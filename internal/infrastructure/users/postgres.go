package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	hasql "golang.yandex/hasql/sqlx"

	domain "go-echo-template/internal/domain/users"
)

type userDB struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Email string    `db:"email"`
}

type PostgresRepo struct {
	cluster *hasql.Cluster
}

func NewPostgresRepo(cluster *hasql.Cluster) *PostgresRepo {
	return &PostgresRepo{
		cluster: cluster,
	}
}

func (r *PostgresRepo) CreateUser(
	ctx context.Context,
	email string,
	createFn func() (*domain.User, error),
) (*domain.User, error) {
	exist, err := r.checkUserExist(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user exist: %w", err)
	}
	if exist {
		return nil, domain.ErrUserAlreadyExist
	}

	user, err := createFn()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	db := r.cluster.Primary().DBx()
	userDB := userDB{
		ID:    user.ID(),
		Name:  user.Name(),
		Email: user.Email(),
	}
	query := `
		INSERT INTO users (id, name, email)
		VALUES (:id, :name, :email)
	`
	_, err = db.NamedExecContext(ctx, query, userDB)
	if err != nil {
		return nil, fmt.Errorf("failed to exec insert user query: %w", err)
	}

	return user, nil
}

func (r *PostgresRepo) UpdateUser(
	ctx context.Context,
	id uuid.UUID,
	updateFn func(*domain.User) (bool, error),
) (*domain.User, error) {
	db := r.cluster.Primary().DBx()
	user, err := r.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	updated, err := updateFn(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	if !updated {
		return user, nil
	}
	userDB := userDB{
		ID:    user.ID(),
		Name:  user.Name(),
		Email: user.Email(),
	}
	query := `
		UPDATE users
		SET name = :name, email = :email
		WHERE id = :id
	`
	_, err = db.NamedExecContext(ctx, query, userDB)
	if err != nil {
		return nil, fmt.Errorf("failed to exec update user query: %w", err)
	}
	return user, nil
}

func (r *PostgresRepo) SaveUser(ctx context.Context, entity domain.User) error {
	db := r.cluster.Primary().DBx()
	user := userDB{
		ID:    entity.ID(),
		Name:  entity.Name(),
		Email: entity.Email(),
	}
	query := `
		INSERT INTO users (id, name, email)
		VALUES (:id, :name, :email)
		ON CONFLICT (id) DO UPDATE SET name = :name, email = :email
	`
	_, err := db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	db := r.cluster.StandbyPreferred().DBx()
	var user userDB
	query := `SELECT id, name, email FROM users WHERE id = $1`
	if err := db.GetContext(ctx, &user, query, id); err != nil {
		return nil, fmt.Errorf("failed to select user: %w", err)
	}
	entity, err := domain.NewUser(user.ID, user.Name, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to init user entity: %w", err)
	}
	return entity, nil
}

func (r *PostgresRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	db := r.cluster.Primary().DBx()
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *PostgresRepo) checkUserExist(ctx context.Context, email string) (bool, error) {
	db := r.cluster.Primary().DBx()

	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	var exist bool
	err := db.GetContext(ctx, &exist, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to select user exist: %w", err)
	}

	return exist, nil
}
