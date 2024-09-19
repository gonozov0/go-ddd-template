package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	hasql "golang.yandex/hasql/sqlx"

	"go-echo-template/internal/domain/users"
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

func (r *PostgresRepo) SaveUser(ctx context.Context, entity users.User) error {
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

func (r *PostgresRepo) GetUser(ctx context.Context, id uuid.UUID) (*users.User, error) {
	db := r.cluster.StandbyPreferred().DBx()
	var user userDB
	query := `SELECT id, name, email FROM users WHERE id = $1`
	if err := db.GetContext(ctx, &user, query, id); err != nil {
		return nil, fmt.Errorf("failed to select user: %w", err)
	}
	entity, err := users.NewUser(user.ID, user.Name, user.Email)
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
