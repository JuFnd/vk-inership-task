package main

import (
	"filmoteka/configs"
	"filmoteka/modules/films/delivery"
	"filmoteka/modules/films/repository"
	"filmoteka/modules/films/usecase"
	"filmoteka/pkg/variables"
	"fmt"
	"log/slog"
	"os"
)

// @title Films service
// @version 1.0
// @description VK Filmoteka films service

// @host localhost:8081
// @BasePath /

// @in header
// @name Films

func main() {
	logFile, err := os.Create("films.log")
	if err != nil {
		fmt.Println("Error creating log file")
		return
	}

	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	configFilms, err := configs.ReadFilmsAppConfig()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	relationalDataBaseConfig, err := configs.ReadRelationalFilmsDataBaseConfig()
	if err != nil {
		logger.Error(variables.ReadFilmsSqlConfigError, err.Error())
		return
	}

	grpcConfig, err := configs.ReadGrpcConfig()

	filmsRepository, err := repository.GetFilmRepository(*relationalDataBaseConfig, logger)
	core := usecase.GetCore(*grpcConfig, filmsRepository, logger)
	if err != nil {
		logger.Error(variables.CoreInitializeError, err)
		return
	}

	api := delivery.GetFilmsApi(core, logger)

	err = api.ListenAndServe(configFilms)
	if err != nil {
		logger.Error(err.Error())
	}
}
