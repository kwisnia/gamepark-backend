package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kwisnia/igdb/v2"
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"github.com/kwisnia/inzynierka-backend/pkg/slice_funcs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"time"
)

var logger = log.Default()
var db *gorm.DB
var client = igdb.NewClient(os.Getenv("IGDB_CLIENT_ID"), os.Getenv("IGDB_ACCESS_TOKEN"), nil)
var externalCategories []schema.ExternalCategory

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client = igdb.NewClient(os.Getenv("IGDB_CLIENT_ID"), os.Getenv("IGDB_ACCESS_TOKEN"), nil)
	database.SetupDB()
	db = database.DB
	db.Find(&externalCategories)
	_, err = client.Games.Get(1, igdb.ComposeOptions(igdb.SetFields("name")))
	if err != nil {
		log.Fatal(err)
	}

	clearDatabase()
	//err = fetchCompanies()
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//err = fetchGenres()
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//err = fetchGames(112500)
	//if err != nil {
	//	logger.Fatal(err)
	//}
	fmt.Println("huh")
	err = associateSimilarGames()
	if err != nil {
		return
	}
}

func fetchCompanies() error {
	companyCount, err := client.Companies.Count()
	if err != nil {
		return err
	}
	for i := 0; i < companyCount; i += 500 {
		logger.Println("Fetching company page ", i)
		opts := igdb.ComposeOptions(
			igdb.SetLimit(500),
			igdb.SetFields("name", "description", "start_date", "slug", "logo.*", "changed_company_id.*"),
			igdb.SetOrder("id", igdb.OrderAscending),
			igdb.SetOffset(i),
		)
		companies, err := client.Companies.Index(opts)
		if err != nil {
			return err
		}
		for _, company := range companies {
			result := db.First(&schema.Company{}, company.ID)
			if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
				var changedCompanyId *uint = nil
				if company.ChangedCompanyID != nil && company.ChangedCompanyID.ID != company.ID {
					changedCompanyId = new(uint)
					*changedCompanyId = uint(company.ChangedCompanyID.ID)
					err := fetchCompany(company.ChangedCompanyID.ID)
					if err != nil {
						return err
					}
				}
				dbCompany := schema.Company{
					Model: gorm.Model{
						ID: uint(company.ID),
					},
					ChangedCompanyID: changedCompanyId,
					CompanyLogo: schema.CompanyLogo{
						Model: gorm.Model{
							ID: uint(company.Logo.ID),
						},
						Image: schema.Image{
							Height:  company.Logo.Height,
							Width:   company.Logo.Width,
							ImageID: company.Logo.ImageID,
							URL:     company.Logo.URL,
						},
					},
					Name:        company.Name,
					Description: &company.Description,
					Slug:        company.Slug,
					StartDate:   time.Unix(int64(company.StartDate), 0),
					IGDBUrl:     company.URL,
				}
				if create := db.Create(&dbCompany); create.Error != nil {
					return create.Error
				}
			}
		}
	}
	return nil
}

func fetchCompany(companyId int) error {
	result := db.First(&schema.Company{}, companyId)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		opts := igdb.ComposeOptions(
			igdb.SetFields("name", "description", "start_date", "slug", "logo.*", "changed_company_id.*"),
		)
		company, err := client.Companies.Get(companyId, opts)
		if err != nil {
			return err
		}
		if company.ChangedCompanyID != nil && company.ChangedCompanyID.ID != companyId {
			err := fetchCompany(company.ChangedCompanyID.ID)
			if err != nil {
				return err
			}
		}
		var changedCompanyId *uint = nil
		if company.ChangedCompanyID != nil {
			changedCompanyId := new(uint)
			*changedCompanyId = uint(company.ChangedCompanyID.ID)
		}
		return db.Create(&schema.Company{
			Model: gorm.Model{
				ID: uint(company.ID),
			},
			Name:             company.Name,
			Description:      &company.Description,
			ChangedCompanyID: changedCompanyId,
			Slug:             company.Slug,
			StartDate:        time.Unix(int64(company.StartDate), 0),
			IGDBUrl:          company.URL,
		}).Error
	}
	return nil
}

func fetchGenres() error {
	genreCount, err := client.Genres.Count()
	if err != nil {
		return err
	}
	for i := 0; i < genreCount; i += 500 {
		opts := igdb.ComposeOptions(
			igdb.SetLimit(500),
			igdb.SetFields("*"),
			igdb.SetOrder("id", igdb.OrderAscending),
			igdb.SetOffset(i),
		)
		genres, err := client.Genres.Index(opts)
		if err != nil {
			return err

		}
		for _, genre := range genres {
			dbGenre := schema.Genre{
				EnumCategory: schema.EnumCategory{
					ID:   uint(genre.ID),
					Name: genre.Name,
				},
			}
			if create := db.Create(&dbGenre); create.Error != nil {
				return create.Error
			}
		}
	}
	return nil
}

