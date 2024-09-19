// Description: The file contains just the example of the Redis repository implementation and isn't used in the project.

package users

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"go-echo-template/internal/domain/users"
)

type redisUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type RedisRepo struct {
	client     redisClient
	expiration time.Duration
}

func NewRedisRepo(clusterMode, tlsEnabled bool, addr, username, password string, expiration time.Duration) *RedisRepo {
	var (
		tlsConfig *tls.Config
		client    redisClient
	)
	if tlsEnabled {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // It's okay in intranet
		}
	}

	if clusterMode {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:     []string{addr},
			Username:  username,
			Password:  password,
			TLSConfig: tlsConfig,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:      addr,
			Username:  username,
			Password:  password,
			TLSConfig: tlsConfig,
		})
	}

	return &RedisRepo{
		client:     client,
		expiration: expiration,
	}
}

func (r *RedisRepo) SaveUser(ctx context.Context, user users.User) error {
	ru := redisUser{
		Name:  user.Name(),
		Email: user.Email(),
	}
	val, err := json.Marshal(ru)
	if err != nil {
		return fmt.Errorf("failed to serialize user: %w", err)
	}

	err = r.client.Set(ctx, user.ID().String(), val, r.expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to save user to redis: %w", err)
	}
	return nil
}

func (r *RedisRepo) GetUser(ctx context.Context, id uuid.UUID) (*users.User, error) {
	val, err := r.client.Get(ctx, id.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, users.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user from redis: %w", err)
	}

	var ru redisUser
	err = json.Unmarshal([]byte(val), &ru)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize user: %w", err)
	}

	return users.NewUser(id, ru.Name, ru.Email)
}

func (r *RedisRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.client.Del(ctx, id.String()).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user from redis: %w", err)
	}
	return nil
}
