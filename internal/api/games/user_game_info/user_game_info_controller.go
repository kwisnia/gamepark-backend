package user_game_info

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserGameInfoHandler(c *gin.Context) {
	slug := c.Param("slug")
	userName := c.GetString("userName")
	gameUserInfo, _ := GetUserGameInfo(slug, userName)
	c.JSON(http.StatusOK, gameUserInfo)
}
