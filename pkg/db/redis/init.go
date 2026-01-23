package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"fmt"
)

type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Ping(ctx context.Context) *redis.StatusCmd
	Close() error
}

type Client struct {
	client redisClient
}

func NewClient(ctx context.Context, config Config) (*Client, error) {
	var (
		client    redisClient
		tlsConfig *tls.Config
	)

	if config.TlsEnabled {
		rootCertPool := x509.NewCertPool()

		pem, err := os.ReadFile(config.TlsRootCert)
		if err != nil {
			return nil, fmt.Errorf("Failed to read certificate")
		}

		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, fmt.Errorf("Failed to append PEM")
		}

		tlsConfig = &tls.Config{
			RootCAs: rootCertPool,
		}
	}

	if len(config.Addrs) < 1 {
		return nil, fmt.Errorf("no host found")
	}

	client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:     config.Addrs,
		Username:  config.Username,
		Password:  config.Password,
		TLSConfig: tlsConfig,
	})

	return &Client{client: client}, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(ctx, key, value, expiration)
}

func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.client.Get(ctx, key)
}

func (c *Client) Del(ctx context.Context, key string) *redis.IntCmd {
	return c.client.Del(ctx, key)
}

func (c *Client) Ping(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}

func (c *Client) Close() error {
	return c.client.Close()
}
