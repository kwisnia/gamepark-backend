package games

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm/clause"
	"time"

	schema "github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"gorm.io/gorm"
)

type Artwork struct {
	gorm.Model
	schema.Image
	GameID uint
}

type Screenshot struct {
	gorm.Model
	schema.Image
	GameID uint
}

//go:generate gomodifytags -file $GOFILE -struct Cover -add-tags json -w

type Cover struct {
	gorm.Model `json:"-"`
	schema.Image
	GameID uint `json:"-"`
}

type GameCategory struct {
	schema.EnumCategory
}

type Genre struct {
	schema.EnumCategory
}

type GameVideo struct {
	gorm.Model
	schema.Video
	GameID uint
}

//go:generate gomodifytags -file $GOFILE -struct Game -add-tags json -w

type Game struct {
	gorm.Model
	AgeRatings            []schema.GameAgeRating   `gorm:"foreignKey:GameID" json:"age_ratings"`
	AggregatedRating      float64                  `json:"aggregated_rating"`
	AggregatedRatingCount int                      `json:"aggregated_rating_count"`
	Artworks              []Artwork                `gorm:"foreignKey:GameID" json:"artworks"`
	Category              GameCategory             `gorm:"foreignKey:CategoryID" json:"category"`
	CategoryID            uint                     `json:"-"`
	Cover                 Cover                    `gorm:"foreignKey:GameID" json:"cover,omitempty"`
	DLCs                  []Game                   `gorm:"foreignKey:DLCBaseReference" json:"dlcs"`
	DLCBaseReference      *uint                    `json:"-"`
	ExpandedGames         []Game                   `gorm:"foreignKey:ExpandedGameReference" json:"expanded_games"`
	ExpandedGameReference *uint                    `json:"-"`
	Expansions            []Game                   `gorm:"foreignKey:ExpansionReference" json:"expansions"`
	ExpansionReference    *uint                    `json:"-"`
	ExternalGames         []schema.ExternalGame    `gorm:"foreignKey:GameID" json:"external_games"`
	FirstReleaseDate      time.Time                `json:"first_release_date"`
	Genres                []Genre                  `gorm:"many2many:game_genres" json:"genres"`
	InvolvedCompanies     []schema.InvolvedCompany `gorm:"foreignKey:GameID" json:"involved_companies"`
	Name                  string                   `json:"name"`
	ParentGame            *Game                    `gorm:"foreignKey:ParentGameID" json:"parent_game"`
	ParentGameID          *uint                    `json:"-"`
	Platforms             []schema.Platform        `gorm:"many2many:game_platforms" json:"platforms"`
	Rating                float64                  `json:"rating"`
	RatingCount           int                      `json:"rating_count"`
	ReleaseDates          []schema.ReleaseDate     `gorm:"foreignKey:GameID" json:"release_dates"`
	Remakes               []Game                   `gorm:"foreignKey:RemakeBaseReference" json:"remakes"`
	RemakeBaseReference   *uint                    `json:"-"`
	Remasters             []Game                   `gorm:"foreignKey:RemasterBaseReference" json:"remasters"`
	RemasterBaseReference *uint                    `json:"-"`
	Screenshots           []Screenshot             `gorm:"foreignKey:GameID" json:"screenshots"`
	SimilarGames          []Game                   `gorm:"many2many:game_similar_games;" json:"similar_games"`
	Slug                  string                   `json:"slug"`
	Storyline             string                   `json:"storyline"`
	Summary               string                   `json:"summary"`
	IGDBUrl               string                   `json:"igdb_url"`
	VersionParent         *Game                    `json:"version_parent"`
	VersionParentID       *uint                    `json:"-"`
	VersionTitle          string                   `json:"version_title"`
	Videos                []GameVideo              `gorm:"foreignKey:GameID" json:"videos"`
}

type GameListElement struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id"`
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Cover      Cover  `gorm:"foreignKey:GameID" json:"cover,omitempty"`
}

func GetPage(pageSize int, after int) ([]GameListElement, error) {
	var games []GameListElement
	fmt.Println(after)
	if err := database.DB.Preload("Cover").Model(&Game{}).
		Order("id").
		Limit(pageSize).Where("id > ?", after).Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}

func GetGameBySlug(slug string) (Game, error) {
	game := Game{}
	if err := database.DB.Preload(clause.Associations).Where("slug = ?", slug).First(&game).Error; err != nil {
		return game, err
	}
	return game, nil
}

func GetGameById(id uint) (Game, error) {
	game := Game{}
	if err := database.DB.Preload(clause.Associations).Where("id = ?", id).First(&game).Error; err != nil {
		return game, err
	}
	return game, nil
}
