package schema

import (
	"gorm.io/gorm"
	"time"
)

//go:generate gomodifytags -file $GOFILE -struct CompanyLogo -add-tags json -w

type CompanyLogo struct {
	gorm.Model
	Image     `json:"image"`
	CompanyID uint `json:"company_id"`
}

//go:generate gomodifytags -file $GOFILE -struct Company -add-tags json -transform camelcase -w

type Company struct {
	gorm.Model       `json:"-"`
	ChangedCompanyID *uint             `json:"changedCompanyID"`
	ChangedCompany   *Company          `json:"changedCompany"`
	CompanyLogo      CompanyLogo       `json:"companyLogo"`
	Name             string            `json:"name"`
	Description      *string           `json:"description"`
	Slug             string            `json:"slug"`
	StartDate        time.Time         `json:"startDate"`
	IGDBUrl          string            `json:"igdbUrl"`
	GamesInvolved    []InvolvedCompany `json:"gamesInvolved"`
}
