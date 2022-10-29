package discussions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/discussions/posts"
	"net/http"
	"strconv"
)

type CreateDiscussionForm struct {
	Title string
	Body  string
}

type ScoreForm struct {
	Score int
}

// TODO: Handle errors in a better way than comparing strings

func GetDiscussionsForGameHandler(c *gin.Context) {
	gameSlug := c.Param("slug")
	userID := c.GetUint("userID")
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page size"})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid after"})
		return
	}
	discussions, err := GetDiscussionsForGame(pageSize, page, gameSlug, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, discussions)
}

func CreateDiscussionHandler(c *gin.Context) {
	gameSlug := c.Param("slug")
	userID := c.GetUint("userID")
	var discussionForm CreateDiscussionForm
	if err := c.ShouldBindJSON(&discussionForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(discussionForm.Title)
	fmt.Println(discussionForm.Body)
	discussion, err := CreateDiscussion(userID, gameSlug, DiscussionForm{
		Title: discussionForm.Title,
		Body:  discussionForm.Body,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, discussion)
}

func GetDiscussionsForUserHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page size"})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid after"})
		return
	}
	discussions, err := GetDiscussionsForUser(pageSize, page, userID)
	if err != nil {
		var errorStatus int
		if err.Error() == "user not found" {
			errorStatus = http.StatusNotFound
		} else {
			errorStatus = http.StatusInternalServerError
		}
		c.JSON(errorStatus, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, discussions)
}

func GetDiscussionHandler(c *gin.Context) {
	discussionID, err := strconv.Atoi(c.Param("discussionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion id"})
		return
	}
	userID := c.GetUint("userID")
	discussion, err := GetDiscussionByID(uint(discussionID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, discussion)
}

func DeleteDiscussionHandler(c *gin.Context) {
	discussionID, err := strconv.Atoi(c.Param("discussionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion id"})
		return
	}
	userID := c.GetUint("userID")
	err = DeleteDiscussion(uint(discussionID), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Discussion deleted"})
}

func ScoreDiscussionHandler(c *gin.Context) {
	discussionID, err := strconv.Atoi(c.Param("discussionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion id"})
		return
	}
	userID := c.GetUint("userID")
	var scoreForm ScoreForm
	if err := c.ShouldBindJSON(&scoreForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = ScoreDiscussion(userID, uint(discussionID), scoreForm.Score)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Discussion scored"})
}

func GetDiscussionPostsHandler(c *gin.Context) {
	discussionID, err := strconv.Atoi(c.Param("discussionId"))
	userID := c.GetUint("userID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion id"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page size"})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid after"})
		return
	}
	discussionPosts, err := posts.GetPostsForDiscussion(pageSize, page, uint(discussionID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, discussionPosts)
}

func CreateDiscussionPostHandler(c *gin.Context) {
	discussionID, err := strconv.Atoi(c.Param("discussionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid discussion id"})
		return
	}
	userID := c.GetUint("userID")
	var postForm posts.PostForm
	if err := c.ShouldBindJSON(&postForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(postForm.Body)
	post, err := posts.CreatePost(userID, uint(discussionID), posts.PostForm{
		OriginalPostID: postForm.OriginalPostID,
		Body:           postForm.Body,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, post)
}

func UpdateDiscussionPostHandler(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post id"})
		return
	}
	userID := c.GetUint("userID")
	var postForm posts.PostForm
	if err := c.ShouldBindJSON(&postForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	post, err := posts.UpdatePost(uint(postID), userID, posts.PostForm{
		Body: postForm.Body,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, post)
}

func DeleteDiscussionPostHandler(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post id"})
		return
	}
	userID := c.GetUint("userID")
	err = posts.DeletePost(uint(postID), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

func ScoreDiscussionPostHandler(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post id"})
		return
	}
	userID := c.GetUint("userID")
	var scoreForm ScoreForm
	if err := c.ShouldBindJSON(&scoreForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = posts.ScorePost(userID, uint(postID), scoreForm.Score)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post scored"})
}

func GetDiscussionPostRepliesHandler(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	userID := c.GetUint("userID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post id"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page size"})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid after"})
		return
	}
	postReplies, err := posts.GetPostReplies(pageSize, page, uint(postID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, postReplies)
}
