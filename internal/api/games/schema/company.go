package schema

import (
	"gorm.io/gorm"
	"time"
)

type CompanyLogo struct {
	gorm.Model
	Image
	CompanyID uint
}

type Company struct {
	gorm.Model
	ChangedCompanyID *uint
	ChangedCompany   *Company
	CompanyLogo      CompanyLogo
	Name             string
	Description      *string
	Slug             string
	StartDate        time.Time
	IGDBUrl          string
	GamesInvolved    []InvolvedCompany
}
