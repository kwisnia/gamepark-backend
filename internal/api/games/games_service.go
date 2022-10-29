package games

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/games/schema"
	"strings"
)

func GetGames(pageSize int, page int, filters []int, sort string, search string) ([]GameListElement, error) {
	parsedSort := strings.Replace(sort, ".", " ", -1)
	offset := (page - 1) * pageSize
	games, err := GetPage(pageSize, offset, filters, parsedSort, search)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func GetBySlug(slug string) (*schema.Game, error) {
	game, err := GetGameBySlug(slug)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func GetShortInfoBySlug(slug string) (*GameListElement, error) {
	game, err := GetGameShortInfoBySlug(slug)
	if err != nil {
		return nil, err
	}
	return &game, nil
}
