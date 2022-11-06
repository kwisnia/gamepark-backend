package dispatcher

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements/verifier"
	"github.com/kwisnia/inzynierka-backend/internal/api/websocket"
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
	for _, achievement := range newAchievements {
		message, err := achievements.PrepareWebSocketMessage(achievement)
		if err != nil {
			fmt.Println(err)
			return
		}
		websocket.ClientHub.Send <- &websocket.Message{
			ReceiverID: userID,
			Data:       message,
		}
	}

}
