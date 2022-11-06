package verifier

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
)

func VerifyCountAchievements(userID uint, conditionType achievements.ConditionType, count int64) ([]achievements.Achievement, error) {
	notCompletedAchievements := achievements.GetNotCompletedAchievementsForType(userID, conditionType)
	fmt.Println(len(notCompletedAchievements))
	completedAchievements := make([]achievements.Achievement, 0)
	for _, achievement := range notCompletedAchievements {
		if count >= int64(achievement.ConditionValue) {
			err := achievements.CompleteAchievement(userID, achievement.ID)
			if err != nil {
				return nil, err
			}
			completedAchievements = append(completedAchievements, achievement)
		}
	}
	return completedAchievements, nil
}

func VerifyProfileFirstEditAchievement(userID uint) ([]achievements.Achievement, error) {
	notCompletedAchievements := achievements.GetNotCompletedAchievementsForType(userID, achievements.ConditionTypeProfileFirstEdit)
	completedAchievements := make([]achievements.Achievement, 0)
	for _, achievement := range notCompletedAchievements {
		err := achievements.CompleteAchievement(userID, achievement.ID)
		if err != nil {
			return nil, err
		}
		completedAchievements = append(completedAchievements, achievement)
	}
	return completedAchievements, nil
}
