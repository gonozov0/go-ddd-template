package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
	"go-ddd-template/internal/infrastructure/users/internal"
	"go-ddd-template/pkg/db/postgres"
	"go-ddd-template/pkg/db/redis"
	loggerutils "go-ddd-template/pkg/logger/utils"
)

type DBRepo struct {
	postgres *internal.PostgresRepo
	redis    *internal.RedisRepo
}

func NewDBRepo(postgresCluster *postgres.Cluster, redisClient *redis.Client) *DBRepo {
	return &DBRepo{
		postgres: internal.NewPostgresRepo(postgresCluster),
		redis:    internal.NewRedisRepo(redisClient),
	}
}

func (r *DBRepo) CreateUser(
	ctx context.Context,
	createFn func() (*domain.User, error),
) (*domain.User, error) {
	user, err := r.postgres.CreateUser(ctx, createFn)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in postgres: %w", err)
	}

	err = r.redis.SaveUser(ctx, user)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to save user to redis cache",
			loggerutils.ErrAttr(fmt.Errorf("failed to save user to redis cache: %w", err)),
		)

		return nil, fmt.Errorf("failed to save user in cache: %w", err)
	}

	return user, nil
}

func (r *DBRepo) GetUser(ctx context.Context, id valueobjects.UserID) (*domain.User, error) {
	user, err := r.redis.GetUser(ctx, id)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			slog.ErrorContext(
				ctx,
				"failed to read user from redis cache",
				loggerutils.ErrAttr(fmt.Errorf("failed to read user from redis cache: %w", err)),
			)

			return nil, fmt.Errorf("failed to get user from cache: %w", err)
		}
	} else {
		return user, nil
	}

	user, err = r.postgres.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read user from postgres: %w", err)
	}

	err = r.redis.SaveUser(ctx, user)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to save user to redis cache",
			loggerutils.ErrAttr(fmt.Errorf("failed to save user to redis cache: %w", err)),
		)

		return nil, fmt.Errorf("failed to save user in cache: %w", err)
	}

	return user, nil
}

func (r *DBRepo) UpdateUser(
	ctx context.Context,
	id valueobjects.UserID,
	updateFn func(*domain.User) error,
) (*domain.User, error) {
	user, err := r.postgres.UpdateUser(ctx, id, updateFn)
	if err != nil {
		return nil, fmt.Errorf("failed to update user in postgres: %w", err)
	}

	err = r.redis.SaveUser(ctx, user)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to update user in redis cache",
			loggerutils.ErrAttr(fmt.Errorf("failed to update user in redis cache: %w", err)),
		)

		return nil, fmt.Errorf("failed to update user in cache: %w", err)
	}

	return user, nil
}

func (r *DBRepo) DeleteUser(ctx context.Context, id valueobjects.UserID, deleteFn func(*domain.User) error) error {
	err := r.postgres.DeleteUser(ctx, id, deleteFn)
	if err != nil {
		return fmt.Errorf("failed to delete user from postgres: %w", err)
	}

	err = r.redis.DeleteUser(ctx, id)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to delete user from redis cache",
			loggerutils.ErrAttr(fmt.Errorf("failed to delete user from redis cache: %w", err)),
		)

		return fmt.Errorf("failed to delete user from cache: %w", err)
	}

	return nil
}
