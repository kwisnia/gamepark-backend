package games

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type GameForm struct {
	ID uint `json:"id"`
}

// get games
func GetGamesHandler(c *gin.Context) {
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
	search := c.DefaultQuery("search", "")
	filters := c.QueryArray("filters")
	parsedFilters := make([]int, len(filters))
	for i, filter := range filters {
		parsedFilters[i], err = strconv.Atoi(filter)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid filter"})
			return
		}
	}
	sort := c.DefaultQuery("sort", "id.asc")
	games, _ := GetGames(pageSize, page, parsedFilters, sort, search)
	c.JSON(http.StatusOK, games)
}

// get game by slug
func GetGameHandler(c *gin.Context) {
	slug := c.Param("slug")
	game, _ := GetBySlug(slug)
	c.JSON(http.StatusOK, game)
}

func GetGameShortInfoHandler(c *gin.Context) {
	slug := c.Param("slug")
	game, err := GetShortInfoBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred"})
	}
	c.JSON(http.StatusOK, game)
}

func GameWebhookCreateHandler(c *gin.Context) {
	var gameForm GameForm
	err := c.ShouldBindJSON(&gameForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	secret := c.GetHeader("X-Secret")
	if secret != config.GetEnv("WEBHOOK_SECRET") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid secret"})
		return
	}
	err = CreateWebhookGame(int(gameForm.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func GameWebhookDeleteHandler(c *gin.Context) {
	var gameForm GameForm
	err := c.ShouldBindJSON(&gameForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	secret := c.GetHeader("X-Secret")
	if secret != config.GetEnv("WEBHOOK_SECRET") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid secret"})
		return
	}
	err = DeleteWebhookGame(int(gameForm.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func GameWebhookUpdateHandler(c *gin.Context) {
	var gameForm GameForm
	err := c.ShouldBindJSON(&gameForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	secret := c.GetHeader("X-Secret")
	if secret != config.GetEnv("WEBHOOK_SECRET") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid secret"})
		return
	}
	err = UpdateWebhookGame(int(gameForm.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
