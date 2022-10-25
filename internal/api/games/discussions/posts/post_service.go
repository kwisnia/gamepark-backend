package posts

import (
	"errors"
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

type PostForm struct {
	content        string
	originalPostID *uint
}

type PostWithUserDetails struct {
	schema.DiscussionPost
	User  user.BasicUserDetails `json:"user"`
	Score int                   `json:"score"`
}

var policy = bluemonday.UGCPolicy()

func CreatePost(userID uint, discussionID uint, postForm PostForm) (*schema.DiscussionPost, error) {
	userCheck := user.GetBasicUserDetailsByID(userID)
	if userCheck == nil {
		return nil, fmt.Errorf("user not found")
	}
	post := schema.DiscussionPost{
		Body:           policy.Sanitize(postForm.content),
		DiscussionID:   discussionID,
		CreatorID:      userID,
		OriginalPostID: postForm.originalPostID,
	}
	if err := Save(&post); err != nil {
		return nil, err
	}
	return &post, nil
}

func GetPostsForDiscussion(pageSize int, page int, discussionID uint) ([]PostWithUserDetails, error) {
	posts, err := GetByDiscussionID(discussionID, pageSize, page)
	if err != nil {
		return nil, err
	}
	postsWithUserDetails := make([]PostWithUserDetails, len(posts))
	for i, post := range posts {
		userDetails := user.GetBasicUserDetailsByID(post.CreatorID)
		if userDetails == nil {
			return nil, fmt.Errorf("user not found")
		}
		score, err := GetScoreByUserAndPost(userDetails.ID, post.ID)
		if err != nil {
			return nil, err
		}
		postsWithUserDetails[i] = PostWithUserDetails{
			DiscussionPost: post,
			User:           *userDetails,
			Score:          score.Score,
		}
	}
	return postsWithUserDetails, nil
}

func GetPostByID(id uint) (*schema.DiscussionPost, error) {
	post, err := GetByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func DeletePost(id uint, userID uint) error {
	post, err := GetByID(id)
	if err != nil {
		return err
	}
	if post.CreatorID != userID {
		return fmt.Errorf("user is not the creator of this post")
	}
	return Delete(post)
}

func UpdatePost(id uint, userID uint, postForm PostForm) (*schema.DiscussionPost, error) {
	post, err := GetByID(id)
	if err != nil {
		return nil, err
	}
	if post.CreatorID != userID {
		return nil, fmt.Errorf("user is not the creator of this post")
	}
	post.Body = policy.Sanitize(postForm.content)
	if err := Save(post); err != nil {
		return nil, err
	}
	return post, nil
}

func ScorePost(userID uint, postID uint, score int) error {
	_, err := GetByID(postID)
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
	scoredPost, err := GetScoreByUserAndPost(userID, postID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if scoredPost != nil {
		err = DeleteScore(scoredPost)
		if err != nil {
			return err
		}
	}
	err = CreatePostScore(&schema.PostScore{
		PostID: postID,
		Score:  score,
		UserID: userID,
	})
	return err
}
