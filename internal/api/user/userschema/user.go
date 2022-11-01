package userschema

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email                 string `gorm:"unique"`
	Password              string
	Username              string                               `gorm:"unique"`
	UserProfile           UserProfile                          `gorm:"foreignkey:UserID"`
	Lists                 []schema.GameList                    `gorm:"foreignkey:Owner;references:ID"`
	Reviews               []schema.GameReview                  `gorm:"foreignkey:Creator;references:ID"`
	Helpfuls              []schema.ReviewHelpful               `gorm:"foreignkey:UserID;references:ID"`
	Discussions           []schema.GameDiscussion              `gorm:"foreignkey:CreatorID;references:ID"`
	Posts                 []schema.DiscussionPost              `gorm:"foreignkey:CreatorID;references:ID"`
	CompletedAchievements []achievements.AchievementCompletion `gorm:"foreignkey:UserID;references:ID"`
}

type UserProfile struct {
	gorm.Model
	DisplayName string
	UserID      uint
	Avatar      *string
}
