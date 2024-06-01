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
}

func NewConfig() Config {
	var config Config

	flag.StringVar(&config.Host, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.DatabaseUri, "d", "", "database uri")
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
