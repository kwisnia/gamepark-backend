package main

import (
	games "github.com/kwisnia/inzynierka-backend/internal/api/games"
	gamesSchema "github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
)

func main() {
	config.LoadConfig()
	database.SetupDB()
	db := database.DB
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&user.UserProfile{})
	db.AutoMigrate(&gamesSchema.AgeRatingOrganization{})
	db.AutoMigrate(&gamesSchema.AgeRating{})
	db.AutoMigrate(&gamesSchema.GameAgeRating{})
	db.AutoMigrate(&gamesSchema.ExternalGame{})
	db.AutoMigrate(&gamesSchema.CompanyLogo{})
	db.AutoMigrate(&gamesSchema.Company{})
	db.AutoMigrate(&gamesSchema.InvolvedCompany{})
	db.AutoMigrate(&gamesSchema.ReleaseRegion{})
	db.AutoMigrate(&gamesSchema.ReleaseDateCategory{})
	db.AutoMigrate(&gamesSchema.ReleaseDate{})
	db.AutoMigrate(&gamesSchema.PlatformLogo{})
	db.AutoMigrate(&gamesSchema.Platform{})
	db.AutoMigrate(&games.Artwork{})
	db.AutoMigrate(&games.GameCategory{})
	db.AutoMigrate(&games.Cover{})
	db.AutoMigrate(&games.Screenshot{})
	db.AutoMigrate(&games.GameVideo{})
	db.AutoMigrate(&games.Genre{})
	db.AutoMigrate(&games.Game{})
}
