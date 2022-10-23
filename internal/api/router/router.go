package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/internal/api/games"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/lists"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/reviews"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/user_game_info"
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
	r.GET("games/:slug/user", middleware.AuthRequired(), user_game_info.GetUserGameInfoHandler)
	r.GET("/me/details", middleware.AuthRequired(), user.GetDetails)
	r.GET("/:userName/details", user.GetDetailsByUsernameHandler)
	r.GET("/:userName/lists", lists.GetUserLists)
	r.GET("/:userName/reviews", reviews.GetReviewsForUserHandler)
	r.GET("/list/:id", lists.GetUserList)
	r.POST("/list", middleware.AuthRequired(), lists.CreateList)
	r.PATCH("/list/:id", middleware.AuthRequired(), lists.UpdateList)
	r.DELETE("/list/:id", middleware.AuthRequired(), lists.DeleteList)
	r.POST("/list/:id/add", middleware.AuthRequired(), lists.AddGameToList)
	r.POST("/list/:id/remove", middleware.AuthRequired(), lists.RemoveGameFromList)
	r.GET("/games/:slug/reviews", middleware.AuthOptional(), reviews.GetReviewsForGameHandler)
	r.POST("/games/:slug/reviews", middleware.AuthRequired(), reviews.CreateReviewHandler)
	//r.PATCH("/games/:slug/reviews/:id", middleware.AuthRequired(), reviews.UpdateReviewHandler)
	r.DELETE("/games/:slug/reviews", middleware.AuthRequired(), reviews.DeleteReviewHandler)
	r.POST("/games/:slug/reviews/:reviewId/helpful", middleware.AuthRequired(), reviews.MarkReviewAsHelpfulHandler)
	r.DELETE("/games/:slug/reviews/:reviewId/helpful", middleware.AuthRequired(), reviews.UnmarkReviewAsHelpfulHandler)

	return r
}
