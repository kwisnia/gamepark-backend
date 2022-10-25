package discussions

import (
	"errors"
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"gorm.io/gorm"
)

type DiscussionForm struct {
	title string
}

type DiscussionWithUserDetails struct {
	schema.GameDiscussion
	User  user.BasicUserDetails `json:"user"`
	Score int                   `json:"score"`
}

func CreateDiscussion(userID uint, game string, discussionForm DiscussionForm) (*schema.GameDiscussion, error) {
	userCheck := user.GetBasicUserDetailsByID(userID)
	if userCheck == nil {
		return nil, fmt.Errorf("user not found")
	}
	discussion := schema.GameDiscussion{
		Title:     discussionForm.title,
		Game:      game,
		CreatorID: userID,
	}
	if err := Save(&discussion); err != nil {
		return nil, err
	}
	return &discussion, nil
}

func GetDiscussionsForUser(pageSize int, page int, userID uint) ([]schema.GameDiscussion, error) {
	discussionCreator := user.GetBasicUserDetailsByID(userID)
	if discussionCreator == nil {
		return nil, fmt.Errorf("user not found")
	}
	discussions, err := GetByUserID(userID, pageSize, page)
	if err != nil {
		return nil, err
	}
	return discussions, nil
}

func GetDiscussionsForGame(pageSize int, page int, game string) ([]schema.GameDiscussion, error) {
	discussions, err := GetByGameSlug(game, pageSize, page)
	if err != nil {
		return nil, err
	}
	discussionsWithUserDetails := make([]DiscussionWithUserDetails, len(discussions))
	for i, discussion := range discussions {
		userDetails := user.GetBasicUserDetailsByID(discussion.CreatorID)
		if userDetails == nil {
			return nil, fmt.Errorf("user not found")
		}
		score, err := GetScoreByUserAndDiscussion(userDetails.ID, discussion.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		discussionsWithUserDetails[i] = DiscussionWithUserDetails{
			GameDiscussion: discussion,
			User:           *userDetails,
			Score:          score.Score,
		}
	}
	return discussions, nil
}

func GetDiscussionByID(id uint) (*schema.GameDiscussion, error) {
	discussion, err := GetByID(id)
	if err != nil {
		return nil, err
	}
	return discussion, nil
}

func DeleteDiscussion(id uint, userID uint) error {
	discussion, err := GetByID(id)
	if err != nil {
		return err
	}
	if discussion.CreatorID != userID {
		return fmt.Errorf("user is not creator of discussion")
	}
	err = Delete(discussion)
	if err != nil {
		return err
	}
	return nil
}

func ScoreDiscussion(userID uint, discussionID uint, score int) error {
	_, err := GetByID(discussionID)
	if err != nil {
		return err
	}
	userDetails := user.GetBasicUserDetailsByID(userID)
	if userDetails == nil {
		return fmt.Errorf("user not found")
	}
	if score > 0 {
		score = 1
	} else if score < 0 {
		score = -1
	} else {
		return nil
	}
	scoredDiscussion, err := GetScoreByUserAndDiscussion(userID, discussionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if scoredDiscussion != nil {
		err = DeleteScore(userID, discussionID)
		if err != nil {
			return err
		}
	}
	err = CreateDiscussionScore(&schema.DiscussionScore{
		DiscussionID: discussionID,
		UserID:       userID,
		Score:        score,
	})
	return nil
}

func DeleteScore(userID uint, discussionID uint) error {
	return RemoveDiscussionScore(&schema.DiscussionScore{
		DiscussionID: discussionID,
		UserID:       userID,
	})
}
