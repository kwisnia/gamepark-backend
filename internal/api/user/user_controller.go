package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/pkg/crypto"
)

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterForm struct {
	LoginForm
	Username string `json:"username"`
}

// register user
func RegisterUser(c *gin.Context) {
	var form RegisterForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user := GetByEmail(form.Email)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email"})
		return
	}
	user = &User{
		Email:    form.Email,
		Password: crypto.HashAndSalt(form.Password),
		Username: form.Username,
		UserProfile: UserProfile{
			FirstName: nil,
			LastName:  nil,
		},
	}
	token, err := crypto.CreateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	Save(user)
	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusCreated, gin.H{"email": user.Email})
}

func LoginUser(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user := GetByEmail(form.Email)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email"})
		return
	}
	if !crypto.ComparePasswords(user.Password, form.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}
	token, err := crypto.CreateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{"email": user.Email})
}
