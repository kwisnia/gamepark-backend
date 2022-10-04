package user

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
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
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type DetailsResponse struct {
	Email       string            `json:"email"`
	Username    string            `json:"username"`
	FirstName   *string           `json:"firstName"`
	LastName    *string           `json:"lastName"`
	DisplayName string            `json:"displayName"`
	ID          uint              `json:"id"`
	Lists       []schema.GameList `json:"lists"`
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
	user = GetByUsername(form.Username)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid username"})
		return
	}
	user = &User{
		Email:    form.Email,
		Password: crypto.HashAndSalt(form.Password),
		Username: form.Username,
		UserProfile: UserProfile{
			FirstName:   nil,
			LastName:    nil,
			DisplayName: form.DisplayName,
		},
	}
	token, err := crypto.CreateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	Save(user)
	c.Header("Authorization", token)
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
	token, err := crypto.CreateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"email": user.Email})
}

func GetDetails(c *gin.Context) {
	userName := c.GetString("userName")
	user := GetByUsername(userName)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, DetailsResponse{
		Email:       user.Email,
		Username:    user.Username,
		FirstName:   user.UserProfile.FirstName,
		LastName:    user.UserProfile.LastName,
		DisplayName: user.UserProfile.DisplayName,
		ID:          user.ID,
		Lists:       user.Lists,
	})
}

func GetDetailsByUsername(c *gin.Context) {
	userName := c.Param("userName")
	user := GetByUsername(userName)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, DetailsResponse{
		Email:       user.Email,
		Username:    user.Username,
		FirstName:   user.UserProfile.FirstName,
		LastName:    user.UserProfile.LastName,
		DisplayName: user.UserProfile.DisplayName,
		ID:          user.ID,
		Lists:       user.Lists,
	})
}
