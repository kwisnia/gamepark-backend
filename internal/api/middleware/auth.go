package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/pkg/crypto"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userName, valid := crypto.ValidateToken(authorizationHeader)
		if !valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			c.Set("userName", *userName)
			c.Next()
		}
	}
}
