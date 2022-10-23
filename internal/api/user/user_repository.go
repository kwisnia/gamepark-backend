package user

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Password    string
	Username    string                 `gorm:"unique"`
	UserProfile UserProfile            `gorm:"foreignkey:UserID"`
	Lists       []schema.GameList      `gorm:"foreignkey:Owner;references:Username"`
	Reviews     []schema.GameReview    `gorm:"foreignkey:Creator;references:Username"`
	Helpfuls    []schema.ReviewHelpful `gorm:"foreignkey:Username;references:Username"`
}

type UserProfile struct {
	gorm.Model
	DisplayName string
	FirstName   *string
	LastName    *string
	UserID      uint
}

func Save(u *User) {
	database.DB.Create(u)
}

func GetByEmail(email string) *User {
	var u User
	if err := database.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return nil
	}
	return &u
}

func GetByUsername(username string) *User {
	var u User
	if err := database.DB.Preload("UserProfile").Preload("Lists").Where("username = ?", username).First(&u).Error; err != nil {
		return nil
	}
	return &u
}
