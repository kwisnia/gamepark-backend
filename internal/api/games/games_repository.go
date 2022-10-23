package games

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type GameListElement struct {
	gorm.Model `json:"-"`
	ID         uint         `json:"id"`
	Slug       string       `json:"slug"`
	Name       string       `json:"name"`
	Cover      schema.Cover `gorm:"foreignKey:GameID" json:"cover,omitempty"`
}

func GetPage(pageSize int, offset int, filters []int, order string, search string) ([]GameListElement, error) {
	var games []GameListElement
	query := database.DB.Preload("Cover").Model(&schema.Game{}).
		Order(order)
	if strings.HasPrefix(order, "rating") {
		query = query.Order("rating_count DESC")
	} else {
		query = query.Order("aggregated_rating_count DESC")
	}
	query = query.Order("id asc").Limit(pageSize).Offset(offset)
	if len(filters) > 0 {
		query = query.Where("category_id IN ?", filters)
	}
	query = query.Where("LOWER(name) LIKE ?", "%"+search+"%")
	if err := query.Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}

func GetGameBySlug(slug string) (schema.Game, error) {
	game := schema.Game{}
	if err := database.DB.Preload("ExternalGames.Category").Preload("InvolvedCompanies.Company").Preload(clause.Associations).Where("slug = ?", slug).First(&game).Error; err != nil {
		return game, err
	}
	return game, nil
}

func GetGameById(id uint) (schema.Game, error) {
	game := schema.Game{}
	if err := database.DB.Preload(clause.Associations).Where("id = ?", id).First(&game).Error; err != nil {
		return game, err
	}
	return game, nil
}

//
//func UpdateGameRating(gameSlug string, newScore float64) {
//	database.DB.Model(&schema.Game{}).Where("slug = ?", gameSlug).Updates(
//		map[string]any{
//			"rating":       gorm.Expr("(rating_count * rating + ?) / (rating_count + 1)", newScore),
//			"rating_count": gorm.Expr("rating_count + ?", 1),
//		},
//	)
//}
