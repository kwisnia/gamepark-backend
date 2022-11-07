package discussions

import (
	"errors"
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements/dispatcher"
	"github.com/kwisnia/inzynierka-backend/internal/api/games"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/discussions/posts"
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
	"regexp"
)

type DiscussionForm struct {
	Title string
	Body  string
}

type DiscussionListItemWithUserDetails struct {
	GameDiscussionListItem
	User       user.BasicUserDetails `json:"user"`
	UserScore  int                   `json:"userScore"`
	PostsCount int64                 `json:"postsCount"`
}

type DiscussionWithUserDetails struct {
	schema.GameDiscussion
	User       user.BasicUserDetails `json:"user"`
	UserScore  int                   `json:"userScore"`
	PostsCount int64                 `json:"postsCount"`
}

type DiscussionWithGameDetails struct {
	schema.GameDiscussion
	Game       games.GameListElement `json:"gameDetails"`
	UserScore  int                   `json:"userScore"`
	PostsCount int64                 `json:"postsCount"`
}

var policy = bluemonday.UGCPolicy()

func CreateDiscussion(userID uint, game string, discussionForm DiscussionForm) (*schema.GameDiscussion, error) {
	policy.AllowAttrs("data-youtube-video").OnElements("div")
	policy.AllowElements("iframe")
	policy.AllowAttrs("width").Matching(bluemonday.Number).OnElements("iframe")
	policy.AllowAttrs("height").Matching(bluemonday.Number).OnElements("iframe")
	policy.AllowAttrs("src").OnElements("iframe")
	policy.AllowAttrs("frameborder").Matching(bluemonday.Number).OnElements("iframe")
	policy.AllowAttrs("allow").Matching(regexp.MustCompile(`[a-z; -]*`)).OnElements("iframe")
	policy.AllowAttrs("allowfullscreen").OnElements("iframe")

	userCheck := user.GetBasicUserDetailsByID(userID)
	if userCheck == nil {
		return nil, fmt.Errorf("user not found")
	}
	discussion := schema.GameDiscussion{
		Title:     discussionForm.Title,
		Body:      policy.Sanitize(discussionForm.Body),
		Game:      game,
		CreatorID: userID,
	}
	if err := Save(&discussion); err != nil {
		return nil, err
	}
	go func() {
		userDiscussionCount, err := GetDiscussionCountForUser(userID)
		if err != nil {
			return
		}
		dispatcher.DispatchAchievementCheck(userID, achievements.ConditionTypeDiscussions, userDiscussionCount)
	}()
	return &discussion, nil
}

func GetDiscussionsForUser(pageSize int, page int, username string, userID uint) ([]DiscussionWithGameDetails, error) {
	discussionCreator := user.GetBasicUserDetailsByUsername(username)
	if discussionCreator == nil {
		return nil, fmt.Errorf("user not found")
	}
	offset := (page - 1) * pageSize
	discussions, err := GetByUserID(discussionCreator.ID, pageSize, offset)
	discussionsWithGameDetails := make([]DiscussionWithGameDetails, len(discussions))
	for i, discussion := range discussions {
		game, err := games.GetGameShortInfoBySlug(discussion.Game)
		if err != nil {
			return nil, err
		}
		score, err := GetScoreByUserAndDiscussion(userID, discussion.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		scoreValue := 0
		if score != nil {
			scoreValue = score.Score
		}
		postsCount, err := posts.GetPostsCountForDiscussion(discussion.ID)
		if err != nil {
			return nil, err
		}
		discussionsWithGameDetails[i] = DiscussionWithGameDetails{
			GameDiscussion: discussion,
			Game:           game,
			UserScore:      scoreValue,
			PostsCount:     postsCount,
		}
	}
	if err != nil {
		return nil, err
	}
	return discussionsWithGameDetails, nil
}

func GetDiscussionsForGame(pageSize int, page int, game string, userID uint) ([]DiscussionListItemWithUserDetails, error) {
	discussions, err := GetByGameSlug(game, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	discussionsWithUserDetails := make([]DiscussionListItemWithUserDetails, len(discussions))
	for i, discussion := range discussions {
		userDetails := user.GetBasicUserDetailsByID(discussion.CreatorID)
		if userDetails == nil {
			return nil, fmt.Errorf("user not found")
		}
		score, err := GetScoreByUserAndDiscussion(userID, discussion.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		scoreValue := 0
		if score != nil {
			scoreValue = score.Score
		}
		postsCount, err := posts.GetPostsCountForDiscussion(discussion.ID)
		if err != nil {
			return nil, err
		}
		discussionsWithUserDetails[i] = DiscussionListItemWithUserDetails{
			GameDiscussionListItem: discussion,
			User:                   *userDetails,
			UserScore:              scoreValue,
			PostsCount:             postsCount,
		}
	}

	return discussionsWithUserDetails, nil
}

func GetDiscussionByID(id uint, userID uint) (*DiscussionWithUserDetails, error) {
	discussion, err := GetByID(id)
	if err != nil {
		return nil, err
	}
	userDetails := user.GetBasicUserDetailsByID(discussion.CreatorID)
	if userDetails == nil {
		return nil, fmt.Errorf("user not found")
	}
	score, err := GetScoreByUserAndDiscussion(userID, discussion.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	scoreValue := 0
	if score != nil {
		scoreValue = score.Score
	}
	postsCount, err := posts.GetPostsCountForDiscussion(discussion.ID)
	if err != nil {
		return nil, err
	}
	discussionWithUserDetails := DiscussionWithUserDetails{
		GameDiscussion: *discussion,
		User:           *userDetails,
		UserScore:      scoreValue,
		PostsCount:     postsCount,
	}
	return &discussionWithUserDetails, nil
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
		err = DeleteScore(scoredDiscussion)
		if err != nil {
			return err
		}
		if scoredDiscussion.Score == score {
			return nil
		}
	}
	err = CreateDiscussionScore(&schema.DiscussionScore{
		DiscussionID: discussionID,
		UserID:       userID,
		Score:        score,
	})
	return nil
}

func DeleteScore(score *schema.DiscussionScore) error {
	return RemoveDiscussionScore(score)
}

func GetDiscussionCountForUser(userID uint) (int64, error) {
	count, err := CountByUser(userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
