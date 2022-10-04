package lists

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/games"
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
)

func Save(l *schema.GameList) {
	database.DB.Create(l)
}

func GetByID(id uint) (*schema.GameList, error) {
	var l schema.GameList
	if err := database.DB.Where("id = ?", id).First(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func GetByOwnerUsername(userName string) ([]schema.GameList, error) {
	var l []schema.GameList
	if err := database.DB.Where("owner = ?", userName).Find(&l).Error; err != nil {
		return nil, err
	}
	return l, nil
}

func GetByGameID(id uint) ([]schema.GameList, error) {
	var l []schema.GameList
	if err := database.DB.Where("games.id = ?", id).Find(&l).Error; err != nil {
		return nil, err
	}
	return l, nil
}

func Delete(l *schema.GameList) {
	database.DB.Delete(l)
}

func Update(l *schema.GameList) {
	database.DB.Save(l)
}

func AddGame(l *schema.GameList, g *schema.Game) error {
	err := database.DB.Model(l).Association("Games").Append(g)
	if err != nil {
		return err
	}
	return nil
}

func RemoveGame(l *schema.GameList, g *schema.Game) error {
	err := database.DB.Model(l).Association("Games").Delete(g)
	if err != nil {
		return err
	}
	return nil
}

func GetGames(l *schema.GameList) ([]games.GameListElement, error) {
	var g []games.GameListElement
	if err := database.DB.Model(l).Preload("Cover").Association("Games").Find(&g); err != nil {
		return nil, err
	}
	return g, nil
}

func GetUsersListsWhereGameIs(gameId uint, userName string) ([]schema.GameList, error) {
	var l []schema.GameList
	if err := database.DB.Model(&schema.GameList{}).Joins(
		"LEFT JOIN list_games ON list_games.game_list_id = game_lists.id",
	).Where("game_lists.owner = ? AND list_games.game_id = ?", userName, gameId).Scan(&l).Error; err != nil {
		return nil, err
	}
	return l, nil
}
