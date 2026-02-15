package redis

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/iubondar/perf-scalability-cource-server/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const keyPrefix = "item:"

// SeedFromPG reads all items from PostgreSQL and stores them in Redis by UUID key.
func SeedFromPG(ctx context.Context, client *redis.Client, pool *pgxpool.Pool) error {
	rows, err := pool.Query(ctx, "SELECT id, name, description, created_at FROM items")
	if err != nil {
		return err
	}
	defer rows.Close()

	pipe := client.Pipeline()
	for rows.Next() {
		var item handlers.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt); err != nil {
			return err
		}
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		pipe.Set(ctx, keyPrefix+item.ID.String(), data, 0)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	_, err = pipe.Exec(ctx)
	return err
}

// ItemKey returns the Redis key for an item UUID.
func ItemKey(id uuid.UUID) string {
	return keyPrefix + id.String()
}
