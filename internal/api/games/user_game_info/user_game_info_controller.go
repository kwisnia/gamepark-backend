package user_game_info

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserGameInfoHandler(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("userID")
	gameUserInfo, _ := GetUserGameInfo(slug, userID)
	c.JSON(http.StatusOK, gameUserInfo)
}
