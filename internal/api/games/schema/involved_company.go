package schema

import (
	"gorm.io/gorm"
)

type InvolvedCompany struct {
	gorm.Model
	GameID     uint
	CompanyID  uint
	Developer  bool
	Publisher  bool
	Porting    bool
	Supporting bool
}
