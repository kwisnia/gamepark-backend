package achievements

import (
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
)

func GetByID(id uint) (Achievement, error) {
	var achievement Achievement
	if err := database.DB.First(&achievement, id).Error; err != nil {
		return achievement, err
	}
	return achievement, nil
}

func GetNotCompletedAchievementsForType(userID uint, conditionType ConditionType) []Achievement {
	var achievements []Achievement

	if err := database.DB.Where("id NOT IN (?) AND conditionType = ?",
		database.DB.Table("achievement_completions").Select("achievement_id").Where("user_id = ?", userID), conditionType).
		Find(&achievements).Error; err != nil {
		return nil
	}
	return achievements
}

func CreateAchievementCompletion(userID uint, achievementID uint) error {
	achievementCompletion := AchievementCompletion{
		UserID:        userID,
		AchievementID: achievementID,
	}
	return database.DB.Create(&achievementCompletion).Error
}
