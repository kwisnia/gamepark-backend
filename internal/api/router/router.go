package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/internal/api/middleware"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.AuthRequired())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}