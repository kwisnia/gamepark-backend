package schema

import "gorm.io/gorm"

//go:generate gomodifytags -file $GOFILE -struct Image -add-tags json -transform camelcase -w

type Image struct {
	Height  int    `json:"height"`
	Width   int    `json:"width"`
	ImageID string `json:"imageId"`
	URL     string `json:"url"`
}

//go:generate gomodifytags -file $GOFILE -struct Video -add-tags json -transform camelcase -w

type Video struct {
	VideoID string `json:"videoId"`
	Name    string `json:"name"`
}

type Artwork struct {
	gorm.Model `json:"-"`
	Image
	GameID uint
}

type Screenshot struct {
	gorm.Model `json:"-"`
	Image
	GameID uint
}

type GameVideo struct {
	gorm.Model `json:"-"`
	Video
	GameID uint
}

type Cover struct {
	gorm.Model `json:"-"`
	Image
	GameID uint `json:"-"`
}
