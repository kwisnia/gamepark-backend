package schema

//go:generate gomodifytags -file $GOFILE -struct EnumCategory -add-tags json -transform camelcase -w

type EnumCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
