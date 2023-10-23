package main

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher/bcrypt"
	pLog "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	pStorages "github.com/SlavaShagalov/my-trello-backend/internal/pkg/storages"
	"log"
	"net/http"
	"os"

	boardsRepository "github.com/SlavaShagalov/my-trello-backend/internal/boards/repository/postgres"
	imagesRepository "github.com/SlavaShagalov/my-trello-backend/internal/images/repository/s3"
	sessionsRepository "github.com/SlavaShagalov/my-trello-backend/internal/sessions/repository/redis"
	usersRepository "github.com/SlavaShagalov/my-trello-backend/internal/users/repository/postgres"
	workspacesRepository "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/repository/postgres"

	authUsecase "github.com/SlavaShagalov/my-trello-backend/internal/auth/usecase"
	boardsUsecase "github.com/SlavaShagalov/my-trello-backend/internal/boards/usecase"
	usersUsecase "github.com/SlavaShagalov/my-trello-backend/internal/users/usecase"
	workspacesUsecase "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/usecase"

	authDel "github.com/SlavaShagalov/my-trello-backend/internal/auth/delivery/http"
	boardsDel "github.com/SlavaShagalov/my-trello-backend/internal/boards/delivery/http"
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"
	usersDel "github.com/SlavaShagalov/my-trello-backend/internal/users/delivery/http"
	workspacesDel "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/delivery/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Logger
	logger, logfile, err := pLog.NewProdLogger("/logs/api.log")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Println(err)
		}
		err = logfile.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	logger.Info("API service starting...")

	// Config
	config.SetDefaultPostgresConfig()
	config.SetDefaultRedisConfig()
	config.SetDefaultS3Config()
	config.SetDefaultValidationConfig()
	viper.SetConfigName("api")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/configs")
	err = viper.ReadInConfig()
	if err != nil {
		logger.Error("Failed to read configuration", zap.Error(err))
		os.Exit(1)
	}

	// Data Storage
	db, err := pStorages.NewPostgres(logger)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Error("Failed to close Postgres connection", zap.Error(err))
		}
		logger.Info("Postgres connection closed")
	}()

	// Sessions Storage
	rdb, err := pStorages.NewRedis(logger, context.Background())
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		err = rdb.Close()
		if err != nil {
			logger.Error("Failed to close Redis client", zap.Error(err))
		}
		logger.Info("Redis client closed")
	}()

	// S3
	s3Client, err := pStorages.NewS3(logger)
	if err != nil {
		os.Exit(1)
	}

	// Hasher
	hasher := pHasher.NewHasher()

	// Repo
	imagesRepo := imagesRepository.New(s3Client, logger)
	sessionsRepo := sessionsRepository.New(rdb, context.Background(), logger)
	usersRepo := usersRepository.New(db, logger)
	workspacesRepo := workspacesRepository.New(db, logger)
	boardsRepo := boardsRepository.New(db, logger)
	//listsRepo := listsRepository.New(storages, logger)
	//cardsRepo := cardsRepository.New(storages, logger)

	// Use cases
	authUC := authUsecase.New(usersRepo, sessionsRepo, hasher, logger)
	usersUC := usersUsecase.New(usersRepo, imagesRepo)
	workspacesUC := workspacesUsecase.New(workspacesRepo)
	boardsUC := boardsUsecase.NewUsecase(boardsRepo, imagesRepo)
	//listsUC := listsUsecase.New(listsRepo)
	//cardsUC := cardsUsecase.New(cardsRepo)

	// Middleware
	checkAuth := mw.NewCheckAuth(authUC, logger)

	router := mux.NewRouter()
	router.Use(
		mw.NewAccessLog(logger),
	)

	// Delivery
	authDel.RegisterHandlers(router, authUC, logger, checkAuth)
	usersDel.RegisterHandlers(router, usersUC, logger, checkAuth)
	workspacesDel.RegisterHandlers(router, workspacesUC, logger, checkAuth)
	boardsDel.RegisterHandlers(router, boardsUC, logger, checkAuth)

	server := http.Server{
		Addr:    constants.ApiAddress,
		Handler: router,
	}

	logger.Info("API service started at", zap.String("address", constants.ApiAddress))
	if err = server.ListenAndServe(); err != nil {
		logger.Error("API server stopped %v", zap.Error(err))
	}
}
