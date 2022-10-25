package posts

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

func Save(r *schema.DiscussionPost) error {
	if err := database.DB.Create(r).Error; err != nil {
		return err
	}
	return nil
}

func GetPageQuery(pageSize int, offset int) *gorm.DB {
	query := database.DB.Preload("Platform").Model(&schema.GameReview{}).
		Limit(pageSize).Offset(offset).Order("created_at ASC")

	return query
}

func GetByID(id uint) (*schema.DiscussionPost, error) {
	var r schema.DiscussionPost
	if err := database.DB.Where("id = ?", id).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func GetByDiscussionID(discussionID uint, pageSize int, offset int) ([]schema.DiscussionPost, error) {
	var posts []schema.DiscussionPost
	query := GetPageQuery(pageSize, offset)
	if err := query.Where("discussion_id = ?", discussionID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func GetByUserID(userID uint, pageSize int, offset int) ([]schema.DiscussionPost, error) {
	var r []schema.DiscussionPost
	query := GetPageQuery(pageSize, offset)
	if err := query.Where("creator_id = ?", userID).Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func CreatePostScore(s *schema.PostScore) error {
	return database.DB.Create(s).Error
}

func DeleteScore(s *schema.PostScore) error {
	return database.DB.Where("user_id = ? AND post_id = ?", s.UserID, s.PostID).Delete(s).Error
}

func GetScoreByUserAndPost(userID uint, postID uint) (*schema.PostScore, error) {
	var postScore schema.PostScore
	if err := database.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&postScore).Error; err != nil {
		return nil, err
	}
	return &postScore, nil
}

func Delete(r *schema.DiscussionPost) error {
	return database.DB.Delete(r).Error
}
