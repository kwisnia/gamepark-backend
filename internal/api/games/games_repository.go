package games

import (
	schema2 "github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type GameListElement struct {
	gorm.Model `json:"-"`
	ID         uint          `json:"id"`
	Slug       string        `json:"slug"`
	Name       string        `json:"name"`
	Cover      schema2.Cover `gorm:"foreignKey:GameID" json:"cover,omitempty"`
}

func GetPage(pageSize int, offset int, filters []int, order string, search string) ([]GameListElement, error) {
	var games []GameListElement
	query := database.DB.Preload("Cover").Model(&schema2.Game{}).
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

func GetGameBySlug(slug string) (schema2.Game, error) {
	game := schema2.Game{}
	if err := database.DB.Preload("ExternalGames.Category").Preload("InvolvedCompanies.Company").Preload(clause.Associations).Where("slug = ?", slug).First(&game).Error; err != nil {
		return game, err
	}
	return game, nil
}

func GetGameById(id uint) (schema2.Game, error) {
	game := schema2.Game{}
	if err := database.DB.Preload(clause.Associations).Where("id = ?", id).First(&game).Error; err != nil {
		return game, err
	}
	return game, nil
}

func GetGameShortInfoBySlug(slug string) (GameListElement, error) {
	game := GameListElement{}
	if err := database.DB.Model(&schema2.Game{}).Preload("Cover").Where("slug = ?", slug).First(&game).Error; err != nil {
		return game, err
	}
	return game, nil
}
