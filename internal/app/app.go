package app

import (
	"context"
	"fmt"
	"net/http"
	"orderPickupPoint/config"
	"orderPickupPoint/internal/service"
	"orderPickupPoint/internal/storage"
	"orderPickupPoint/internal/storage/postgres"
	"orderPickupPoint/internal/transport"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("something wrong with config")
	}

	dbConnPool, err := postgres.InitDb()
	if err != nil {
		fmt.Println("something wrong with database")
	}
	defer dbConnPool.Close()

	repos := storage.NewRepositories(dbConnPool)
	services := service.NewServices(&service.Deps{
		Repos: repos,
		Cfg:   cfg,
	})
	handler := transport.NewHandler(services)
	router := handler.InitRouter()

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("server error: ", err)
		}
	}()
	http.ListenAndServe(cfg.ServerAddress, router)

	fmt.Println("server started on", cfg.ServerAddress)

	<-quit
	fmt.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("server shutdown failed: %v\n", err)
	} else {
		fmt.Println("server exited properly")
	}

}
