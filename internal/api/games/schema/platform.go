package schema

import "gorm.io/gorm"

type PlatformLogo struct {
	gorm.Model
	Image
	PlatformID uint
}

type Platform struct {
	gorm.Model
	Name         string
	Abbreviation string
	Generation   int
	Logo         PlatformLogo `gorm:"foreignKey:PlatformID"`
	Slug         string
	IGDBUrl      string
}
