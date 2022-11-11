package dashboard

import (
	"encoding/json"
	"errors"
	"github.com/kwisnia/inzynierka-backend/internal/api/dashboard/activity"
	"github.com/kwisnia/inzynierka-backend/internal/api/followers"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/discussions"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/discussions/posts"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/reviews"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"gorm.io/gorm"
	"time"
)

type Activity struct {
	Type      activity.ActivityType `json:"type"`
	User      user.BasicUserDetails `json:"user"`
	Data      any                   `json:"data"`
	CreatedAt time.Time             `json:"createdAt"`
}

func GetNewestActivitiesFromFollowedUsers(userID uint, pageSize int, page int) ([]Activity, error) {
	followedUsers, err := followers.GetAllFollowingsForUser(userID)
	if err != nil {
		return nil, err
	}
	var followedUsersIDs []uint
	for _, followedUser := range followedUsers {
		followedUsersIDs = append(followedUsersIDs, followedUser.Followed)
	}
	activities, err := activity.GetActivitiesForUsers(followedUsersIDs, pageSize, page)
	if err != nil {
		return nil, err
	}
	parsedActivities := make([]Activity, 0)
	for _, userActivity := range activities {
		userDetails := user.GetBasicUserDetailsByID(userActivity.UserID)
		if userDetails == nil {
			continue
		}
		switch userActivity.Activity {
		case activity.NewReview:
			parsedActivity, err := GetReviewActivity(userActivity)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
				continue
			}
			parsedActivity.User = *userDetails
			parsedActivity.CreatedAt = userActivity.CreatedAt
			parsedActivities = append(parsedActivities, *parsedActivity)
		case activity.NewDiscussion:
			parsedActivity, err := GetDiscussionActivity(userActivity, userID)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
				continue
			}
			parsedActivity.User = *userDetails
			parsedActivity.CreatedAt = userActivity.CreatedAt
			parsedActivities = append(parsedActivities, *parsedActivity)
		case activity.NewPost:
			parsedActivity, err := GetPostActivity(userActivity, userID)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
				continue
			}
			parsedActivity.User = *userDetails
			parsedActivity.CreatedAt = userActivity.CreatedAt
			parsedActivities = append(parsedActivities, *parsedActivity)
		}
	}
	return parsedActivities, nil
}

func GetReviewActivity(userActivity activity.UserActivity) (*Activity, error) {
	reviewActivityData := activity.ReviewActivityData{}
	if err := json.Unmarshal(userActivity.ActivityData, &reviewActivityData); err != nil {
		return nil, err
	}
	review, err := reviews.GetReviewWithGameDetailsById(reviewActivityData.ReviewID)
	if err != nil {
		return nil, err
	}
	return &Activity{
		Type: userActivity.Activity,
		Data: map[string]interface{}{
			"review": review,
		},
	}, nil
}

func GetDiscussionActivity(userActivity activity.UserActivity, userID uint) (*Activity, error) {
	discussionActivityData := activity.DiscussionActivityData{}
	if err := json.Unmarshal(userActivity.ActivityData, &discussionActivityData); err != nil {
		return nil, err
	}
	// TODO: Get discussion with game details
	discussion, err := discussions.GetDiscussionWithGameDetails(discussionActivityData.DiscussionID, userID)
	if err != nil {
		return nil, err
	}
	return &Activity{
		Type: userActivity.Activity,
		Data: map[string]interface{}{
			"discussion": discussion,
		},
	}, nil
}

func GetPostActivity(userActivity activity.UserActivity, userID uint) (*Activity, error) {
	postActivityData := activity.PostActivityData{}
	if err := json.Unmarshal(userActivity.ActivityData, &postActivityData); err != nil {
		return nil, err
	}
	discussion, err := discussions.GetDiscussionShortInfo(postActivityData.DiscussionID)
	if err != nil {
		return nil, err
	}
	post, err := posts.GetPostByID(postActivityData.PostID, userID)
	if err != nil {
		return nil, err
	}
	return &Activity{
		Type: userActivity.Activity,
		Data: map[string]interface{}{
			"discussion": discussion,
			"post":       post,
		},
	}, nil
}
