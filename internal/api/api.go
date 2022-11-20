package api

import (
	"github.com/joho/godotenv"
	"github.com/kwisnia/inzynierka-backend/internal/api/router"
	"github.com/kwisnia/inzynierka-backend/internal/api/websocket"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/awsconf"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/igdb"
	"log"
)

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := router.Setup()
	database.SetupDB()
	igdb.SetupClient()
	igdb.SetupWebhooks()
	awsconf.ConnectAws()
	go websocket.ClientHub.Run()
	r.Run()
}
