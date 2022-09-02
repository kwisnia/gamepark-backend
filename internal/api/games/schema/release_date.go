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

type ReleaseDate struct {
	gorm.Model
	GameID     uint
	Region     ReleaseRegion
	RegionID   uint
	Human      string
	Date       time.Time
	PlatformID uint
	Category   ReleaseDateCategory
	CategoryID uint
}
