package main

import (
	"todocc/config"
	"todocc/package/handler"
	"todocc/package/repository"
	"todocc/package/service"
	"todocc/server"
	"todocc/util/logrus"

	_ "github.com/lib/pq"
)

func main() {
	logger := logrus.GetLogger()

	configs, err := config.InitConfig()
	if err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     configs.DBHost,
		Port:     configs.DBPort,
		Username: configs.DBUsername,
		DBName:   configs.DBName,
		SSLMode:  configs.DBSSLMode,
		Password: configs.DBPassword,
	})

	logger.Infof("configs: %v", configs)

	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(configs.ServerPort, handlers.InitRoutes()); err != nil {
		logger.Fatalf("error occurred while running http server: %s", err.Error())
	}
}
