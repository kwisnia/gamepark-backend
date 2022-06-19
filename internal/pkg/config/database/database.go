package database

import (
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB() {
	databaseConfig := config.Config.Database
	dsn := "user=" + databaseConfig.Username + " password=" + databaseConfig.Password + " dbname=" + databaseConfig.Dbname + " host=" + databaseConfig.Host + " port=" + databaseConfig.Port + " sslmode=disable TimeZone=Europe/Warsaw"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
}

