package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rawfish-dev/angrypros-api/config"
	"github.com/rawfish-dev/angrypros-api/handlers"
	"github.com/rawfish-dev/angrypros-api/services/auth"
	"github.com/rawfish-dev/angrypros-api/services/feed"
	"github.com/rawfish-dev/angrypros-api/services/storage"
	timeS "github.com/rawfish-dev/angrypros-api/services/time"
)

func main() {
	appConfig := config.NewAppConfig(os.Getenv("APP_ENVIRONMENT"), ".")

	authService, err := auth.NewService(appConfig.GoogleConfig)
	if err != nil {
		panic(fmt.Sprintf("could not initialise auth service due to %s", err))
	}

	storageService, err := storage.NewService(appConfig.PostgresConfig)
	if err != nil {
		panic(fmt.Sprintf("could not initialise storage service due to %s", err))
	}

	feedService := feed.NewService(storageService)
	if err != nil {
		panic(fmt.Sprintf("could not initialise auth service due to %s", err))
	}

	timeService := timeS.NewService()

	server, err := handlers.NewServer(appConfig, authService,
		feedService, storageService, timeService)
	if err != nil {
		panic(fmt.Sprintf("could not initialise server due to %s", err))
	}

	server.SetupRoutes()
	http.ListenAndServe(":8080", server)
}
