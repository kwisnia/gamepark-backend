package games

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// get games
func GetGames(c *gin.Context) {
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page size"})
		return
	}
	afterId, err := strconv.Atoi(c.DefaultQuery("afterId", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid afterId"})
		return
	}
	games, _ := getGames(pageSize, afterId)
	c.JSON(http.StatusOK, games)
}

// get game by slug
func GetGame(c *gin.Context) {
	slug := c.Param("slug")
	game, _ := getBySlug(slug)
	c.JSON(http.StatusOK, game)
}
