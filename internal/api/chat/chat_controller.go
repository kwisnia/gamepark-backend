package chat

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MessageForm struct {
	Receiver uint   `json:"receiver"`
	Content  string `json:"content"`
}

func GetChatHistoryHandler(c *gin.Context) {
	gameSlug := c.Param("user")
	userID := c.GetUint("userID")
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
	messages, err := GetChatHistory(userID, gameSlug, pageSize, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func GetChatReceiversHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	receivers, err := GetUsersChatReceivers(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, receivers)
}
