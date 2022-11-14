package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/api/chat"
	"github.com/kwisnia/inzynierka-backend/internal/api/dashboard"
	"github.com/kwisnia/inzynierka-backend/internal/api/file"
	"github.com/kwisnia/inzynierka-backend/internal/api/followers"
	"github.com/kwisnia/inzynierka-backend/internal/api/games"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/discussions"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/lists"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/reviews"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/user_game_info"
	"github.com/kwisnia/inzynierka-backend/internal/api/middleware"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"github.com/kwisnia/inzynierka-backend/internal/api/websocket"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/login", user.LoginUserHandler)
	r.POST("/register", user.RegisterUserHandler)
	r.GET("/users", user.GetUsersHandler)
	r.GET("/games", games.GetGamesHandler)
	r.GET("/games/test", games.GameWebhookCreateHandler)
	r.GET("/games/:slug", games.GetGameHandler)
	r.GET("/games/:slug/info", games.GetGameShortInfoHandler)
	r.GET("games/:slug/user", middleware.AuthRequired(), user_game_info.GetUserGameInfoHandler)
	r.GET("/me/details", middleware.AuthRequired(), user.GetDetailsHandler)
	r.PATCH("/me/details", middleware.AuthRequired(), user.UpdateUserProfileHandler)
	r.PATCH("/me/details/banner", middleware.AuthRequired(), user.UpdateUserBannerPositionHandler)
	r.GET("/:userName/details", user.GetDetailsByUsernameHandler)
	r.GET("/:userName/lists", lists.GetUserListsHandler)
	r.GET("/:userName/reviews", middleware.AuthOptional(), reviews.GetReviewsForUserHandler)
	r.GET("/:userName/discussions", middleware.AuthOptional(), discussions.GetDiscussionsForUserHandler)
	r.GET("/:userName/followers", followers.GetUserFollowersHandler)
	r.GET("/:userName/following", followers.GetUserFollowingHandler)
	r.GET("/:userName/achievements", user.GetUserAchievementsHandler)
	r.GET("/list/:id", lists.GetUserListHandler)
	r.POST("/list", middleware.AuthRequired(), lists.CreateListHandler)
	r.PATCH("/list/:id", middleware.AuthRequired(), lists.UpdateListHandler)
	r.DELETE("/list/:id", middleware.AuthRequired(), lists.DeleteListHandler)
	r.POST("/list/:id/add", middleware.AuthRequired(), lists.AddGameToListHandler)
	r.POST("/list/:id/remove", middleware.AuthRequired(), lists.RemoveGameFromListHandler)
	r.GET("/games/:slug/reviews", middleware.AuthOptional(), reviews.GetReviewsForGameHandler)
	r.POST("/games/:slug/reviews", middleware.AuthRequired(), reviews.CreateReviewHandler)
	//r.PATCH("/games/:slug/reviews/:id", middleware.AuthRequired(), reviews.UpdateReviewHandler)
	r.DELETE("/games/:slug/reviews", middleware.AuthRequired(), reviews.DeleteReviewHandler)
	r.POST("/games/:slug/reviews/:reviewID/helpful", middleware.AuthRequired(), reviews.MarkReviewAsHelpfulHandler)
	r.DELETE("/games/:slug/reviews/:reviewID/helpful", middleware.AuthRequired(), reviews.UnmarkReviewAsHelpfulHandler)
	r.GET("/ws", middleware.QueryAuthRequired(), websocket.WebSockerConnectionHandler)
	r.GET("/chat/:user", middleware.AuthRequired(), chat.GetChatHistoryHandler)
	r.GET("/chat/history", middleware.AuthRequired(), chat.GetChatReceiversHandler)
	r.POST("/image", middleware.AuthRequired(), file.UploadImageHandler)
	r.GET("/achievements", achievements.GetAllAchievementsHandler)
	r.GET("/dashboard", middleware.AuthRequired(), dashboard.GetActivitiesFromFollowedHandler)
	getDiscussionRoutes(r)
	getFollowRoutes(r)
	return r
}

func getDiscussionRoutes(router *gin.Engine) {
	discussionRoutes := router.Group("/games/:slug/discussions")
	discussionRoutes.GET("", middleware.AuthOptional(), discussions.GetDiscussionsForGameHandler)
	discussionRoutes.POST("", middleware.AuthRequired(), discussions.CreateDiscussionHandler)
	discussionRoutes.GET("/:discussionId", middleware.AuthOptional(), discussions.GetDiscussionHandler)
	discussionRoutes.DELETE("/:discussionId", middleware.AuthRequired(), discussions.DeleteDiscussionHandler)
	discussionRoutes.POST("/:discussionId/score", middleware.AuthRequired(), discussions.ScoreDiscussionHandler)
	discussionRoutes.GET("/:discussionId/posts", middleware.AuthOptional(), discussions.GetDiscussionPostsHandler)
	discussionRoutes.POST("/:discussionId/posts", middleware.AuthRequired(), discussions.CreateDiscussionPostHandler)
	discussionRoutes.PATCH("/:discussionId/posts/:postId", middleware.AuthRequired(), discussions.UpdateDiscussionPostHandler)
	discussionRoutes.DELETE("/:discussionId/posts/:postId", middleware.AuthRequired(), discussions.DeleteDiscussionPostHandler)
	discussionRoutes.POST("/:discussionId/posts/:postId/score", middleware.AuthRequired(), discussions.ScoreDiscussionPostHandler)
	discussionRoutes.GET("/:discussionId/posts/:postId/replies", middleware.AuthOptional(), discussions.GetDiscussionPostRepliesHandler)
}

func getFollowRoutes(router *gin.Engine) {
	followRoutes := router.Group("/follow")
	followRoutes.POST("/:userName", middleware.AuthRequired(), followers.FollowUserHandler)
	followRoutes.DELETE("/:userName", middleware.AuthRequired(), followers.UnfollowUserHandler)
	followRoutes.GET("/:userName", middleware.AuthRequired(), followers.CheckFollowConnectionHandler)
}
