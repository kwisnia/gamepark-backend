package followers

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/user/userschema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

func GetPageQuery(pageSize int, offset int) *gorm.DB {
	query := database.DB.Model(&userschema.Following{}).
		Limit(pageSize).Offset(offset)
	return query
}

func GetFollowersByUser(userID uint, pageSize int, offset int) ([]userschema.Following, error) {
	var followers []userschema.Following
	query := GetPageQuery(pageSize, offset)
	if err := query.Where("followed = ?", userID).Find(&followers).Error; err != nil {
		return nil, err
	}
	return followers, nil
}

func GetFollowingByUser(userID uint, pageSize int, offset int) ([]userschema.Following, error) {
	var following []userschema.Following
	query := GetPageQuery(pageSize, offset)
	if err := query.Where("user_id = ?", userID).Find(&following).Error; err != nil {
		return nil, err
	}
	return following, nil
}

func CreateFollowing(following *userschema.Following) error {
	return database.DB.Create(following).Error
}

func DeleteFollowing(following *userschema.Following) error {
	return database.DB.Where("user_id = ? AND followed = ?", following.UserID, following.Followed).Delete(following).Error
}

func GetFollowConnection(userID uint, followedID uint) (*userschema.Following, error) {
	var following userschema.Following
	if err := database.DB.Where("user_id = ? AND followed = ?", userID, followedID).First(&following).Error; err != nil {
		return nil, err
	}
	return &following, nil
}
