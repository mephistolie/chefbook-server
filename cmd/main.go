package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mephistolie/chefbook-server/internal/app/chefbook"
	"github.com/mephistolie/chefbook-server/internal/handlers"
	"github.com/mephistolie/chefbook-server/internal/repositories"
	"github.com/mephistolie/chefbook-server/internal/services"
	"github.com/siruspen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variable: %s", err.Error())
	}

	db, err := repositories.NewPostgresDB(repositories.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repository := repositories.NewRepository(db)
	service := services.NewService(repository)
	handler := handlers.NewHandler(service)

	srv := new(chefbook.Server)
	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
