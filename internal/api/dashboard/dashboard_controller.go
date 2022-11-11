package dashboard

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetActivitiesFromFollowedHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
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
	activities, err := GetNewestActivitiesFromFollowedUsers(userID, pageSize, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}
