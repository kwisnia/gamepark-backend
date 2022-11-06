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
