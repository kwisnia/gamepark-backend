package schema

import (
	"gorm.io/gorm"
	"time"
)

//go:generate gomodifytags -file $GOFILE -struct Cover -add-tags json -transform camelcase -w

type GameCategory struct {
	EnumCategory
}

type Genre struct {
	EnumCategory
}

//go:generate gomodifytags -file $GOFILE -struct Game -add-tags json -transform camelcase -w

type Game struct {
	gorm.Model            `json:"-"`
	AgeRatings            []GameAgeRating   `gorm:"foreignKey:GameID" json:"ageRatings"`
	AggregatedRating      float64           `json:"aggregatedRating"`
	AggregatedRatingCount int               `json:"aggregatedRatingCount"`
	Artworks              []Artwork         `gorm:"foreignKey:GameID" json:"artworks"`
	Category              GameCategory      `gorm:"foreignKey:CategoryID" json:"category"`
	CategoryID            uint              `json:"-"`
	Cover                 Cover             `gorm:"foreignKey:GameID" json:"cover,omitempty"`
	DLCs                  []Game            `gorm:"foreignKey:DLCBaseReference" json:"dlcs"`
	DLCBaseReference      *uint             `json:"-"`
	ExpandedGames         []Game            `gorm:"foreignKey:ExpandedGameReference" json:"expandedGames"`
	ExpandedGameReference *uint             `json:"-"`
	Expansions            []Game            `gorm:"foreignKey:ExpansionReference" json:"expansions"`
	ExpansionReference    *uint             `json:"-"`
	ExternalGames         []ExternalGame    `gorm:"foreignKey:GameID" json:"externalGames"`
	FirstReleaseDate      time.Time         `json:"firstReleaseDate"`
	Genres                []Genre           `gorm:"many2many:game_genres" json:"genres"`
	InvolvedCompanies     []InvolvedCompany `gorm:"foreignKey:GameID" json:"involvedCompanies"`
	Name                  string            `json:"name"`
	ParentGame            *Game             `gorm:"foreignKey:ParentGameID" json:"parentGame"`
	ParentGameID          *uint             `json:"-"`
	Platforms             []Platform        `gorm:"many2many:game_platforms" json:"platforms"`
	Rating                float64           `json:"rating"`
	RatingCount           int               `json:"ratingCount"`
	ReleaseDates          []ReleaseDate     `gorm:"foreignKey:GameID" json:"releaseDates"`
	Remakes               []Game            `gorm:"foreignKey:RemakeBaseReference" json:"remakes"`
	RemakeBaseReference   *uint             `json:"-"`
	Remasters             []Game            `gorm:"foreignKey:RemasterBaseReference" json:"remasters"`
	RemasterBaseReference *uint             `json:"-"`
	Screenshots           []Screenshot      `gorm:"foreignKey:GameID" json:"screenshots"`
	SimilarGames          []Game            `gorm:"many2many:game_similar_games;" json:"similarGames"`
	Slug                  string            `json:"slug" gorm:"unique"`
	Storyline             string            `json:"storyline"`
	Summary               string            `json:"summary"`
	IGDBUrl               string            `json:"igdbUrl"`
	VersionParent         *Game             `json:"versionParent"`
	VersionParentID       *uint             `json:"-"`
	VersionTitle          string            `json:"versionTitle"`
	Videos                []GameVideo       `gorm:"foreignKey:GameID" json:"videos"`
	Lists                 []GameList        `gorm:"many2many:list_games" json:"lists"`
	Reviews               []GameReview      `gorm:"foreignKey:Game;references:Slug" json:"reviews"`
	Discussions           []GameDiscussion  `gorm:"foreignKey:Game;references:Slug" json:"discussions"`
}

type GameList struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	AvatarUrl   string         `json:"avatarUrl"`
	Owner       uint           `json:"owner"`
	Games       []Game         `gorm:"many2many:list_games;" json:"games"`
}
