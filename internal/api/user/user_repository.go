package user

import (
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string
	Password    string
	Username    string
	UserProfile UserProfile `gorm:"foreignkey:UserID"`
}

type UserProfile struct {
	gorm.Model
	FirstName *string
	LastName  *string
	UserID    uint
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
