package main

import (
	"log"
	"os"

	"github.com/iubondar/perf-scalability-cource-server/internal/config"
	"github.com/iubondar/perf-scalability-cource-server/internal/database"
	"github.com/iubondar/perf-scalability-cource-server/internal/router"
	"github.com/iubondar/perf-scalability-cource-server/internal/server"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
}

func main() {
	cfg, err := config.NewConfig(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigrations(cfg.DatabaseDSN); err != nil {
		log.Fatal(err)
	}

	pool, err := database.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	router, err := router.NewRouter(pool)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(cfg.RunAddress, router)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
