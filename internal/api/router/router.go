package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/internal/api/games"
	"github.com/kwisnia/inzynierka-backend/internal/api/middleware"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/login", user.LoginUser)
	r.POST("/register", user.RegisterUser)
	r.GET("/games", games.GetGames)
	r.GET("/games/:slug", games.GetGame)
	return r
}
