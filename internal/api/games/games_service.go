package games

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"strings"
)

func getGames(pageSize int, page int, filters []int, sort string, search string) ([]GameListElement, error) {
	parsedSort := strings.Replace(sort, ".", " ", -1)
	offset := (page - 1) * pageSize
	games, err := GetPage(pageSize, offset, filters, parsedSort, search)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func getBySlug(slug string) (*schema.Game, error) {
	game, err := GetGameBySlug(slug)
	if err != nil {
		return nil, err
	}
	return &game, nil
}
