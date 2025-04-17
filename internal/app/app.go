package app

import (
	"fmt"
	"net/http"
	"orderPickupPoint/config"
	"orderPickupPoint/internal/service"
	"orderPickupPoint/internal/storage"
	"orderPickupPoint/internal/storage/postgres"
	"orderPickupPoint/internal/transport"
)

func Run() {
	fmt.Println("starting app")

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("something wrong with config")
	}

	dbConnPool, err := postgres.InitDb()
	if err != nil {
		fmt.Println("something wrong with database")
	}

	repos := storage.NewRepositories(dbConnPool)
	services := service.NewServices(&service.Deps{
		Repos: repos,
		Cfg:   cfg,
	})
	handler := transport.NewHandler(services)
	router := handler.InitRouter()

	http.ListenAndServe(cfg.ServerAddress, router)
}
