package lists

import (
	"errors"
	"github.com/kwisnia/inzynierka-backend/internal/api/games"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
)

type GameListDetails struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Games       []games.GameListElement `json:"games"`
}

func GetUserLists(userName string) ([]schema.GameList, error) {
	userCheck := user.GetByUsername(userName)
	if userCheck == nil {
		return nil, errors.New("user not found")
	}
	lists, err := GetByOwnerID(userCheck.ID)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func GetListDetails(listID int) (*GameListDetails, error) {
	list, err := GetByID(uint(listID))
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

func CreateList(owner uint, listDetails ListForm) error {
	list := schema.GameList{
		Name:        listDetails.Name,
		Description: listDetails.Description,
		Owner:       owner,
	}
	Save(&list)
	return nil
}

func AddGameToList(listId int, gameSlug string, requestingUserID uint) error {
	list, err := GetByID(uint(listId))
	if err != nil {
		return err
	}
	if list.Owner != requestingUserID {
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

func DeleteList(listId int, requestingUserID uint) error {
	list, err := GetByID(uint(listId))
	if err != nil {
		return err
	}
	if list.Owner != requestingUserID {
		return errors.New("you are not the owner of this list")
	}
	Delete(list)
	return nil
}

func UpdateList(listId int, listDetails ListForm, requestingUserID uint) error {
	list, err := GetByID(uint(listId))
	if err != nil {
		return err
	}
	if list.Owner != requestingUserID {
		return errors.New("you are not the owner of this list")
	}
	list.Name = listDetails.Name
	list.Description = listDetails.Description
	Update(list)
	return nil
}

func DeleteGameFromList(listId int, gameSlug string, requestingUserID uint) error {
	list, err := GetByID(uint(listId))
	if err != nil {
		return err
	}
	if list.Owner != requestingUserID {
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
