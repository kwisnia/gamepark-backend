package user_game_info

import (
	"errors"
	"github.com/kwisnia/inzynierka-backend/internal/api/games"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/lists"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/reviews"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"gorm.io/gorm"
)

type UserGameDetails struct {
	Lists  []schema.GameList              `json:"lists"`
	Review *reviews.ReviewWithUserDetails `json:"review"`
}

func GetUserGameInfo(slug string, userID uint) (*UserGameDetails, error) {
	game, err := games.GetGameBySlug(slug)
	if err != nil {
		return nil, err
	}
	listsWhereGameIs, err := lists.GetUsersListsWhereGameIs(game.ID, userID)
	if err != nil {
		return nil, err
	}
	userGameReview, err := reviews.GetUserGameReview(game.Slug, userID)
	if err != nil {
		return nil, err
	}
	userDetails := user.GetBasicUserDetailsByID(userID)
	_, err = reviews.GetHelpfulByUserAndReview(userID, userGameReview.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &UserGameDetails{
		Lists: listsWhereGameIs,
		Review: &reviews.ReviewWithUserDetails{
			GameReview:      *userGameReview,
			User:            *userDetails,
			MarkedAsHelpful: err == nil,
		},
	}, nil
}
