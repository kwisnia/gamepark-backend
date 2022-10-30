package lists

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ListForm struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type ChangeListContentForm struct {
	Slug string `json:"slug" binding:"required"`
}

func GetUserListsHandler(c *gin.Context) {
	userName := c.Param("userName")
	lists, err := GetUserLists(userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, lists)
}

func GetUserListHandler(c *gin.Context) {
	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid list id"})
		return
	}
	list, err := GetListDetails(listID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func AddGameToListHandler(c *gin.Context) {
	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid list id"})
		return
	}
	var game ChangeListContentForm
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid game id"})
		return
	}
	userID := c.GetUint("userID")
	err = AddGameToList(listID, game.Slug, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game added to list"})
}

func RemoveGameFromListHandler(c *gin.Context) {
	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid list id"})
		return
	}
	var game ChangeListContentForm
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userID := c.GetUint("userID")
	err = DeleteGameFromList(listID, game.Slug, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game removed from list"})
}

func CreateListHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	var list ListForm
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := CreateList(userID, list)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func DeleteListHandler(c *gin.Context) {
	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid list id"})
		return
	}
	userID := c.GetUint("userID")
	err = DeleteList(listID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "List deleted"})
}

func UpdateListHandler(c *gin.Context) {
	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid list id"})
		return
	}
	userID := c.GetUint("userID")
	var list ListForm
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = UpdateList(listID, list, userID)
	if err != nil {
		if err.Error() == "you are not the owner of this list" {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, list)
}
