package user

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user/userschema"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/pkg/crypto"
)

type LoginForm struct {
	Username string `json:"username" binding:"required" binding:"min=3" binding:"max=30"`
	Password string `json:"password" binding:"required" binding:"min=8" binding:"max=50"`
}

type RegisterForm struct {
	LoginForm
	DisplayName string `json:"displayName" binding:"required" binding:"min=3" binding:"max=30"`
}

type ProfileEditForm struct {
	DisplayName  string                `binding:"required" form:"displayName"`
	Bio          string                `form:"bio"`
	Banner       *multipart.FileHeader `form:"banner"`
	Avatar       *multipart.FileHeader `form:"avatar"`
	RemoveBanner bool                  `form:"removeBanner"`
}

type BannerPositionForm struct {
	Position *float32 `json:"position" binding:"required" binding:"min=0" binding:"max=100"`
}

type DetailsResponse struct {
	Username       string                       `json:"username"`
	DisplayName    string                       `json:"displayName"`
	ID             uint                         `json:"id"`
	Lists          []schema.GameList            `json:"lists"`
	Avatar         *string                      `json:"avatar"`
	FollowerCount  uint                         `json:"followerCount"`
	Bio            string                       `json:"bio"`
	FollowingCount uint                         `json:"followingCount"`
	Banner         *string                      `json:"banner"`
	BannerPosition float32                      `json:"bannerPosition"`
	UserScore      int                          `json:"userScore"`
	UserUnlocks    userschema.UserFeatureUnlock `json:"userUnlocks"`
}

// register user
func RegisterUserHandler(c *gin.Context) {
	var form RegisterForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user := GetByUsername(form.Username)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid username"})
		return
	}
	user = &userschema.User{
		Password: crypto.HashAndSalt(form.Password),
		Username: form.Username,
		UserProfile: userschema.UserProfile{
			DisplayName: form.DisplayName,
			Avatar:      nil,
		},
		UserFeatureUnlock: userschema.UserFeatureUnlock{
			Banner:         false,
			AnimatedAvatar: false,
			AnimatedBanner: false,
			Avatar:         true,
		},
	}
	SaveNewUser(user)
	token, err := crypto.CreateToken(user.Username, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Header("Authorization", token)
	c.JSON(http.StatusCreated, gin.H{"username": user.Username})
}

func LoginUserHandler(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user := GetByUsername(strings.ToLower(form.Username))
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username"})
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
	c.JSON(http.StatusOK, gin.H{"username": user.Username})
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
	user := GetBasicUserDetailsByUsername(userName)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, user)
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

func UpdateUserProfileHandler(c *gin.Context) {
	var form ProfileEditForm
	if err := c.ShouldBindWith(&form, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(form)
	userName := c.GetString("userName")
	userID := c.GetUint("userID")
	user := GetByUsername(userName)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	err := UpdateUserProfile(userID, userName, form)
	if err != nil {
		if err.Error() == "invalid file type" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file type"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func UpdateUserBannerPositionHandler(c *gin.Context) {
	var form BannerPositionForm
	if err := c.ShouldBindJSON(&form); err != nil {
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
	err := UpdateUserBannerPosition(userID, *form.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Banner position updated"})
}

func GetUserAchievementsHandler(c *gin.Context) {
	userName := c.Param("userName")
	achievements, err := GetUserAchievements(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, achievements)
}
