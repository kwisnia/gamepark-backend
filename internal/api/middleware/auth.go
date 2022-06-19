package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/pkg/crypto"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if !crypto.ValidateToken(authorizationHeader) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			c.Next()
		}
	}
}
