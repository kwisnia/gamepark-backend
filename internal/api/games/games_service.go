package games

import (
	"errors"
	"github.com/kwisnia/igdb/v2"
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	igdb2 "github.com/kwisnia/inzynierka-backend/internal/pkg/config/igdb"
	"github.com/kwisnia/inzynierka-backend/pkg/slice_funcs"
	"gorm.io/gorm"
	"strings"
	"time"
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

func CreateWebhookGame(gameId int) error {
	err := fetchGame(gameId, false)
	if err != nil {
		return err
	}
	return nil
}

func UpdateWebhookGame(gameId int) error {
	err := fetchGame(gameId, true)
	if err != nil {
		return err
	}
	return nil
}

func DeleteWebhookGame(gameId int) error {
	game, err := GetGameById(uint(gameId))
	if err != nil {
		return err
	}
	return DeleteGame(game.ID)
}

func fetchGame(gameId int, refetchGame bool) error {
	dbGame, err := prepareGame(gameId, false, refetchGame)
	if err != nil {
		return err
	}
	if dbGame.ParentGameID != nil {
		if err := fetchGame(int(*dbGame.ParentGameID), false); err != nil {
			if errors.Is(err, igdb.ErrNoResults) {
				dbGame.ParentGameID = nil
			} else {
				return err
			}
		}
	}
	if dbGame.VersionParentID != nil {
		if err := fetchGame(int(*dbGame.VersionParentID), false); err != nil {
			if errors.Is(err, igdb.ErrNoResults) {
				dbGame.VersionParentID = nil
			} else {
				return err
			}
		}
	}
	if refetchGame {
		err = UpdateGame(dbGame)
	} else {
		err = CreateGame(dbGame)
	}
	if err != nil {
		return err
	}
	return nil
}

func prepareGame(gameId int, getParentVersion bool, refetchGame bool) (*schema.Game, error) {
	dbGame := schema.Game{}
	_, err := GetGameById(uint(gameId))
	if (err != nil && errors.Is(err, gorm.ErrRecordNotFound)) || refetchGame {
		var parentGameId *uint = nil
		var versionParentId *uint = nil
		opts := igdb.ComposeOptions(
			igdb.SetFields("age_ratings.*", "artworks.*", "cover.*", "external_games.*", "involved_companies.*",
				"platforms.*", "platforms.platform_logo.*", "release_dates.*", "screenshots.*", "videos.*", "*"),
		)
		game, err := igdb2.IgdbClient.Games.Get(gameId, opts)
		if err != nil {
			if !errors.Is(err, igdb.ErrNoResults) {
				game, err = igdb2.IgdbClient.Games.Get(gameId, opts)
			} else {
				return nil, err
			}
		}
		if game.ParentGame != nil && game.ParentGame.ID != game.ID {
			parentGameId = new(uint)
			*parentGameId = uint(game.ParentGame.ID)
			if err != nil {
				return nil, err
			}
		}
		// get version parent game if exists
		if game.VersionParent != nil && game.VersionParent.ID != game.ID {
			versionParentId = new(uint)
			*versionParentId = uint(game.VersionParent.ID)
			if getParentVersion {
				if err := fetchGame(game.VersionParent.ID, false); err != nil {
					if errors.Is(err, igdb.ErrNoResults) {
						versionParentId = nil
					} else {
						return nil, err
					}
				}
			}
		}
		if err != nil {
			return nil, err
		}
		dbGame = *createBaseGameObject(game, parentGameId, versionParentId)
	}
	return &dbGame, nil
}

func createBaseGameObject(game *igdb.Game, parentGameId *uint, versionParentId *uint) *schema.Game {
	return &schema.Game{
		Model: gorm.Model{
			ID: uint(game.ID),
		},
		AgeRatings: slice_funcs.Map(game.AgeRatings, func(ageRating igdb.AgeRatingWrapper) schema.GameAgeRating {
			return schema.GameAgeRating{
				Model: gorm.Model{
					ID: uint(ageRating.ID),
				},
				AgeRatingID:    uint(ageRating.Rating),
				OrganizationID: uint(ageRating.Category),
				Synopsys:       &ageRating.Synopsis,
			}
		}),
		AggregatedRating:      game.AggregatedRating,
		AggregatedRatingCount: game.AggregatedRatingCount,
		Artworks: slice_funcs.Map(game.Artworks, func(artwork igdb.ArtworkWrapper) schema.Artwork {
			return schema.Artwork{
				Model: gorm.Model{
					ID: uint(artwork.ID),
				},
				Image: schema.Image{
					Height:  artwork.Height,
					Width:   artwork.Width,
					ImageID: artwork.ImageID,
					URL:     artwork.URL,
				},
			}
		}),
		CategoryID: uint(game.Category),
		Cover: schema.Cover{
			Model: gorm.Model{
				ID: uint(game.Cover.ID),
			},
			Image: schema.Image{
				Height:  game.Cover.Height,
				Width:   game.Cover.Width,
				ImageID: game.Cover.ImageID,
				URL:     game.Cover.URL,
			},
		},
		ExternalGames: slice_funcs.Map(filterExternalCategories(game.ExternalGames), func(externalGame igdb.ExternalGameWrapper) schema.ExternalGame {
			return schema.ExternalGame{
				Model: gorm.Model{
					ID: uint(externalGame.ID),
				},
				CategoryID: uint(externalGame.Category),
				UID:        externalGame.UID,
				URL:        externalGame.Url,
			}
		}),
		FirstReleaseDate: time.Unix(int64(game.FirstReleaseDate), 0),
		Genres: slice_funcs.Map(game.Genres, func(genre igdb.GenreWrapper) schema.Genre {
			return schema.Genre{
				EnumCategory: schema.EnumCategory{
					ID: uint(genre.ID),
				},
			}
		}),
		InvolvedCompanies: slice_funcs.Map(game.InvolvedCompanies, func(involvedCompany igdb.InvolvedCompanyWrapper) schema.InvolvedCompany {
			return schema.InvolvedCompany{
				Model: gorm.Model{
					ID: uint(involvedCompany.ID),
				},
				CompanyID:  uint(involvedCompany.Company.ID),
				Developer:  involvedCompany.Developer,
				Publisher:  involvedCompany.Publisher,
				Porting:    involvedCompany.Porting,
				Supporting: involvedCompany.Supporting,
			}
		}),
		Name:         game.Name,
		ParentGameID: parentGameId,
		Platforms: slice_funcs.Map(game.Platforms, func(platform igdb.PlatformWrapper) schema.Platform {
			return schema.Platform{
				ID:           uint(platform.ID),
				Name:         platform.Name,
				Abbreviation: platform.Abbreviation,
				Generation:   platform.Generation,
				Logo: schema.PlatformLogo{
					Model: gorm.Model{
						ID: uint(platform.PlatformLogo.ID),
					},
					Image: schema.Image{
						Height:  platform.PlatformLogo.Height,
						Width:   platform.PlatformLogo.Width,
						ImageID: platform.PlatformLogo.ImageID,
						URL:     platform.PlatformLogo.URL,
					},
				},
				Slug:    platform.Slug,
				IGDBUrl: platform.URL,
			}
		}),
		Rating:      0,
		RatingCount: 0,
		ReleaseDates: slice_funcs.Map(game.ReleaseDates, func(releaseDate igdb.ReleaseDateWrapper) schema.ReleaseDate {
			return schema.ReleaseDate{
				Model: gorm.Model{
					ID: uint(releaseDate.ID),
				},
				RegionID:   uint(releaseDate.Region),
				Human:      releaseDate.Human,
				Date:       time.Unix(int64(releaseDate.Date), 0),
				PlatformID: uint(releaseDate.Platform.ID),
				CategoryID: uint(releaseDate.Category),
			}
		}),
		Screenshots: slice_funcs.Map(game.Screenshots, func(screenshot igdb.ScreenshotWrapper) schema.Screenshot {
			return schema.Screenshot{
				Model: gorm.Model{
					ID: uint(screenshot.ID),
				},
				Image: schema.Image{
					Height:  screenshot.Height,
					Width:   screenshot.Width,
					ImageID: screenshot.ImageID,
					URL:     screenshot.URL,
				},
			}
		}),
		Slug:            game.Slug,
		Storyline:       game.Storyline,
		Summary:         game.Summary,
		IGDBUrl:         game.URL,
		VersionParentID: versionParentId,
		VersionTitle:    game.VersionTitle,
		Videos: slice_funcs.Map(game.Videos, func(video igdb.GameVideoWrapper) schema.GameVideo {
			return schema.GameVideo{
				Model: gorm.Model{
					ID: uint(video.ID),
				},
				Video: schema.Video{
					VideoID: video.VideoID,
					Name:    video.Name,
				},
			}
		}),
	}
}

func filterExternalCategories(externalGames []igdb.ExternalGameWrapper) []igdb.ExternalGameWrapper {
	externalCategories, err := GetExternalCategories()
	if err != nil {
		return externalGames
	}
	return slice_funcs.Filter(externalGames, func(externalGame igdb.ExternalGameWrapper) bool {
		for _, category := range externalCategories {
			if uint(externalGame.Category) == category.ID {
				return true
			}
		}
		return false
	})
}
