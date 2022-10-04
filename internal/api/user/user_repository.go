package user

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Password    string
	Username    string            `gorm:"unique"`
	UserProfile UserProfile       `gorm:"foreignkey:UserID"`
	Lists       []schema.GameList `gorm:"foreignkey:Owner;references:Username"`
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
	if err := database.DB.Preload(clause.Associations).Where("username = ?", username).First(&u).Error; err != nil {
		return nil
	}
	return &u
}
