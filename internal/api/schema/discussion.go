package schema

import (
	"gorm.io/gorm"
	"time"
)

type GameDiscussion struct {
	ID        uint              `gorm:"primarykey" json:"id"`
	CreatedAt time.Time         `json:"-"`
	UpdatedAt time.Time         `json:"-"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"-"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Game      string            `json:"game"`
	Score     int               `json:"score"`
	CreatorID uint              `json:"-"`
	Scores    []DiscussionScore `gorm:"foreignKey:DiscussionID" json:"-"`
	Posts     []DiscussionPost  `gorm:"foreignKey:DiscussionID" json:"-"`
}

type DiscussionScore struct {
	DiscussionID uint
	UserID       uint
	Score        int
}

func (h *DiscussionScore) AfterCreate(tx *gorm.DB) error {
	err := tx.Model(&GameDiscussion{}).Where("id = ?", h.DiscussionID).Updates(map[string]any{
		"score": gorm.Expr("score + ?", h.Score),
	}).Error
	return err
}

func (h *DiscussionScore) AfterDelete(tx *gorm.DB) error {
	err := tx.Model(&GameDiscussion{}).Where("id = ?", h.DiscussionID).Updates(map[string]any{
		"score": gorm.Expr("score + ?", h.Score*-1),
	}).Error
	return err
}

type DiscussionPost struct {
	ID             uint             `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time        `json:"-"`
	UpdatedAt      time.Time        `json:"-"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`
	Body           string           `json:"body"`
	CreatorID      uint             `json:"-"`
	DiscussionID   uint             `json:"-"`
	Score          int              `json:"score"`
	Scores         []PostScore      `gorm:"foreignKey:PostID" json:"-"`
	OriginalPostID *uint            `json:"originalPostID"`
	Replies        []DiscussionPost `gorm:"foreignKey:OriginalPostID" json:"-"`
}

type PostScore struct {
	PostID uint
	UserID uint
	Score  int
}

func (h *PostScore) AfterCreate(tx *gorm.DB) error {
	err := tx.Model(&DiscussionPost{}).Where("id = ?", h.PostID).Updates(map[string]any{
		"score": gorm.Expr("score + ?", h.Score),
	}).Error
	return err
}

func (h *PostScore) AfterDelete(tx *gorm.DB) error {
	err := tx.Model(&DiscussionPost{}).Where("id = ?", h.PostID).Updates(map[string]any{
		"score": gorm.Expr("score + ?", h.Score*-1),
	}).Error
	return err
}
