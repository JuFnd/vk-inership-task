package main

import (
	"filmoteka/configs"
	delivery_grpc "filmoteka/modules/authorization/delivery/grpc"
	"filmoteka/modules/authorization/delivery/http"
	"filmoteka/modules/authorization/usecase"
	"filmoteka/pkg/variables"
	"fmt"
	"log/slog"
	"os"
)

// @title Authorization service
// @version 1.0
// @description VK Filmoteka authorization service

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

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

	relationalDataBaseConfig, err := configs.ReadRelationalAuthDataBaseConfig()
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

	grpcServer, err := delivery_grpc.NewServer(relationalDataBaseConfig, cacheDatabaseConfig, logger)
	if err != nil {
		logger.Error(variables.ListenAndServeError)
		return
	}

	api := delivery.GetAuthorizationApi(core, logger)

	errs := make(chan error, 2)
	go func() {
		errs <- api.ListenAndServe(authAppConfig)
	}()

	go func() {
		errs <- grpcServer.ListenAndServeGrpc()
	}()

	err = <-errs
	if err != nil {
		logger.Error(variables.ListenAndServeError, err.Error())
	}
}
