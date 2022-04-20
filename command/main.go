package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"todocc/config"
	"todocc/package/handler"
	"todocc/package/repository"
	"todocc/package/service"
	"todocc/server"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// logger := logrus.GetLogger()

	configs, err := config.InitConfig()
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     configs.DBHost,
		Port:     configs.DBPort,
		Username: configs.DBUsername,
		DBName:   configs.DBName,
		SSLMode:  configs.DBSSLMode,
		Password: configs.DBPassword,
	})

	logrus.Infof("configs: %v", configs)

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
