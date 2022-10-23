package reviews

import (
	"errors"
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/games"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"gorm.io/gorm"
)

type ReviewWithUserDetails struct {
	schema.GameReview
	User            user.BasicUserDetails `json:"user"`
	MarkedAsHelpful bool                  `json:"markedAsHelpful"`
}

func CreateReview(username string, gameSlug string, form ReviewForm) (*schema.GameReview, error) {
	_, err := games.GetGameBySlug(gameSlug)
	if err != nil {
		return nil, err
	}
	check := user.GetByUsername(username)
	if check == nil {
		return nil, err
	}
	review, err := GetByGameAndUser(gameSlug, username)
	if err != nil {
		review = &schema.GameReview{
			Game:             gameSlug,
			Creator:          username,
			Rating:           form.Rating,
			Title:            form.Title,
			Body:             form.Body,
			ContainsSpoilers: form.ContainsSpoilers,
			PlatformID:       form.PlatformID,
			GameCompletionID: form.GameCompletionID,
		}

	} else {
		review.Rating = form.Rating
		review.Title = form.Title
		review.Body = form.Body
		review.ContainsSpoilers = form.ContainsSpoilers
		review.PlatformID = form.PlatformID
		review.GameCompletionID = form.GameCompletionID
	}
	err = Save(review)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func GetReviewsForUser(pageSize int, page int, filters []int, username string) ([]schema.GameReview, error) {
	userReviews, err := GetByUserUsername(username, pageSize, page, filters)
	if err != nil {
		return nil, err
	}
	return userReviews, nil
}

func GetReviewsForGame(pageSize int, page int, filters []int, gameSlug string, userName string) ([]ReviewWithUserDetails, error) {
	offset := (page - 1) * pageSize
	gameReviews, err := GetByGameSlug(gameSlug, pageSize, offset, filters)
	if err != nil {
		return nil, err
	}
	reviewsWithUserDetails := make([]ReviewWithUserDetails, len(gameReviews))

	for i, review := range gameReviews {
		userDetails := user.GetBasicUserDetails(review.Creator)
		if userDetails == nil {
			return nil, err
		}
		_, err := GetHelpfulByUserAndReview(userName, review.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		reviewsWithUserDetails[i] = ReviewWithUserDetails{
			GameReview:      review,
			User:            *userDetails,
			MarkedAsHelpful: err == nil,
		}
	}
	return reviewsWithUserDetails, nil
}

func GetUserGameReview(gameSlug string, username string) (*schema.GameReview, error) {
	review, err := GetByGameAndUser(gameSlug, username)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func GetReview(reviewID uint, username string) (*ReviewWithUserDetails, error) {
	review, err := GetByID(reviewID)
	if err != nil {
		return nil, err
	}
	userDetails := user.GetBasicUserDetails(review.Creator)
	if userDetails == nil {
		return nil, err
	}
	_, err = GetHelpfulByUserAndReview(username, review.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &ReviewWithUserDetails{
		GameReview:      *review,
		User:            *userDetails,
		MarkedAsHelpful: err == nil,
	}, nil
}

func DeleteReview(username string, gameSlug string) error {
	review, err := GetByGameAndUser(gameSlug, username)
	if err != nil {
		return err
	}
	Delete(review)
	return nil
}

func MarkReviewAsHelpful(username string, reviewID uint) error {
	user := user.GetByUsername(username)
	if user == nil {
		return fmt.Errorf("user not found")
	}
	_, err := GetByID(reviewID)
	if err != nil {
		return fmt.Errorf("review not found")
	}
	_, err = GetHelpfulByUserAndReview(username, reviewID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		CreateHelpful(&schema.ReviewHelpful{
			ReviewID: reviewID,
			Username: username,
		})
	}
	return nil
}

func UnmarkReviewAsHelpful(username string, reviewID uint) error {
	user := user.GetByUsername(username)
	if user == nil {
		return fmt.Errorf("user not found")
	}
	_, err := GetByID(reviewID)
	if err != nil {
		return fmt.Errorf("review not found")
	}
	_, err = GetHelpfulByUserAndReview(username, reviewID)
	if err == nil {
		RemoveHelpful(&schema.ReviewHelpful{
			ReviewID: reviewID,
			Username: username,
		})
	}
	return nil
}
