package games

import (
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

type Cover struct {
	gorm.Model
	schema.Image
	GameID uint
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

type Game struct {
	gorm.Model
	AgeRatings            []schema.GameAgeRating `gorm:"foreignKey:GameID"`
	AggregatedRating      float64
	AggregatedRatingCount int
	Artworks              []Artwork    `gorm:"foreignKey:GameID"`
	Category              GameCategory `gorm:"foreignKey:CategoryID"`
	CategoryID            uint
	Cover                 Cover `gorm:"foreignKey:CoverID"`
	CoverID               uint
	DLCs                  []Game `gorm:"foreignKey:DLCBaseReference"`
	DLCBaseReference      *uint
	ExpandedGames         []Game `gorm:"foreignKey:ExpandedGameReference"`
	ExpandedGameReference *uint
	Expansions            []Game `gorm:"foreignKey:ExpansionReference"`
	ExpansionReference    *uint
	ExternalGames         []schema.ExternalGame `gorm:"foreignKey:GameID"`
	FirstReleaseDate      time.Time
	Genres                []Genre                  `gorm:"foreignKey:ID"`
	InvolvedCompanies     []schema.InvolvedCompany `gorm:"foreignKey:GameID"`
	Name                  string
	ParentGame            *Game `gorm:"foreignKey:ParentGameID"`
	ParentGameID          *uint
	Rating                float64
	RatingCount           int
	Remakes               []Game `gorm:"foreignKey:RemakeBaseReference"`
	RemakeBaseReference   *uint
	Remasters             []Game `gorm:"foreignKey:RemasterBaseReference"`
	RemasterBaseReference *uint
	Screenshots           []Screenshot `gorm:"foreignKey:GameID"`
	SimilarGames          []Game       `gorm:"many2many:game_similar_games;"`
	Slug                  string
	Storyline             string
	Summary               string
	IGDBUrl               string
	VersionParent         *Game
	VersionParentID       *uint
	VersionTitle          string
	Videos                []GameVideo `gorm:"foreignKey:GameID"`
}
