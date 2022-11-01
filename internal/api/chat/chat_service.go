package chat

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"github.com/kwisnia/inzynierka-backend/internal/api/user/userschema"
)

func GetChatHistory(requestingUser uint, secondUser string, pageSize int, page int) ([]userschema.Message, error) {
	secondParticipant := user.GetByUsername(secondUser)
	if secondParticipant == nil {
		return nil, fmt.Errorf("user not found")
	}
	offset := pageSize * (page - 1)
	return GetMessagesBetweenUsers(pageSize, offset, requestingUser, secondParticipant.ID)
}

func GetUsersChatReceivers(userID uint) ([]user.BasicUserDetails, error) {
	uniqueUsers, err := GetUniqueUserChatHistory(userID)
	fmt.Println("usery", len(uniqueUsers))
	if err != nil {
		return nil, err
	}
	var users []user.BasicUserDetails
	for _, uniqueUser := range uniqueUsers {
		userDetails := user.GetBasicUserDetailsByID(uniqueUser)
		if userDetails == nil {
			return nil, err
		}
		users = append(users, *userDetails)
	}
	return users, nil
}

func SaveNewMessage(creator uint, form MessageForm) error {
	msg := userschema.Message{
		SenderID:   creator,
		ReceiverID: form.Receiver,
		Content:    form.Content,
	}
	return CreateMessage(msg)
}
