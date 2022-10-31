package dispatcher

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements/verifier"
)

func DispatchAchievementCheck(userID uint, achievementType achievements.ConditionType, count int64) {
	var newAchievements []achievements.Achievement
	var err error
	if achievementType == achievements.ConditionTypeProfileFirstEdit {
		newAchievements, err = verifier.VerifyProfileFirstEditAchievement(userID)
	} else {
		newAchievements, err = verifier.VerifyCountAchievements(userID, achievementType, count)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(newAchievements)
}
