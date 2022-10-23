package reviews

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

func Save(r *schema.GameReview) error {
	if err := database.DB.Create(r).Error; err != nil {
		return err
	}
	return nil
}

func GetPageQuery(pageSize int, offset int, filters []int) *gorm.DB {
	fmt.Println(offset)
	query := database.DB.Preload("Platform").Model(&schema.GameReview{}).
		Limit(pageSize).Offset(offset).Order("created_at DESC")
	if len(filters) > 0 {
		query = query.Where("platform_id IN ?", filters)
	}
	return query
}

func GetByID(id uint) (*schema.GameReview, error) {
	var r schema.GameReview
	if err := database.DB.Where("id = ?", id).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func GetByGameSlug(slug string, pageSize int, offset int, filters []int) ([]schema.GameReview, error) {
	var r []schema.GameReview
	query := GetPageQuery(pageSize, offset, filters)
	if err := query.Where("game = ?", slug).Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func GetByUserUsername(username string, pageSize int, offset int, filters []int) ([]schema.GameReview, error) {
	var r []schema.GameReview
	query := GetPageQuery(pageSize, offset, filters)
	if err := query.Where("creator = ?", username).Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func GetByGameAndUser(gameSlug string, username string) (*schema.GameReview, error) {
	var r schema.GameReview
	if err := database.DB.Preload("Platform").Preload("GameCompletion").Where("game = ? AND creator = ?", gameSlug, username).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func CreateHelpful(r *schema.ReviewHelpful) {
	database.DB.Create(r)
}

func RemoveHelpful(r *schema.ReviewHelpful) {
	database.DB.Where("username = ? AND review_id = ?", r.Username, r.ReviewID).Delete(r)
}

func GetHelpfulByUserAndReview(username string, reviewID uint) (*schema.ReviewHelpful, error) {
	var r schema.ReviewHelpful
	if err := database.DB.Where("username = ? AND review_id = ?", username, reviewID).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func Delete(r *schema.GameReview) {
	database.DB.Delete(r)
}
