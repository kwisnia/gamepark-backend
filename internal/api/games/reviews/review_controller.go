package reviews

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReviewForm struct {
	Rating           float64 `json:"rating"`
	Title            string  `json:"title"`
	Body             string  `json:"body"`
	ContainsSpoilers bool    `json:"containsSpoilers"`
	PlatformID       *uint   `json:"platform"`
	GameCompletionID uint    `json:"completionStatus"`
}

func CreateReviewHandler(c *gin.Context) {
	var form ReviewForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(form.PlatformID)
	username := c.GetString("userName")
	gameSlug := c.Param("slug")
	review, err := CreateReview(username, gameSlug, form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, review)
}

func GetReviewHandler(c *gin.Context) {
	username := c.GetString("userName")
	reviewId, err := strconv.Atoi(c.Param("reviewId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid review id"})
		return
	}
	review, err := GetReview(uint(reviewId), username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, review)
}

func DeleteReviewHandler(c *gin.Context) {
	username := c.GetString("userName")
	slug := c.Param("slug")
	err := DeleteReview(username, slug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Review deleted"})
}

func GetReviewsForGameHandler(c *gin.Context) {
	gameSlug := c.Param("slug")
	userName := c.GetString("userName")
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
	for i, filter := range filters {
		parsedFilters[i], err = strconv.Atoi(filter)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid filter"})
			return
		}
	}
	reviews, err := GetReviewsForGame(pageSize, page, parsedFilters, gameSlug, userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func GetReviewsForUserHandler(c *gin.Context) {
	username := c.GetString("userName")
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
	for i, filter := range filters {
		parsedFilters[i], err = strconv.Atoi(filter)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid filter"})
			return
		}
	}
	reviews, err := GetReviewsForUser(pageSize, page, parsedFilters, username)
	c.JSON(http.StatusOK, reviews)
}

func MarkReviewAsHelpfulHandler(c *gin.Context) {
	username := c.GetString("userName")
	reviewId, err := strconv.Atoi(c.Param("reviewId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid review id"})
		return
	}
	err = MarkReviewAsHelpful(username, uint(reviewId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Review marked as helpful"})
}

func UnmarkReviewAsHelpfulHandler(c *gin.Context) {
	username := c.GetString("userName")
	reviewId, err := strconv.Atoi(c.Param("reviewId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid review id"})
		return
	}
	err = UnmarkReviewAsHelpful(username, uint(reviewId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Review marked as unhelpful"})
}
