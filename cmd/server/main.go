package main

import (
	"log"
	"os"

	"github.com/iubondar/perf-scalability-cource-server/internal/config"
	"github.com/iubondar/perf-scalability-cource-server/internal/router"
	"github.com/iubondar/perf-scalability-cource-server/internal/server"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	router, err := router.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	config, err := config.NewConfig(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Создаем и запускаем сервер
	srv := server.New(config.RunAddress, router)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
