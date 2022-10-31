package schema

import (
	"gorm.io/gorm"
)

type ExternalCategory struct {
	EnumCategory
}

//go:generate gomodifytags -file $GOFILE -struct ExternalGame -add-tags json -transform camelcase -w

type ExternalGame struct {
	gorm.Model `json:"-"`
	Category   ExternalCategory `gorm:"foreignKey:CategoryID" json:"category"`
	CategoryID uint             `json:"categoryID"`
	GameID     uint             `json:"gameID"`
	UID        string           `json:"uid"`
	URL        string           `json:"url"`
}