func fetchGames(offset int) error {
	gamesCount, err := client.Games.Count()
	if err != nil {
		return err
	}
	for i := offset; i < gamesCount; i += 500 {
		opts := igdb.ComposeOptions(
			igdb.SetLimit(500),
			igdb.SetFields("age_ratings.*", "artworks.*", "cover.*", "external_games.*", "involved_companies.*",
				"platforms.*", "platforms.platform_logo.*", "release_dates.*", "screenshots.*", "videos.*", "*"),
			igdb.SetOrder("id", igdb.OrderAscending),
			igdb.SetOffset(i),
		)
		logger.Println("Fetching games", i, "to", i+500)
		igdbGames, err := client.Games.Index(opts)
		if err != nil {
			return err
		}

		for _, game := range igdbGames {
			result := db.First(&schema.Game{}, game.ID)
			if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
				var parentGameId *uint = nil
				var versionParentId *uint = nil
				// get parent game if exists
				if game.ParentGame != nil && game.ParentGame.ID != game.ID {
					parentGameId = new(uint)
					*parentGameId = uint(game.ParentGame.ID)
					if err := fetchGame(game.ParentGame.ID); err != nil {
						if errors.Is(err, igdb.ErrNoResults) {
							parentGameId = nil
						} else {
							return err
						}
					}
				}
				// get version parent game if exists
				if game.VersionParent != nil && game.VersionParent.ID != game.ID {
					versionParentId = new(uint)
					*versionParentId = uint(game.VersionParent.ID)
					if err := fetchGame(game.VersionParent.ID); err != nil {
						if errors.Is(err, igdb.ErrNoResults) {
							versionParentId = nil
						} else {
							return err
						}
					}
				}
				dbGame := createBaseGameObject(game, parentGameId, versionParentId)
				if create := db.Create(&dbGame); create.Error != nil {
					return create.Error
				}
				// dlc
				for _, dlc := range game.DLCS {
					dlcGame, err := prepareGame(dlc.ID, true)
					if err != nil {
						if errors.Is(err, igdb.ErrNoResults) {
							continue
						}
						return err
					}
					parentGameIdRef := new(uint)
					*parentGameIdRef = uint(game.ID)
					dlcGame.DLCBaseReference = parentGameIdRef
					if create := db.Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "id"}},
						DoUpdates: clause.AssignmentColumns([]string{"dlc_base_reference"}),
					}).Create(dlcGame); create.Error != nil {
						return create.Error
					}
					if err != nil {
						return err
					}
				}
				// expanded
				for _, expanded := range game.ExpandedGames {
					expandedGame, err := prepareGame(expanded.ID, true)
					if err != nil {
						if errors.Is(err, igdb.ErrNoResults) {
							continue
						}
						return err
					}
					parentGameIdRef := new(uint)
					*parentGameIdRef = uint(game.ID)
					expandedGame.ExpandedGameReference = parentGameIdRef
					if create := db.Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "id"}},
						DoUpdates: clause.AssignmentColumns([]string{"expanded_game_reference"}),
					}).Create(expandedGame); create.Error != nil {
						return create.Error
					}
					if err != nil {
						return err
					}
				}
				// expansions
				for _, expansion := range game.Expansions {
					expansionGame, err := prepareGame(expansion.ID, true)
					if err != nil {
						if errors.Is(err, igdb.ErrNoResults) {
							continue
						}
						return err
					}
					parentGameIdRef := new(uint)
					*parentGameIdRef = uint(game.ID)
					expansionGame.ExpansionReference = parentGameIdRef
					if create := db.Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "id"}},
						DoUpdates: clause.AssignmentColumns([]string{"expansion_reference"}),
					}).Create(expansionGame); create.Error != nil {
						return create.Error
					}
					if err != nil {
						return err
					}
				}
				// remakes
				for _, remake := range game.Remakes {
					remakeGame, err := prepareGame(remake.ID, true)
					if err != nil {
						if errors.Is(err, igdb.ErrNoResults) {
							continue
						}
						return err
					}
					parentGameIdRef := new(uint)
					*parentGameIdRef = uint(game.ID)
					remakeGame.RemakeBaseReference = parentGameIdRef
					if create := db.Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "id"}},
						DoUpdates: clause.AssignmentColumns([]string{"remake_base_reference"}),
					}).Create(remakeGame); create.Error != nil {
						return create.Error
					}
					if err != nil {
						return err
					}
				}
				// remasters
				for _, remaster := range game.Remasters {
					remasterGame, err := prepareGame(remaster.ID, true)
					if err != nil {
						if errors.Is(err, igdb.ErrNoResults) {
							continue
						}
						return err
					}
					parentGameIdRef := new(uint)
					*parentGameIdRef = uint(game.ID)
					remasterGame.RemasterBaseReference = parentGameIdRef
					if create := db.Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "id"}},
						DoUpdates: clause.AssignmentColumns([]string{"remaster_base_reference"}),
					}).Create(remasterGame); create.Error != nil {
						return create.Error
					}
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func fetchGame(gameId int) error {
	dbGame, err := prepareGame(gameId, false)
	if err != nil {
		return err
	}
	if dbGame.ParentGameID != nil {
		if err := fetchGame(int(*dbGame.ParentGameID)); err != nil {
			if errors.Is(err, igdb.ErrNoResults) {
				dbGame.ParentGameID = nil
			} else {
				return err
			}
		}
	}
	if dbGame.VersionParentID != nil {
		if err := fetchGame(int(*dbGame.VersionParentID)); err != nil {
			if errors.Is(err, igdb.ErrNoResults) {
				dbGame.VersionParentID = nil
			} else {
				return err
			}
		}
	}
	if create := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&dbGame); create.Error != nil {
		return create.Error
	}
	return nil
}

