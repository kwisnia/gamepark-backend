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

func CreateReview(userID uint, gameSlug string, form ReviewForm) (*schema.GameReview, error) {
	_, err := games.GetGameBySlug(gameSlug)
	if err != nil {
		return nil, err
	}
	check := user.GetByID(userID)
	if check == nil {
		return nil, err
	}
	review, err := GetByGameAndUser(gameSlug, userID)
	if err != nil {
		review = &schema.GameReview{
			Game:             gameSlug,
			Creator:          userID,
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
	userCheck := user.GetByUsername(username)
	if userCheck == nil {
		return nil, fmt.Errorf("user not found")
	}
	userReviews, err := GetByUserID(userCheck.ID, pageSize, page, filters)
	if err != nil {
		return nil, err
	}
	return userReviews, nil
}

func GetReviewsForGame(pageSize int, page int, filters []int, gameSlug string, userID uint) ([]ReviewWithUserDetails, error) {
	offset := (page - 1) * pageSize
	gameReviews, err := GetByGameSlug(gameSlug, pageSize, offset, filters)
	if err != nil {
		return nil, err
	}
	reviewsWithUserDetails := make([]ReviewWithUserDetails, len(gameReviews))

	for i, review := range gameReviews {
		userDetails := user.GetBasicUserDetailsByID(review.Creator)
		if userDetails == nil {
			return nil, err
		}
		_, err := GetHelpfulByUserAndReview(userID, review.ID)
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

func GetUserGameReview(gameSlug string, userID uint) (*schema.GameReview, error) {
	review, err := GetByGameAndUser(gameSlug, userID)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func GetReview(reviewID uint, userID uint) (*ReviewWithUserDetails, error) {
	review, err := GetByID(reviewID)
	if err != nil {
		return nil, err
	}
	userDetails := user.GetBasicUserDetailsByID(review.Creator)
	if userDetails == nil {
		return nil, err
	}
	_, err = GetHelpfulByUserAndReview(userID, review.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &ReviewWithUserDetails{
		GameReview:      *review,
		User:            *userDetails,
		MarkedAsHelpful: err == nil,
	}, nil
}

func DeleteReview(userID uint, gameSlug string) error {
	review, err := GetByGameAndUser(gameSlug, userID)
	if err != nil {
		return err
	}
	Delete(review)
	return nil
}

func MarkReviewAsHelpful(userID uint, reviewID uint) error {
	userCheck := user.GetByID(userID)
	if userCheck == nil {
		return fmt.Errorf("user not found")
	}
	_, err := GetByID(reviewID)
	if err != nil {
		return fmt.Errorf("review not found")
	}
	_, err = GetHelpfulByUserAndReview(userID, reviewID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		CreateHelpful(&schema.ReviewHelpful{
			ReviewID: reviewID,
			UserID:   userID,
		})
	}
	return nil
}

func UnmarkReviewAsHelpful(userID uint, reviewID uint) error {
	userCheck := user.GetByID(userID)
	if userCheck == nil {
		return fmt.Errorf("user not found")
	}
	_, err := GetByID(reviewID)
	if err != nil {
		return fmt.Errorf("review not found")
	}
	_, err = GetHelpfulByUserAndReview(userID, reviewID)
	if err == nil {
		RemoveHelpful(&schema.ReviewHelpful{
			ReviewID: reviewID,
			UserID:   userID,
		})
	}
	return nil
}

func GetReviewCountForUser(userID uint) (int64, error) {
	count, err := CountByUser(userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
