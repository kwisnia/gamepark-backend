package schema

import (
	"gorm.io/gorm"
)

//go:generate gomodifytags -file $GOFILE -struct InvolvedCompany -add-tags json -transform camelcase -w

type InvolvedCompany struct {
	gorm.Model `json:"-"`
	GameID     uint    `json:"gameID"`
	CompanyID  uint    `json:"companyID"`
	Company    Company `gorm:"foreignKey:CompanyID" json:"company"`
	Developer  bool    `json:"developer"`
	Publisher  bool    `json:"publisher"`
	Porting    bool    `json:"porting"`
	Supporting bool    `json:"supporting"`
}
