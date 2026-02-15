package redis

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// New creates a Redis client, verifies connectivity, and seeds from PostgreSQL.
func New(addr string, pool *pgxpool.Pool) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	if err := SeedFromPG(context.Background(), client, pool); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}
