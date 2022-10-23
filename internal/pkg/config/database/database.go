package database

import (
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func SetupDB() {
	databaseConfig := config.Config.Database
	dsn := "user=" + databaseConfig.Username + " password=" + databaseConfig.Password + " dbname=" + databaseConfig.Dbname + " host=" + databaseConfig.Host + " port=" + databaseConfig.Port + " sslmode=disable TimeZone=Europe/Warsaw"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer,
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	DB = db
}
