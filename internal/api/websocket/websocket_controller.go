package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func WebSockerConnectionHandler(c *gin.Context) {
	fmt.Println("WebSockerConnectionHandler")
	userID := c.MustGet("userID").(uint)
	fmt.Println("connectionId", userID)
	err := CreateConnection(c.Writer, c.Request, userID)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}
}

func GetOnlineUsers(c *gin.Context) {
	fmt.Println("GetOnlineUsers")
	connections := GetConnections()
	fmt.Println("connections", connections)
	c.JSON(200, gin.H{"status": "success", "connections": connections})
}
