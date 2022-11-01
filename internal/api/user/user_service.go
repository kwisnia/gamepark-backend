package user

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/uploader"
	"mime/multipart"
)

type BasicUserDetails struct {
	Username    string  `json:"username"`
	DisplayName string  `json:"displayName"`
	ID          uint    `json:"id"`
	Avatar      *string `json:"avatar"`
}

func GetUserDetails(userName string) *DetailsResponse {
	user := GetByUsername(userName)
	if user == nil {
		return nil
	}
	return &DetailsResponse{
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.UserProfile.DisplayName,
		ID:          user.ID,
		Lists:       user.Lists,
		Avatar:      user.UserProfile.Avatar,
	}
}

func GetBasicUserDetailsByID(userID uint) *BasicUserDetails {
	user := GetByID(userID)
	if user == nil {
		return nil
	}
	return &BasicUserDetails{
		Username:    user.Username,
		DisplayName: user.UserProfile.DisplayName,
		ID:          user.ID,
		Avatar:      user.UserProfile.Avatar,
	}
}

func UploadUserAvatar(userID uint, username string, file *multipart.FileHeader) error {
	user := GetByID(userID)
	if user == nil {
		return fmt.Errorf("user not found")
	}
	if user.Username != username {
		return fmt.Errorf("invalid permissions")
	}
	if file == nil {
		return fmt.Errorf("invalid file")
	}
	filePath, err := uploader.UploadFile("gamepark-images", *file)
	if err != nil {
		return err
	}
	fmt.Println("filePath", filePath)
	user.UserProfile.Avatar = &filePath
	UpdateUser(user)
	return nil
}

func GetUsers(pageSize int, page int, search string) ([]BasicUserDetails, error) {
	offset := pageSize * (page - 1)
	users, err := GetBySearch(pageSize, offset, search)
	if err != nil {
		return nil, err
	}
	var usersDetails = make([]BasicUserDetails, len(users))
	for i, user := range users {
		usersDetails[i] = BasicUserDetails{
			Username:    user.Username,
			DisplayName: user.UserProfile.DisplayName,
			ID:          user.ID,
			Avatar:      user.UserProfile.Avatar,
		}
	}
	return usersDetails, nil
}
