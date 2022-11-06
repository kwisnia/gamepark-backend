package achievements

import (
	"encoding/json"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
)

func CompleteAchievement(userID uint, achievementID uint) error {
	_, err := GetByID(achievementID)
	if err != nil {
		return err
	}
	return CreateAchievementCompletion(userID, achievementID)
}

func GetTotalScoreForUser(userID uint) int {
	return getAchievementsScoreForUser(userID)
}

func PrepareWebSocketMessage(achievement Achievement) ([]byte, error) {
	receiverMessage := map[string]any{
		"badge":       achievement.Icon,
		"score":       achievement.Points,
		"title":       achievement.Name,
		"messageType": "newAchievement",
	}
	return json.Marshal(receiverMessage)
}

func GetAllAchievements() []Achievement {
	var achievements []Achievement
	if err := database.DB.Find(&achievements).Error; err != nil {
		return nil
	}
	return achievements
}

func GetAchievementsForUser(userID uint) ([]Achievement, error) {
	return getAchievementsForUser(userID)
}
