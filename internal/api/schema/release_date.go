package schema

import (
	"gorm.io/gorm"
	"time"
)

type ReleaseRegion struct {
	EnumCategory
}

type ReleaseDateCategory struct {
	EnumCategory
}

//go:generate gomodifytags -file $GOFILE -struct ReleaseDate -add-tags json -transform camelcase -w

type ReleaseDate struct {
	gorm.Model `json:"-"`
	GameID     uint                `json:"gameID"`
	Region     ReleaseRegion       `json:"region"`
	RegionID   uint                `json:"regionID"`
	Human      string              `json:"human"`
	Date       time.Time           `json:"date"`
	PlatformID uint                `json:"platformID"`
	Category   ReleaseDateCategory `json:"category"`
	CategoryID uint                `json:"categoryID"`
}
