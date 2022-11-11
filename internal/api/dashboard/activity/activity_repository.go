package activity

import (
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

func CreateActivity(activity *UserActivity) error {
	return database.DB.Create(activity).Error
}

func GetPageQuery(pageSize int, offset int) *gorm.DB {
	query := database.DB.Model(&UserActivity{}).
		Limit(pageSize).Offset(offset)
	return query
}

func GetActivitiesByUser(userID uint, pageSize int, offset int) ([]UserActivity, error) {
	var activities []UserActivity
	query := GetPageQuery(pageSize, offset)
	if err := query.Order("created_at desc").Where("user_id = ?", userID).Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func GetActivitiesByUsers(userIDs []uint, pageSize int, offset int) ([]UserActivity, error) {
	var activities []UserActivity
	query := GetPageQuery(pageSize, offset)
	if err := query.Order("created_at desc").Where("user_id IN ?", userIDs).Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}
