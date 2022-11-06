package user

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/api/user/userschema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/uploader"
	"net/http"
)

type BasicUserDetails struct {
	Username       string                       `json:"username"`
	DisplayName    string                       `json:"displayName"`
	ID             uint                         `json:"id"`
	Avatar         *string                      `json:"avatar"`
	FollowerCount  uint                         `json:"followerCount"`
	FollowingCount uint                         `json:"followingCount"`
	Bio            string                       `json:"bio"`
	Banner         *string                      `json:"banner"`
	BannerPosition float32                      `json:"bannerPosition"`
	UserScore      int                          `json:"userScore"`
	UserUnlocks    userschema.UserFeatureUnlock `json:"userUnlocks"`
}

func GetUserDetails(userName string) *DetailsResponse {
	user := GetByUsername(userName)
	if user == nil {
		return nil
	}
	return &DetailsResponse{
		Email:          user.Email,
		Username:       user.Username,
		DisplayName:    user.UserProfile.DisplayName,
		ID:             user.ID,
		Lists:          user.Lists,
		Avatar:         user.UserProfile.Avatar,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		Bio:            user.UserProfile.Bio,
		Banner:         user.UserProfile.Banner,
		BannerPosition: user.UserProfile.BannerPosition,
		UserScore:      achievements.GetTotalScoreForUser(user.ID),
		UserUnlocks:    user.UserFeatureUnlock,
	}
}

func GetBasicUserDetailsByUsername(username string) *BasicUserDetails {
	user := GetByUsername(username)
	if user == nil {
		return nil
	}
	return &BasicUserDetails{
		Username:       user.Username,
		DisplayName:    user.UserProfile.DisplayName,
		ID:             user.ID,
		Avatar:         user.UserProfile.Avatar,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		Bio:            user.UserProfile.Bio,
		Banner:         user.UserProfile.Banner,
		BannerPosition: user.UserProfile.BannerPosition,
		UserScore:      achievements.GetTotalScoreForUser(user.ID),
		UserUnlocks:    user.UserFeatureUnlock,
	}
}

func GetBasicUserDetailsByID(userID uint) *BasicUserDetails {
	user := GetByID(userID)
	if user == nil {
		return nil
	}
	return &BasicUserDetails{
		Username:       user.Username,
		DisplayName:    user.UserProfile.DisplayName,
		ID:             user.ID,
		Avatar:         user.UserProfile.Avatar,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		Bio:            user.UserProfile.Bio,
		Banner:         user.UserProfile.Banner,
		BannerPosition: user.UserProfile.BannerPosition,
		UserScore:      achievements.GetTotalScoreForUser(user.ID),
		UserUnlocks:    user.UserFeatureUnlock,
	}
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
			Username:       user.Username,
			DisplayName:    user.UserProfile.DisplayName,
			ID:             user.ID,
			Avatar:         user.UserProfile.Avatar,
			FollowerCount:  user.FollowerCount,
			FollowingCount: user.FollowingCount,
			Bio:            user.UserProfile.Bio,
			Banner:         user.UserProfile.Banner,
			BannerPosition: user.UserProfile.BannerPosition,
			UserScore:      achievements.GetTotalScoreForUser(user.ID),
		}
	}
	return usersDetails, nil
}

func UpdateUserProfile(userID uint, username string, userProfileForm ProfileEditForm) error {
	user := GetByID(userID)
	if user == nil {
		return fmt.Errorf("user not found")
	}
	userUnlocks := user.UserFeatureUnlock
	if user.Username != username {
		return fmt.Errorf("invalid permissions")
	}
	user.UserProfile.DisplayName = userProfileForm.DisplayName
	user.UserProfile.Bio = userProfileForm.Bio
	if userProfileForm.Avatar != nil {
		validFile := true
		if !userUnlocks.AnimatedAvatar {
			fileContent, _ := userProfileForm.Avatar.Open()
			buf := make([]byte, 512)

			_, err := fileContent.Read(buf)
			if err != nil {
				return err
			}
			mimeType := http.DetectContentType(buf)
			// if mimeType is an animated image then set valid file as false
			if mimeType == "image/gif" {
				validFile = false
			}
			defer fileContent.Close()
		}
		if validFile {
			avatarFilePath, err := uploader.UploadFile("gamepark-images", *userProfileForm.Avatar)
			if err != nil {
				return err
			}
			user.UserProfile.Avatar = &avatarFilePath
		} else {
			return fmt.Errorf("invalid file type")
		}
	}
	if userProfileForm.Banner != nil {
		validFile := true
		if !userUnlocks.Banner {
			fileContent, _ := userProfileForm.Banner.Open()
			buf := make([]byte, 512)

			_, err := fileContent.Read(buf)
			if err != nil {
				return err
			}
			mimeType := http.DetectContentType(buf)
			if mimeType == "image/gif" {
				validFile = false
			}
			defer fileContent.Close()
		}
		if validFile {
			bannerFilePath, err := uploader.UploadFile("gamepark-images", *userProfileForm.Banner)
			if err != nil {
				return err
			}
			user.UserProfile.Banner = &bannerFilePath
		} else {
			return fmt.Errorf("invalid file type")
		}
	}
	if userProfileForm.RemoveBanner {
		user.UserProfile.Banner = nil
	}
	return UpdateUser(user)
}

func UpdateUserBannerPosition(userID uint, bannerPosition float32) error {
	user := GetByID(userID)
	if user == nil {
		return fmt.Errorf("user not found")
	}
	user.UserProfile.BannerPosition = bannerPosition
	return UpdateUser(user)
}

func GetUserAchievements(username string) ([]achievements.Achievement, error) {
	user := GetByUsername(username)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return achievements.GetAchievementsForUser(user.ID)
}
