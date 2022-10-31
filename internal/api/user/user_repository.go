package user

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	schema2 "github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email                 string `gorm:"unique"`
	Password              string
	Username              string                               `gorm:"unique"`
	UserProfile           UserProfile                          `gorm:"foreignkey:UserID"`
	Lists                 []schema2.GameList                   `gorm:"foreignkey:Owner;references:ID"`
	Reviews               []schema2.GameReview                 `gorm:"foreignkey:Creator;references:ID"`
	Helpfuls              []schema2.ReviewHelpful              `gorm:"foreignkey:UserID;references:ID"`
	Discussions           []schema2.GameDiscussion             `gorm:"foreignkey:CreatorID;references:ID"`
	Posts                 []schema2.DiscussionPost             `gorm:"foreignkey:CreatorID;references:ID"`
	CompletedAchievements []achievements.AchievementCompletion `gorm:"foreignkey:UserID;references:ID"`
}

type UserProfile struct {
	gorm.Model
	DisplayName string
	UserID      uint
	Avatar      *string
}

func SaveNewUser(u *User) {
	database.DB.Create(u)
}

func UpdateUser(u *User) {
	database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(u)
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

func GetByID(id uint) *User {
	var u User
	if err := database.DB.Preload("UserProfile").Preload("Lists").Where("id = ?", id).First(&u).Error; err != nil {
		return nil
	}
	return &u
}
