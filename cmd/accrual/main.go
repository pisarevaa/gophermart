package main

import (
	"math/rand/v2"
	"net/http"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type OrderUri struct {
	Number    string `uri:"number"              binding:"required"`
	IsInvalid bool   `uri:"isInvalid,omitempty"`
}

type Config struct {
	Host string `env:"RUN_ADDRESS"`
}

type OrderReponse struct {
	Number  string `json:"number"            binding:"required"`
	Status  string `json:"status"            binding:"required"`
	Accrual int64  `json:"accrual,omitempty"`
}

func NewLogger() *zap.SugaredLogger {
	logger := zap.NewExample().Sugar()
	defer logger.Sync() //nolint:errcheck // ignore check
	return logger
}

func NewConfig() Config {
	var config Config
	flag.StringVar(&config.Host, "a", "localhost:8085", "host and port to run server")
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
	return config
}

func GetOrder(c *gin.Context) {
	var orderUri OrderUri
	if err := c.ShouldBindUri(&orderUri); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var statuses []string

	if orderUri.IsInvalid {
		statuses = []string{
			"REGISTERED",
			"INVALID",
			"PROCESSING",
		}
	} else {
		statuses = []string{
			"REGISTERED",
			"PROCESSED",
			"PROCESSING",
		}
	}
	randomInt := rand.IntN(len(statuses))
	randomStatus := statuses[randomInt]

	orderResponse := OrderReponse{
		Number: orderUri.Number,
		Status: randomStatus,
	}

	if randomStatus == "PROCESSED" {
		orderResponse.Accrual = rand.Int64N(500)
	}

	c.JSON(http.StatusOK, orderResponse)
}

func main() {
	cfg := NewConfig()
	logger := NewLogger()
	r := gin.Default()
	r.GET("/api/orders/:number", GetOrder)
	logger.Info("Run Server")
	logger.Fatal(endless.ListenAndServe(cfg.Host, r))
}
