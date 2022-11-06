package main

import (
	"github.com/joho/godotenv"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.SetupDB()
	db := database.DB
	newAchievements := []achievements.Achievement{
		{
			Name:          "First profile edit",
			Description:   "Edit your profile for the first time",
			Icon:          "https://gamepark-images.s3.eu-central-1.amazonaws.com/one-badge.svg",
			Points:        5,
			ConditionType: achievements.ConditionTypeProfileFirstEdit,
		},
		{
			Name:           "First review",
			Description:    "Write your first review",
			Icon:           "https://gamepark-images.s3.eu-central-1.amazonaws.com/star-badge.svg",
			Points:         10,
			ConditionType:  achievements.ConditionTypeReviews,
			ConditionValue: 1,
		},
		{
			Name:           "First discussion",
			Description:    "Create your first discussion",
			ConditionType:  achievements.ConditionTypeDiscussions,
			Icon:           "https://gamepark-images.s3.eu-central-1.amazonaws.com/bolt-badge.svg",
			Points:         10,
			ConditionValue: 1,
		},
		{
			Name:           "First helpful",
			Description:    "Mark someone's review as helpful for the first time",
			ConditionType:  achievements.ConditionTypeHelpfuls,
			Icon:           "https://gamepark-images.s3.eu-central-1.amazonaws.com/heart-badge.svg",
			Points:         5,
			ConditionValue: 1,
		},
		{
			ConditionType:  achievements.ConditionTypeLists,
			Name:           "First list",
			Description:    "Create your first list",
			Icon:           "https://gamepark-images.s3.eu-central-1.amazonaws.com/diamond-badge.svg",
			Points:         5,
			ConditionValue: 1,
		},
		{
			ConditionType:  achievements.ConditionTypePosts,
			Name:           "First post",
			Description:    "Write your first post",
			Icon:           "https://gamepark-images.s3.eu-central-1.amazonaws.com/like-badge.svg",
			Points:         5,
			ConditionValue: 1,
		},
	}
	db.Create(&newAchievements)
}
