package router

import (
	"github.com/go-chi/chi"
	"github.com/iubondar/perf-scalability-cource-server/internal/handlers"
	"github.com/iubondar/perf-scalability-cource-server/internal/payload"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(pool *pgxpool.Pool) (chi.Router, error) {
	router := chi.NewRouter()

	helloWorldHandler := handlers.NewHelloWorldHandler()
	router.Get("/hello", helloWorldHandler.Handle)

	// payload
	cpuSleep := payload.NewGetrusagePayload()
	ioSleep := payload.NewIOPayload()
	router.Get("/payload", handlers.SleepHandler(cpuSleep, ioSleep))

	router.Get("/pg", handlers.PgHandler(pool))

	return router, nil
}
