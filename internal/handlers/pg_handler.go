package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Item struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type PgResponse struct {
	WallTimeMSec float64 `json:"wall_time_msec"`
	Item         *Item   `json:"item,omitempty"`
}

func PgHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		num := 1
		if numStr := r.URL.Query().Get("num"); numStr != "" {
			var err error
			num, err = strconv.Atoi(numStr)
			if err != nil || num < 1 {
				JSONError(w, map[string]string{"error": "invalid num parameter"}, http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()

		var id uuid.UUID
		err := pool.QueryRow(ctx, "SELECT id FROM items ORDER BY RANDOM() LIMIT 1").Scan(&id)
		if err != nil {
			JSONError(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}

		var lastItem *Item
		for i := 0; i < num; i++ {
			var item Item
			err := pool.QueryRow(ctx, "SELECT id, name, description, created_at FROM items WHERE id = $1", id).Scan(
				&item.ID, &item.Name, &item.Description, &item.CreatedAt,
			)
			if err != nil {
				JSONError(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
				return
			}
			lastItem = &item
		}

		JSONResponse(w, PgResponse{
			WallTimeMSec: float64(time.Since(start).Nanoseconds()) / nanosecToMillisec,
			Item:         lastItem,
		})
	}
}
