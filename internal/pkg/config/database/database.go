package database

import (
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func SetupDB() {
	dsn := config.GetEnv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer,
	//	logger.Config{
	//		SlowThreshold:             time.Second, // Slow SQL threshold
	//		LogLevel:                  logger.Info, // Log level
	//		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	//		Colorful:                  false,       // Disable color
	//	},
	//)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
}
