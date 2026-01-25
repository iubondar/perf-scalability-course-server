package router

import (
	"github.com/go-chi/chi"
	"github.com/iubondar/perf-scalability-cource-server/internal/handler"
)

func NewRouter() (chi.Router, error) {
	router := chi.NewRouter()

	helloWorldHandler := handler.NewHelloWorldHandler()
	router.Get("/", helloWorldHandler.Handle)

	return router, nil
}
