package config

import (
	"flag"

	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

type Config struct {
	RunAddress  string `env:"RUN_ADDRESS"`
	DatabaseDSN string `env:"DATABASE_DSN"`
	RedisAddr   string `env:"REDIS_ADDR"`
}

// для локальной разработки
const (
	defaultRunAddress  = "localhost:8000"
	defaultDatabaseDSN = "postgres://localhost:5432/perf?sslmode=disable"
	defaultRedisAddr  = "localhost:6379"
)

func NewConfig(progname string, args []string) (*Config, error) {
	var c Config

	// https://eli.thegreenplace.net/2020/testing-flag-parsing-in-go-programs/
	// Загружаем значения из переданных аргументов командной строки
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)

	flags.StringVar(&c.RunAddress, "a", defaultRunAddress, "address to run server")
	flags.StringVar(&c.DatabaseDSN, "d", defaultDatabaseDSN, "postgres connection DSN")
	flags.StringVar(&c.RedisAddr, "r", defaultRedisAddr, "redis address")

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	// Переписываем значения из переменных окружения
	err = env.Parse(&c)
	if err != nil {
		return nil, err
	}

	zap.L().Sugar().Debugln(
		"Config: ",
		"RunAddress", c.RunAddress,
		"DatabaseDSN", c.DatabaseDSN,
		"RedisAddr", c.RedisAddr,
	)

	return &c, nil
}
