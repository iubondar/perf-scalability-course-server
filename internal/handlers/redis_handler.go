package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func RedisHandler(client *redis.Client) http.HandlerFunc {
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

		id := uuid.New()
		key := "item:" + id.String()

		var lastItem *Item
		for i := 0; i < num; i++ {
			val, err := client.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
					JSONError(w, map[string]string{"error": "item not found"}, http.StatusNotFound)
					return
				}
				JSONError(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
				return
			}
			var item Item
			if err := json.Unmarshal([]byte(val), &item); err != nil {
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
