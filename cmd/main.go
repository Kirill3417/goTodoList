package main

import (
	"awesomeProject"
	"awesomeProject/pkg/handler"
	"awesomeProject/pkg/repository"
	"awesomeProject/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка инициализации конфигов: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.name"),
		viper.GetInt("db.port"),
	)

	if err != nil {
		logrus.Fatalf("Ошибка инициализации БД: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(todo.Server)
	if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Произошла ошибка при запуске сервера: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
