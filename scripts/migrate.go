package main

import (
	"github.com/joho/godotenv"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user/userschema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"log"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.SetupDB()
	db := database.DB
	db.AutoMigrate(&userschema.User{})
	db.AutoMigrate(&userschema.UserProfile{})
	db.AutoMigrate(&schema.AgeRatingOrganization{})
	db.AutoMigrate(&schema.AgeRating{})
	db.AutoMigrate(&schema.GameAgeRating{})
	db.AutoMigrate(&schema.ExternalGame{})
	db.AutoMigrate(&schema.CompanyLogo{})
	db.AutoMigrate(&schema.Company{})
	db.AutoMigrate(&schema.InvolvedCompany{})
	db.AutoMigrate(&schema.ReleaseRegion{})
	db.AutoMigrate(&schema.ReleaseDateCategory{})
	db.AutoMigrate(&schema.ReleaseDate{})
	db.AutoMigrate(&schema.PlatformLogo{})
	db.AutoMigrate(&schema.Platform{})
	db.AutoMigrate(&schema.Artwork{})
	db.AutoMigrate(&schema.GameCategory{})
	db.AutoMigrate(&schema.Cover{})
	db.AutoMigrate(&schema.Screenshot{})
	db.AutoMigrate(&schema.GameVideo{})
	db.AutoMigrate(&schema.Genre{})
	db.AutoMigrate(&schema.Game{})
	db.AutoMigrate(&schema.GameList{})
	db.AutoMigrate(&schema.GameCompletion{})
	db.AutoMigrate(&schema.GameReview{})
	db.AutoMigrate(&schema.ReviewHelpful{})
	db.AutoMigrate(&schema.GameDiscussion{})
	db.AutoMigrate(&schema.DiscussionScore{})
	db.AutoMigrate(&schema.DiscussionPost{})
	db.AutoMigrate(&schema.PostScore{})
	db.AutoMigrate(&achievements.Achievement{})
	db.AutoMigrate(&achievements.AchievementCompletion{})
	db.AutoMigrate(&userschema.Message{})
	db.AutoMigrate(&userschema.Following{})
}
