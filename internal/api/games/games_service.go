package games

func getGames(pageSize int, afterId int) ([]GameListElement, error) {
	games, err := GetPage(pageSize, afterId)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func getBySlug(slug string) (*Game, error) {
	game, err := GetGameBySlug(slug)
	if err != nil {
		return nil, err
	}
	return &game, nil
}
