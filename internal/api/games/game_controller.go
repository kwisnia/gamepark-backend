package games

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GamePage struct {
	data       []GameListElement
	nextCursor uint
}

// get games
func GetGames(c *gin.Context) {
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page size"})
		return
	}
	afterId, err := strconv.Atoi(c.DefaultQuery("after", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid after"})
		return
	}
	games, _ := getGames(pageSize, afterId)
	response := GamePage{
		data:       games,
		nextCursor: games[len(games)-1].ID,
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       response.data,
		"nextCursor": response.nextCursor,
	})
}

// get game by slug
func GetGame(c *gin.Context) {
	slug := c.Param("slug")
	game, _ := getBySlug(slug)
	c.JSON(http.StatusOK, game)
}
