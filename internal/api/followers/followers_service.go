package followers

import (
	"errors"
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"github.com/kwisnia/inzynierka-backend/internal/api/user/userschema"
	"gorm.io/gorm"
)

func FollowUser(userID uint, followUsername string) error {
	userCheck := user.GetByUsername(followUsername)
	if userCheck == nil {
		return fmt.Errorf("user not found")
	}
	followCheck, err := GetFollowConnection(userID, userCheck.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if followCheck != nil {
		return nil
	}
	err = CreateFollowing(&userschema.Following{
		UserID:   userID,
		Followed: userCheck.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func UnfollowUser(userID uint, followUsername string) error {
	userCheck := user.GetByUsername(followUsername)
	if userCheck == nil {
		return fmt.Errorf("user not found")
	}
	followCheck, err := GetFollowConnection(userID, userCheck.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if followCheck == nil {
		return nil
	}
	if err = DeleteFollowing(followCheck); err != nil {
		return err
	}
	return nil
}

func GetUserFollowers(username string, pageSize int, page int) ([]user.BasicUserDetails, error) {
	userCheck := user.GetByUsername(username)
	if userCheck == nil {
		return nil, fmt.Errorf("user not found")
	}
	offset := pageSize * (page - 1)
	followers, err := GetFollowersByUser(userCheck.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	users := make([]user.BasicUserDetails, len(followers))
	for i, follower := range followers {
		userDetails := user.GetBasicUserDetailsByID(follower.UserID)
		if userDetails == nil {
			return nil, fmt.Errorf("user not found")
		}
		users[i] = *userDetails
	}
	return users, nil
}

func GetUserFollowing(username string, pageSize int, page int) ([]user.BasicUserDetails, error) {
	userCheck := user.GetByUsername(username)
	if userCheck == nil {
		return nil, fmt.Errorf("user not found")
	}
	offset := pageSize * (page - 1)
	following, err := GetFollowingByUser(userCheck.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	users := make([]user.BasicUserDetails, len(following))
	for i, followingUser := range following {
		userDetails := user.GetBasicUserDetailsByID(followingUser.Followed)
		if userDetails == nil {
			return nil, fmt.Errorf("user not found")
		}
		users[i] = *userDetails
	}
	return users, nil
}

func CheckFollowConnection(userID uint, followUsername string) bool {
	userCheck := user.GetByUsername(followUsername)
	if userCheck == nil {
		return false
	}
	followCheck, err := GetFollowConnection(userID, userCheck.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if followCheck == nil {
		return false
	}
	return true
}

func GetAllFollowingsForUser(userID uint) ([]userschema.Following, error) {
	return GetAllFollowingsByUser(userID)
}
