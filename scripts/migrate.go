package main

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/user"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
)

func main() {
	config.LoadConfig()
	database.SetupDB()
	db := database.DB
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&user.UserProfile{})
}
