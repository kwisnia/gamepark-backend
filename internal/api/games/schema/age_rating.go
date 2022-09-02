package schema

import (
	"gorm.io/gorm"
)

type AgeRatingOrganization struct {
	EnumCategory
}

type AgeRating struct {
	EnumCategory
}

type GameAgeRating struct {
	gorm.Model
	GameID         uint
	AgeRatingID    uint
	AgeRating      AgeRating
	OrganizationID uint
	Organization   AgeRatingOrganization `gorm:"foreignKey:OrganizationID"`
	Synopsys       *string
}
