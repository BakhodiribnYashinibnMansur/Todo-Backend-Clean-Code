package main

import (
	"log"

	"todocc/config"
	"todocc/package/handler"
	"todocc/package/repository"
	"todocc/package/service"
	"todocc/server"

	_ "github.com/lib/pq"
)

func main() {

	configs, err := config.InitConfig()
	if err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     configs.Host,
		Port:     configs.Port,
		Username: configs.Username,
		DBName:   configs.DBName,
		SSLMode:  configs.SSLMode,
		Password: configs.Password,
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}
