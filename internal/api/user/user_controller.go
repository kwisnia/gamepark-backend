package user

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user/userschema"
	"net/http"
	"strconv"

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
	DisplayName string            `json:"displayName"`
	ID          uint              `json:"id"`
	Lists       []schema.GameList `json:"lists"`
	Avatar      *string           `json:"avatar"`
}

// register user
func RegisterUserHandler(c *gin.Context) {
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
	user = &userschema.User{
		Email:    form.Email,
		Password: crypto.HashAndSalt(form.Password),
		Username: form.Username,
		UserProfile: userschema.UserProfile{
			DisplayName: form.DisplayName,
			Avatar:      nil,
		},
	}
	SaveNewUser(user)
	token, err := crypto.CreateToken(user.Username, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusCreated, gin.H{"email": user.Email})
}

func LoginUserHandler(c *gin.Context) {
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
	token, err := crypto.CreateToken(user.Username, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"email": user.Email})
}

func GetDetailsHandler(c *gin.Context) {
	userName := c.GetString("userName")
	user := GetUserDetails(userName)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetDetailsByUsernameHandler(c *gin.Context) {
	userName := c.Param("userName")
	user := GetUserDetails(userName)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	user.Email = ""
	c.JSON(http.StatusOK, user)
}

func UploadUserAvatarHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userName := c.GetString("userName")
	userID := c.GetUint("userID")
	user := GetByUsername(userName)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	err = UploadUserAvatar(userID, userName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded"})
}

func GetUsersHandler(c *gin.Context) {
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
	search := c.DefaultQuery("search", "")
	users, err := GetUsers(pageSize, page, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
