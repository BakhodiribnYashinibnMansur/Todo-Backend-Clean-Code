package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Configs struct {
	ServerPort string `default:"8080"`
	DBHost     string `default:"localhost"`
	DBPort     string `default:"2000"`
	DBUsername string `default:"postgres"`
	DBName     string `default:"todocc"`
	DBPassword string `default:"0224"`
	DBSSLMode  string `default:"false"`
}

func InitConfig() (cnfg *Configs, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()

	if err != nil {
		return cnfg, fmt.Errorf("fatal error config file: %w ", err)
	}

	if err := godotenv.Load(); err != nil {
		return cnfg, fmt.Errorf("error loading env variables: %s", err.Error())
	}

	cnfg = &Configs{
		ServerPort: viper.GetString("port"),
		DBHost:     viper.GetString("db.host"),
		DBPort:     viper.GetString("db.port"),
		DBUsername: viper.GetString("db.username"),
		DBName:     viper.GetString("db.dbname"),
		DBSSLMode:  viper.GetString("db.sslmode"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
	return
}
