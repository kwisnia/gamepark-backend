package chat

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/user/userschema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
)

func GetPageQuery(pageSize int, offset int) *gorm.DB {
	query := database.DB.Model(&userschema.Message{}).
		Limit(pageSize).Offset(offset).Order("created_at DESC")
	return query
}

func CreateMessage(message userschema.Message) error {
	return database.DB.Create(&message).Error
}

func GetMessagesBetweenUsers(pageSize int, offset int, userID1 uint, userID2 uint) ([]userschema.Message, error) {
	var messages []userschema.Message
	query := GetPageQuery(pageSize, offset)
	if err := query.Where("sender_id = ? AND receiver_id = ? OR sender_id = ? AND receiver_id = ?", userID1, userID2, userID2, userID1).
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func GetUniqueUserChatHistory(userID uint) ([]uint, error) {
	// get all users that have sent a message to the user or the user has sent a message to them
	var users []uint
	if err := database.DB.Raw("SELECT DISTINCT sender_id FROM messages WHERE receiver_id = ? UNION SELECT DISTINCT receiver_id FROM messages WHERE sender_id = ?", userID, userID).Scan(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
