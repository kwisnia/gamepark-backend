package games

import (
	"fmt"
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
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid after"})
		return
	}
	filters := c.QueryArray("filters")
	parsedFilters := make([]int, len(filters))
	fmt.Println(filters)
	for i, filter := range filters {
		parsedFilters[i], err = strconv.Atoi(filter)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid filter"})
			return
		}
	}
	sort := c.DefaultQuery("sort", "id.asc")
	games, _ := getGames(pageSize, page, parsedFilters, sort)
	c.JSON(http.StatusOK, games)
}

// get game by slug
func GetGame(c *gin.Context) {
	slug := c.Param("slug")
	game, _ := getBySlug(slug)
	c.JSON(http.StatusOK, game)
}
