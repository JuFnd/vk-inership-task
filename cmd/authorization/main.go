package main

import (
	"filmoteka/configs"
	"filmoteka/modules/authorization/delivery"
	"filmoteka/modules/authorization/usecase"
	"filmoteka/pkg/variables"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logFile, err := os.Create("authorization.log")
	if err != nil {
		fmt.Println("Error creating log file")
		return
	}

	logger := slog.New(slog.NewJSONHandler(logFile, nil))
	authAppConfig, err := configs.ReadAuthAppConfig()
	if err != nil {
		logger.Error(variables.ReadAuthConfigError, err.Error())
		return
	}

	relationalDataBaseConfig, err := configs.ReadRelationalDataBaseConfig()
	if err != nil {
		logger.Error(variables.ReadAuthSqlConfigError, err.Error())
		return
	}

	cacheDatabaseConfig, err := configs.ReadCacheDatabaseConfig()
	if err != nil {
		logger.Error(variables.ReadAuthCacheConfigError, err.Error())
		return
	}

	core, err := usecase.GetCore(relationalDataBaseConfig, cacheDatabaseConfig, logger)
	if err != nil {
		logger.Error(variables.CoreInitializeError, err)
		return
	}

	api := delivery.GetAuthorizationApi(core, logger)

	errs := make(chan error, 2)
	go func() {
		errs <- api.ListenAndServe(authAppConfig)
	}()

	err = <-errs
	if err != nil {
		logger.Error(variables.ListenAndServeError, err.Error())
	}
}
