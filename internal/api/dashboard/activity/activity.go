package activity

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ActivityType string

const (
	NewReview     ActivityType = "new_review"
	NewPost       ActivityType = "new_post"
	NewDiscussion ActivityType = "new_discussion"
)

type ReviewActivityData struct {
	ReviewID uint `json:"reviewID"`
}

type PostActivityData struct {
	PostID       uint `json:"postID"`
	DiscussionID uint `json:"discussionID"`
}

type DiscussionActivityData struct {
	DiscussionID uint `json:"discussionID"`
}

type UserActivity struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `json:"userID"`
	Activity     ActivityType   `json:"activity"`
	ActivityData datatypes.JSON `json:"activityData"`
}
