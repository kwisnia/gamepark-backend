package games

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"strings"
)

type GameWithSimilarGames struct {
	schema.Game
	SimilarGames []GameListElement `json:"similarGames"`
}

func GetGames(pageSize int, page int, filters []int, sort string, search string) ([]GameListElement, error) {
	parsedSort := strings.Replace(sort, ".", " ", -1)
	offset := (page - 1) * pageSize
	games, err := GetPage(pageSize, offset, filters, parsedSort, search)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func GetBySlug(slug string) (*GameWithSimilarGames, error) {
	game, err := GetGameBySlug(slug)
	if err != nil {
		return nil, err
	}
	similarGames, err := GetSimilarGames(game.ID)
	if err != nil {
		return nil, err
	}
	similarGamesIds := make([]uint, len(similarGames))
	for i, similarGame := range similarGames {
		if similarGame.GameID == game.ID {
			similarGamesIds[i] = similarGame.SimilarGameID
		} else {
			similarGamesIds[i] = similarGame.GameID
		}
	}
	similarGamesShortInfo, err := GetGameShortInfosByIds(similarGamesIds)
	if err != nil {
		return nil, err
	}
	gameWithSimilarGames := GameWithSimilarGames{
		Game:         game,
		SimilarGames: similarGamesShortInfo,
	}

	return &gameWithSimilarGames, nil
}

func GetShortInfoBySlug(slug string) (*GameListElement, error) {
	game, err := GetGameShortInfoBySlug(slug)
	if err != nil {
		return nil, err
	}
	return &game, nil
}
