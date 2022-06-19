package api

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/router"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config"
)

func Run()  {
	config.LoadConfig()
	r := router.Setup()
	r.Run()
}