package lists

import (
	"fmt"
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

func GetUserLists(c *gin.Context) {
	userName := c.Param("userName")
	lists, err := getUserLists(userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, lists)
}

func GetUserList(c *gin.Context) {
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid list id"})
		return
	}
	list, err := getListDetails(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func AddGameToList(c *gin.Context) {
	listId, err := strconv.Atoi(c.Param("id"))
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
	userName := c.GetString("userName")
	err = addGameToList(listId, game.Slug, userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game added to list"})
}

func RemoveGameFromList(c *gin.Context) {
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid list id"})
		return
	}
	var game ChangeListContentForm
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userName := c.GetString("userName")
	err = deleteGameFromList(listId, game.Slug, userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Game removed from list"})
}

func CreateList(c *gin.Context) {
	userName := c.GetString("userName")
	var list ListForm
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := createList(userName, list)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func DeleteList(c *gin.Context) {
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid list id"})
		return
	}
	userName := c.GetString("userName")
	err = deleteList(listId, userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "List deleted"})
}

func UpdateList(c *gin.Context) {
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid list id"})
		return
	}
	userName := c.GetString("userName")
	var list ListForm
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = updateList(listId, list, userName)
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

func GetUserGameInfo(c *gin.Context) {
	slug := c.Param("slug")
	userName := c.GetString("userName")
	gameUserInfo, _ := getUserGameInfo(slug, userName)
	fmt.Println(gameUserInfo)
	c.JSON(http.StatusOK, gameUserInfo)
}
