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

//go:generate gomodifytags -file $GOFILE -struct GameAgeRating -add-tags json -transform camelcase -w

type GameAgeRating struct {
	gorm.Model     `json:"-"`
	GameID         uint                  `json:"gameID"`
	AgeRatingID    uint                  `json:"ageRatingID"`
	AgeRating      AgeRating             `json:"ageRating"`
	OrganizationID uint                  `json:"organizationID"`
	Organization   AgeRatingOrganization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Synopsys       *string               `json:"synopsys"`
}
