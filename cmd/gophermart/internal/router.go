package server

import (
	"github.com/gin-gonic/gin"
	// "github.com/pisarevaa/metrics/internal/server/storage"
)

const readTimeout = 5
const writeTimout = 10

func NewRouter() {
	r := gin.Default()
	r.GET("/ping", Hello)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080"
}
