package schema

import "gorm.io/gorm"

//go:generate gomodifytags -file $GOFILE -struct PlatformLogo -add-tags json -w

type PlatformLogo struct {
	gorm.Model `json:"-"`
	Image      `json:"image"`
	PlatformID uint `json:"platform_id"`
}

//go:generate gomodifytags -file $GOFILE -struct Platform -add-tags json -transform camelcase -w

type Platform struct {
	gorm.Model   `json:"-"`
	Name         string       `json:"name"`
	Abbreviation string       `json:"abbreviation"`
	Generation   int          `json:"generation"`
	Logo         PlatformLogo `gorm:"foreignKey:PlatformID" json:"logo"`
	Slug         string       `json:"slug"`
	IGDBUrl      string       `json:"igdbUrl"`
}
