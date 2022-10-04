package lists

import (
	"errors"
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/games"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
)

type GameListDetails struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Games       []games.GameListElement `json:"games"`
}

type UserGameDetails struct {
	Lists []schema.GameList `json:"lists"`
}

func getUserLists(userName string) ([]schema.GameList, error) {
	fmt.Println(userName)
	lists, err := GetByOwnerUsername(userName)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func getListDetails(listId int) (*GameListDetails, error) {
	list, err := GetByID(uint(listId))
	if err != nil {
		return nil, err
	}
	listGames, err := GetGames(list)
	if err != nil {
		return nil, err
	}
	return &GameListDetails{
		Name:        list.Name,
		Description: list.Description,
		Games:       listGames,
	}, nil
}

func createList(owner string, listDetails ListForm) error {
	list := schema.GameList{
		Name:        listDetails.Name,
		Description: listDetails.Description,
		Owner:       owner,
	}
	Save(&list)
	return nil
}

func addGameToList(listId int, gameSlug string, requestingUser string) error {
	list, err := GetByID(uint(listId))
	if err != nil {
		return err
	}
	if list.Owner != requestingUser {
		return errors.New("you are not the owner of this list")
	}
	game, err := games.GetGameBySlug(gameSlug)
	if err != nil {
		return err
	}
	err = AddGame(list, &game)
	if err != nil {
		return err
	}
	return nil
}

func deleteList(listId int, requestingUser string) error {
	list, err := GetByID(uint(listId))
	if err != nil {
		return err
	}
	if list.Owner != requestingUser {
		return errors.New("you are not the owner of this list")
	}
	Delete(list)
	return nil
}

func updateList(listId int, listDetails ListForm, requestingUser string) error {
	list, err := GetByID(uint(listId))
	if err != nil {
		return err
	}
	if list.Owner != requestingUser {
		return errors.New("you are not the owner of this list")
	}
	list.Name = listDetails.Name
	list.Description = listDetails.Description
	Update(list)
	return nil
}

func deleteGameFromList(listId int, gameSlug string, requestingUser string) error {
	list, err := GetByID(uint(listId))
	if err != nil {
		return err
	}
	if list.Owner != requestingUser {
		return errors.New("you are not the owner of this list")
	}
	game, err := games.GetGameBySlug(gameSlug)
	if err != nil {
		return err
	}
	err = RemoveGame(list, &game)
	if err != nil {
		return err
	}
	return nil
}

func getUserGameInfo(slug string, userName string) (*UserGameDetails, error) {
	game, err := games.GetGameBySlug(slug)
	if err != nil {
		return nil, err
	}
	listsWhereGameIs, err := GetUsersListsWhereGameIs(game.ID, userName)
	fmt.Println(listsWhereGameIs)
	if err != nil {
		return nil, err
	}
	return &UserGameDetails{
		Lists: listsWhereGameIs,
	}, nil
}
