package configs

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host                 string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	GinMode              string `env:"GIN_MODE"`
	SecretKey            string `env:"SECRET_KEY"`
	TokenExpSec          int64  `env:"TOKEN_EXP"`
	TaskInterval         int64  `env:"TASK_INTERVAL"`
}

func NewConfig() Config {
	var config Config

	flag.StringVar(&config.Host, "a", "localhost:8085", "address and port to run server")
	flag.StringVar(&config.GinMode, "g", "debug", "gin server logs mode")
	flag.StringVar(
		&config.DatabaseURI,
		"d",
		"postgres://gophermart:CC7B02B06C4C1CF81FAE7D8C46C429EC@localhost:5432/gophermart?sslmode=disable",
		"database uri",
	)
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://localhost:8080", "charging system address")
	flag.StringVar(&config.SecretKey, "k", "7fd315fd5f381bb9035d003dbd904102", "secret key to hash password")
	flag.Int64Var(&config.TokenExpSec, "t", 7200, "time in sec to expire token")
	flag.Int64Var(&config.TaskInterval, "i", 1, "time in sec to update order statuses")
	flag.Parse()
	if len(flag.Args()) > 0 {
		log.Fatal("used not declared arguments")
	}

	var envConfig Config
	err := env.Parse(&envConfig)
	if err != nil {
		log.Fatal(err)
	}

	if envConfig.GinMode != "" {
		config.GinMode = envConfig.GinMode
	}
	if envConfig.Host != "" {
		config.Host = envConfig.Host
	}
	if envConfig.DatabaseURI != "" {
		config.DatabaseURI = envConfig.DatabaseURI
	}
	if envConfig.AccrualSystemAddress != "" {
		config.AccrualSystemAddress = envConfig.AccrualSystemAddress
	}
	if envConfig.SecretKey != "" {
		config.SecretKey = envConfig.SecretKey
	}
	if envConfig.TokenExpSec != 0 {
		config.TokenExpSec = envConfig.TokenExpSec
	}
	if envConfig.TaskInterval != 0 {
		config.TaskInterval = envConfig.TaskInterval
	}
	return config
}