func prepareGame(gameId int, getParentVersion bool) (*schema.Game, error) {
	dbGame := schema.Game{}
	result := db.First(&dbGame, gameId)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var parentGameId *uint = nil
		var versionParentId *uint = nil
		opts := igdb.ComposeOptions(
			igdb.SetFields("age_ratings.*", "artworks.*", "cover.*", "external_games.*", "involved_companies.*",
				"platforms.*", "platforms.platform_logo.*", "release_dates.*", "screenshots.*", "videos.*", "*"),
		)
		game, err := client.Games.Get(gameId, opts)
		if err != nil {
			if !errors.Is(err, igdb.ErrNoResults) {
				println("Jebło")
				game, err = client.Games.Get(gameId, opts)
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
				if err := fetchGame(game.VersionParent.ID); err != nil {
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
	return slice_funcs.Filter(externalGames, func(externalGame igdb.ExternalGameWrapper) bool {
		for _, category := range externalCategories {
			if uint(externalGame.Category) == category.ID {
				return true
			}
		}
		return false
	})
}

func associateSimilarGames() error {
	gamesCount, err := client.Games.Count()
	fmt.Println("Games count: ", gamesCount)
	if err != nil {
		return err
	}
	for i := 0; i < gamesCount; i += 500 {
		fmt.Println("Processing games from ", i, " to ", i+500)
		opts := igdb.ComposeOptions(
			igdb.SetLimit(500),
			igdb.SetFields("similar_games"),
			igdb.SetOrder("id", igdb.OrderAscending),
			igdb.SetOffset(i),
		)
		igdbGames, err := client.Games.Index(opts)
		if err != nil {
			return err
		}
		for _, igdbGame := range igdbGames {
			var game schema.Game
			dbGame := db.Model(&game).Where("id = ?", igdbGame.ID).First(&game)
			if dbGame.Error != nil {
				if errors.Is(dbGame.Error, gorm.ErrRecordNotFound) {
					continue
				}
				return dbGame.Error
			}
			association := db.Model(&game).Association("SimilarGames")
			if association.Error != nil {
				return association.Error
			}
			for _, similarGame := range igdbGame.SimilarGames {
				var similarGameModel schema.Game
				existingGameCheck := db.Model(&game).Where("id = ?", similarGame.ID).First(&similarGameModel)
				if existingGameCheck.Error != nil {
					if errors.Is(existingGameCheck.Error, gorm.ErrRecordNotFound) {
						continue
					}
					return existingGameCheck.Error
				}
				var test schema.Game
				existingAssociation := db.Model(&game).Where("id = ?", similarGame.ID).Association("SimilarGames").Find(&test)
				if existingAssociation != nil {
					if !errors.Is(existingAssociation, gorm.ErrRecordNotFound) {
						return existingAssociation
					}
				}
				if test.ID != 0 {
					continue
				}
				existingAssociation = db.Model(&similarGameModel).Where("id = ?", game.ID).Association("SimilarGames").Find(&test)
				if existingAssociation != nil {
					if !errors.Is(existingAssociation, gorm.ErrRecordNotFound) {
						return existingAssociation
					}
				}
				if test.ID != 0 {
					continue
				}
				err := association.Append(&schema.Game{
					Model: gorm.Model{
						ID: uint(similarGame.ID),
					},
				})
				if err != nil {
					fmt.Println("Error appending similar game: ", err)
				}
			}
		}
	}
	return nil
}

func clearDatabase() {
	//if result := db.Exec("DELETE FROM game_platforms"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM platform_logos"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM platforms"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM release_dates"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM artworks"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	////db.Exec("DELETE FROM companies")
	////db.Exec("DELETE FROM company_logos")
	//if result := db.Exec("DELETE FROM covers"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM external_games"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM game_age_ratings"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM game_genres"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM game_similar_games"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM game_videos"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM games"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	////db.Exec("DELETE FROM genres")
	//if result := db.Exec("DELETE FROM involved_companies"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	//if result := db.Exec("DELETE FROM screenshots"); result.Error != nil {
	//	logger.Fatal(result.Error)
	//}
	if result := db.Exec("DELETE FROM game_similar_games"); result.Error != nil {
		logger.Fatal(result.Error)
	}
}
