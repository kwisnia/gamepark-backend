package posts

import (
	"errors"
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements/dispatcher"
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"gorm.io/gorm"
)

type PostForm struct {
	Body           string
	OriginalPostID *uint
}

type PostWithUserDetails struct {
	schema.DiscussionPost
	User       user.BasicUserDetails `json:"user"`
	UserScore  int                   `json:"userScore"`
	ReplyCount int64                 `json:"replyCount"`
}

func CreatePost(userID uint, discussionID uint, postForm PostForm) (*schema.DiscussionPost, error) {
	userCheck := user.GetBasicUserDetailsByID(userID)
	if userCheck == nil {
		return nil, fmt.Errorf("user not found")
	}
	post := schema.DiscussionPost{
		Body:           postForm.Body,
		DiscussionID:   discussionID,
		CreatorID:      userID,
		OriginalPostID: postForm.OriginalPostID,
	}
	if err := SaveNewPost(&post); err != nil {
		return nil, err
	}
	go func() {
		userPostsCount, err := GetPostCountForUser(userID)
		if err != nil {
			return
		}
		dispatcher.DispatchAchievementCheck(userID, achievements.ConditionTypePosts, userPostsCount)
	}()
	return &post, nil
}

func GetPostsForDiscussion(pageSize int, page int, discussionID uint, userID uint) ([]PostWithUserDetails, error) {
	posts, err := GetWithoutRepliesByDiscussionID(discussionID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	postsWithUserDetails := make([]PostWithUserDetails, len(posts))
	for i, post := range posts {
		userDetails := user.GetBasicUserDetailsByID(post.CreatorID)
		if userDetails == nil {
			return nil, fmt.Errorf("user not found")
		}
		score, err := GetScoreByUserAndPost(userID, post.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		scoreValue := 0
		if score != nil {
			scoreValue = score.Score
		}
		repliesCount, err := GetReplyCountForPost(post.ID)
		if err != nil {
			return nil, err
		}
		postsWithUserDetails[i] = PostWithUserDetails{
			DiscussionPost: post,
			User:           *userDetails,
			UserScore:      scoreValue,
			ReplyCount:     repliesCount,
		}
	}
	return postsWithUserDetails, nil
}

func GetPostReplies(pageSize int, page int, postID uint, userID uint) ([]PostWithUserDetails, error) {
	posts, err := GetRepliesForPost(postID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	postsWithUserDetails := make([]PostWithUserDetails, len(posts))
	for i, post := range posts {
		userDetails := user.GetBasicUserDetailsByID(post.CreatorID)
		if userDetails == nil {
			return nil, fmt.Errorf("user not found")
		}
		score, err := GetScoreByUserAndPost(userID, post.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		scoreValue := 0
		if score != nil {
			scoreValue = score.Score
		}
		repliesCount, err := GetReplyCountForPost(post.ID)
		if err != nil {
			return nil, err
		}
		postsWithUserDetails[i] = PostWithUserDetails{
			DiscussionPost: post,
			User:           *userDetails,
			UserScore:      scoreValue,
			ReplyCount:     repliesCount,
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
	post.Body = postForm.Body
	if err := Update(post); err != nil {
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
		if scoredPost.Score == score {
			return nil
		}
	}
	err = CreatePostScore(&schema.PostScore{
		PostID: postID,
		Score:  score,
		UserID: userID,
	})
	return err
}

func GetPostCountForUser(userID uint) (int64, error) {
	count, err := CountByUser(userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
