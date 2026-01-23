package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"go-ddd-template/internal/domain/shared/valueobjects"
	domain "go-ddd-template/internal/domain/users"
	pkgredis "go-ddd-template/pkg/db/redis"
)

type userRedis struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RedisRepo struct {
	client     *pkgredis.Client
	expiration time.Duration
}

const usersExpiration = 15 * time.Minute

func NewRedisRepo(client *pkgredis.Client) *RedisRepo {
	return &RedisRepo{
		client:     client,
		expiration: usersExpiration,
	}
}

func (r *RedisRepo) SaveUser(
	ctx context.Context,
	user *domain.User,
) error {
	ru := userRedis{
		Name:  user.GetName().String(),
		Email: user.GetEmail().String(),
	}

	val, err := json.Marshal(ru)
	if err != nil {
		return fmt.Errorf("failed to serialize user: %w", err)
	}

	err = r.client.Set(ctx, user.GetID().String(), val, r.expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to save user to redis: %w", err)
	}

	return nil
}

func (r *RedisRepo) GetUser(ctx context.Context, id valueobjects.UserID) (*domain.User, error) {
	val, err := r.client.Get(ctx, id.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user from redis: %w", err)
	}

	var user userRedis

	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize user: %w", err)
	}

	return domain.NewUser(valueobjects.UserID(id), domain.Name(user.Name), valueobjects.Email(user.Email)), nil
}

func (r *RedisRepo) DeleteUser(ctx context.Context, id valueobjects.UserID) error {
	err := r.client.Del(ctx, id.String()).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user from redis: %w", err)
	}

	return nil
}
