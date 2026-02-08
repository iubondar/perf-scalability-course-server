package router

import (
	"github.com/go-chi/chi"
	"github.com/iubondar/perf-scalability-cource-server/internal/handlers"
	"github.com/iubondar/perf-scalability-cource-server/internal/payload"
)

func NewRouter() (chi.Router, error) {
	router := chi.NewRouter()

	helloWorldHandler := handlers.NewHelloWorldHandler()
	router.Get("/hello", helloWorldHandler.Handle)

	// payload
	cpuSleep := payload.NewGetrusagePayload()
	ioSleep := payload.NewIOPayload()
	router.Get("/payload", handlers.SleepHandler(cpuSleep, ioSleep))

	return router, nil
}
