package schema

import (
	"gorm.io/gorm"
)

type ExternalCategory struct {
	EnumCategory
}

type ExternalGame struct {
	gorm.Model
	Category   ExternalCategory `gorm:"foreignKey:CategoryID"`
	CategoryID uint
	GameID     uint
	UID        string
	URL        string
}
