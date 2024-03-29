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
		userName, userID, valid := crypto.ValidateToken(authorizationHeader)
		if !valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			c.Set("userName", *userName)
			c.Set("userID", *userID)
			c.Next()
		}
	}
}

func AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.Next()
			return
		}
		userName, userID, valid := crypto.ValidateToken(authorizationHeader)
		if valid {
			c.Set("userName", *userName)
			c.Set("userID", *userID)
		}
		c.Next()
	}
}

func QueryAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Query("authorization")
		if authorizationHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userName, userID, valid := crypto.ValidateToken(authorizationHeader)
		if !valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			c.Set("userName", *userName)
			c.Set("userID", *userID)
			c.Next()
		}
	}
}
