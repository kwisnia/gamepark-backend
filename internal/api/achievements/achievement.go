package achievements

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
	"time"
)

type ConditionType string

const BannerUnlock = 10
const AnimatedAvatarUnlock = 20
const AnimatedBannerUnlock = 35

const (
	ConditionTypeReviews          ConditionType = "reviews"
	ConditionTypeDiscussions      ConditionType = "discussions"
	ConditionTypePosts            ConditionType = "posts"
	ConditionTypeLists            ConditionType = "lists"
	ConditionTypeHelpfuls         ConditionType = "helpfuls"
	ConditionTypeProfileFirstEdit ConditionType = "profile_first_edit"
)

type Achievement struct {
	ID             uint                    `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time               `json:"-"`
	UpdatedAt      time.Time               `json:"-"`
	DeletedAt      gorm.DeletedAt          `gorm:"index" json:"-"`
	Name           string                  `json:"name"`
	Description    string                  `json:"description"`
	Icon           string                  `json:"icon"`
	Points         int                     `json:"points"`
	ConditionType  ConditionType           `json:"-"`
	ConditionValue int                     `json:"-"`
	Completions    []AchievementCompletion `gorm:"foreignKey:AchievementID" json:"-"`
}

type AchievementCompletion struct {
	CreatedAt     time.Time `json:"completedAt"`
	UserID        uint      `json:"userID"`
	AchievementID uint      `json:"achievementID"`
}

func (ac *AchievementCompletion) AfterCreate(tx *gorm.DB) (err error) {
	var achievement Achievement
	if err := tx.First(&achievement, ac.AchievementID).Error; err != nil {
		return err
	}
	fmt.Println("To się robi?", ac.UserID)
	err = database.DB.Exec("UPDATE user_profiles SET user_score = user_score + ? WHERE user_id = ?", achievement.Points, ac.UserID).Error
	if err != nil {
		return err
	}
	// if new user score is higher than BannerUnlock, unlock banner
	database.DB.Exec("UPDATE user_feature_unlocks SET banner = (case when (select user_score from user_profiles where user_id = ?) >= ? then true else false end),"+
		"animated_avatar = (case when (select user_score from user_profiles where user_id = ?) >= ? then true else false end),"+
		"animated_banner = (case when (select user_score from user_profiles where user_id = ?) >= ? then true else false end) WHERE user_id = ?", ac.UserID, BannerUnlock, ac.UserID, AnimatedAvatarUnlock, ac.UserID, AnimatedBannerUnlock, ac.UserID)
	return nil
}
