package schema

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type GameCompletion struct {
	EnumCategory
}

type GameReview struct {
	ID               uint            `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time       `json:"-"`
	UpdatedAt        time.Time       `json:"-"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
	Rating           float64         `json:"rating"`
	HelpfulCount     int             `json:"helpfulCount"`
	PlatformID       *uint           `json:"-"`
	Platform         *Platform       `gorm:"foreignKey:PlatformID" json:"platform"`
	GameCompletionID uint            `json:"gameCompletionID"`
	GameCompletion   GameCompletion  `gorm:"foreignKey:GameCompletionID" json:"-"`
	Title            string          `json:"title"`
	Creator          uint            `json:"creator"`
	ContainsSpoilers bool            `json:"containsSpoilers"`
	Body             string          `json:"body"`
	Game             string          `json:"game"`
	Helpfuls         []ReviewHelpful `gorm:"foreignKey:ReviewID" json:"-"`
}

func (r *GameReview) AfterCreate(tx *gorm.DB) error {
	err := tx.Model(&Game{}).Where("slug = ?", r.Game).Updates(map[string]any{
		"rating_count": gorm.Expr("rating_count + ?", 1),
		"rating":       gorm.Expr("(rating_count * rating + ?) / (rating_count + 1)", r.Rating),
	}).Error
	return err
}

func (r *GameReview) AfterDelete(tx *gorm.DB) error {
	err := tx.Model(&Game{}).Where("slug = ?", r.Game).Updates(map[string]any{
		"rating_count": gorm.Expr("rating_count - ?", 1),
		"rating":       gorm.Expr("CASE WHEN rating_count > 1 THEN (rating_count * rating - ?) / (rating_count - 1) ELSE 0 END", r.Rating),
	}).Error
	return err
}

type ReviewHelpful struct {
	ReviewID uint
	UserID   uint
}

func (h *ReviewHelpful) AfterCreate(tx *gorm.DB) error {
	err := tx.Model(&GameReview{}).Where("id = ?", h.ReviewID).Updates(map[string]any{
		"helpful_count": gorm.Expr("helpful_count + ?", 1),
	}).Error
	return err
}

func (h *ReviewHelpful) AfterDelete(tx *gorm.DB) error {
	err := tx.Model(&GameReview{}).Where("id = ?", h.ReviewID).Updates(map[string]any{
		"helpful_count": gorm.Expr("helpful_count - ?", 1),
	}).Error
	fmt.Println(err)
	return err
}
