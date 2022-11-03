package followers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FollowUserHandler(c *gin.Context) {
	username := c.Param("userName")
	userID := c.MustGet("userID").(uint)
	err := FollowUser(userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func UnfollowUserHandler(c *gin.Context) {
	username := c.Param("userName")
	userID := c.MustGet("userID").(uint)
	err := UnfollowUser(userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func GetUserFollowersHandler(c *gin.Context) {
	username := c.Param("userName")
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
	followers, err := GetUserFollowers(username, pageSize, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, followers)
}

func GetUserFollowingHandler(c *gin.Context) {
	username := c.Param("userName")
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
	following, err := GetUserFollowing(username, pageSize, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, following)
}

func CheckFollowConnectionHandler(c *gin.Context) {
	username := c.Param("userName")
	userID := c.MustGet("userID").(uint)
	following := CheckFollowConnection(userID, username)
	c.JSON(http.StatusOK, following)
}
