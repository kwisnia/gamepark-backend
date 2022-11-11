package discussions

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

type GameDiscussionListItem struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Title     string `json:"title"`
	Game      string `json:"game"`
	Score     int    `json:"score"`
	CreatorID uint   `json:"creator_id"`
}

func Save(r *schema.GameDiscussion) error {
	if err := database.DB.Create(r).Error; err != nil {
		return err
	}
	return nil
}

func GetPageQuery(pageSize int, offset int) *gorm.DB {
	query := database.DB.Model(&schema.GameDiscussion{}).
		Limit(pageSize).Offset(offset).Order("created_at ASC")

	return query
}

func GetByID(id uint) (*schema.GameDiscussion, error) {
	var r schema.GameDiscussion
	if err := database.DB.Where("id = ?", id).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func GetByGameSlug(slug string, pageSize int, offset int) ([]GameDiscussionListItem, error) {
	var discussions []GameDiscussionListItem
	query := GetPageQuery(pageSize, offset)
	if err := query.Where("game = ?", slug).Find(&discussions).Error; err != nil {
		return nil, err
	}
	return discussions, nil
}

func GetShortInfo(id uint) (*GameDiscussionListItem, error) {
	var discussion GameDiscussionListItem
	if err := database.DB.Model(&schema.GameDiscussion{}).Where("id = ?", id).First(&discussion).Error; err != nil {
		return nil, err
	}
	return &discussion, nil
}

func GetByUserID(userID uint, pageSize int, offset int) ([]schema.GameDiscussion, error) {
	var r []schema.GameDiscussion
	query := GetPageQuery(pageSize, offset)
	if err := query.Where("creator_id = ?", userID).Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func GetByGameAndUser(gameSlug string, userID uint) ([]schema.GameDiscussion, error) {
	var discussions []schema.GameDiscussion
	if err := database.DB.Where("game = ? AND creator_id = ?", gameSlug, userID).Find(&discussions).Error; err != nil {
		return nil, err
	}
	return discussions, nil
}

func CreateDiscussionScore(s *schema.DiscussionScore) error {
	return database.DB.Create(s).Error
}

func RemoveDiscussionScore(s *schema.DiscussionScore) error {
	return database.DB.Where("user_id = ? AND discussion_id = ?", s.UserID, s.DiscussionID).Delete(s).Error
}

func GetScoreByUserAndDiscussion(userID uint, discussionID uint) (*schema.DiscussionScore, error) {
	var discussion schema.DiscussionScore
	if err := database.DB.Where("user_id = ? AND discussion_id = ?", userID, discussionID).First(&discussion).Error; err != nil {
		return nil, err
	}
	return &discussion, nil
}

func Delete(r *schema.GameDiscussion) error {
	return database.DB.Delete(r).Error
}

func CountByUser(userID uint) (int64, error) {
	var count int64
	if err := database.DB.Model(&schema.GameDiscussion{}).Where("creator_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
