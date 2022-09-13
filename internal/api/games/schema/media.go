package schema

//go:generate gomodifytags -file $GOFILE -struct Image -add-tags json -w

type Image struct {
	Height  int    `json:"height"`
	Width   int    `json:"width"`
	ImageID string `json:"imageId"`
	URL     string `json:"url"`
}

//go:generate gomodifytags -file $GOFILE -struct Video -add-tags json -w

type Video struct {
	VideoID string `json:"videoId"`
	Name    string `json:"name"`
}
