package discussions

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

func Save(r *schema.GameDiscussion) error {
	if err := database.DB.Create(r).Error; err != nil {
		return err
	}
	return nil
}

func GetPageQuery(pageSize int, offset int) *gorm.DB {
	fmt.Println(offset)
	query := database.DB.Preload("Platform").Model(&schema.GameReview{}).
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

func GetByGameSlug(slug string, pageSize int, offset int) ([]schema.GameDiscussion, error) {
	var discussions []schema.GameDiscussion
	query := GetPageQuery(pageSize, offset)
	if err := query.Where("game = ?", slug).Find(&discussions).Error; err != nil {
		return nil, err
	}
	return discussions, nil
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
