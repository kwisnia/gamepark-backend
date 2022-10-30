package achievements

import (
	"gorm.io/gorm"
	"time"
)

type ConditionType string

const (
	ConditionTypeReviews     ConditionType = "reviews"
	ConditionTypeDiscussions ConditionType = "discussions"
	ConditionTypePosts       ConditionType = "posts"
	ConditionTypeLists       ConditionType = "lists"
	ConditionTypeHelpfuls    ConditionType = "helpfuls"
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
	ConditionType  ConditionType           `json:"conditionType"`
	ConditionValue int                     `json:"conditionValue"`
	Completions    []AchievementCompletion `gorm:"foreignKey:AchievementID" json:"-"`
}

type AchievementCompletion struct {
	CreatedAt     time.Time `json:"completedAt"`
	UserID        uint      `json:"userID"`
	AchievementID uint      `json:"achievementID"`
}
