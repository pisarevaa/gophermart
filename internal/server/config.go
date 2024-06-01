package server

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host                 string `env:"RUN_ADDRESS"`
	DatabaseUri          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	GinMode              string `env:"GIN_MODE"`
}

func NewConfig() Config {
	var config Config

	flag.StringVar(&config.Host, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.GinMode, "g", "debug", "gin server logs mode")
	flag.StringVar(&config.DatabaseUri, "d", "postgres://gophermart:CC7B02B06C4C1CF81FAE7D8C46C429EC@localhost:5432/gophermart?sslmode=disable", "database uri")
	flag.StringVar(&config.AccrualSystemAddress, "r", "", "charging system address")
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
	if envConfig.DatabaseUri != "" {
		config.DatabaseUri = envConfig.DatabaseUri
	}
	if envConfig.AccrualSystemAddress != "" {
		config.AccrualSystemAddress = envConfig.AccrualSystemAddress
	}
	return config
}
