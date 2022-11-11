package activity

import (
	"encoding/json"
	"gorm.io/datatypes"
)

func CreateNewActivity(userID uint, activity ActivityType, activityData interface{}) error {
	activityJSON, err := json.Marshal(activityData)
	if err != nil {
		return err
	}
	newActivity := UserActivity{
		UserID:       userID,
		Activity:     activity,
		ActivityData: datatypes.JSON(activityJSON),
	}
	return CreateActivity(&newActivity)
}

func GetActivitiesForUser(userID uint, pageSize int, page int) ([]UserActivity, error) {
	offset := pageSize * (page - 1)
	return GetActivitiesByUser(userID, pageSize, offset)
}

func GetActivitiesForUsers(userIDs []uint, pageSize int, page int) ([]UserActivity, error) {
	offset := pageSize * (page - 1)
	return GetActivitiesByUsers(userIDs, pageSize, offset)
}
