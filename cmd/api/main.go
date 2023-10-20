package main

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pDb "github.com/SlavaShagalov/my-trello-backend/internal/pkg/db"
	pHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher/bcrypt"
	pLog "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	"log"
	"net/http"
	"os"

	sessionsRepository "github.com/SlavaShagalov/my-trello-backend/internal/sessions/repository/redis"
	usersRepository "github.com/SlavaShagalov/my-trello-backend/internal/users/repository/postgres"
	workspacesRepository "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/repository/postgres"

	authUsecase "github.com/SlavaShagalov/my-trello-backend/internal/auth/usecase"
	workspacesUsecase "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/usecase"

	authDel "github.com/SlavaShagalov/my-trello-backend/internal/auth/delivery/http"
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"
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
	db, err := pDb.NewPostgres(logger)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Error("Failed to close DB connection", zap.Error(err))
		}
		logger.Info("DB connection closed")
	}()

	// Sessions Storage
	rdb, err := pDb.NewRedis(logger, context.Background())
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

	// Hasher
	hasher := pHasher.NewHasher()

	// Repo
	sessionsRepo := sessionsRepository.New(rdb, context.Background(), logger)
	usersRepo := usersRepository.New(db, logger)
	workspacesRepo := workspacesRepository.New(db, logger)
	//boardsRepo := boardsRepository.New(db, logger)
	//listsRepo := listsRepository.New(db, logger)
	//cardsRepo := cardsRepository.New(db, logger)

	// Use cases
	authUC := authUsecase.New(usersRepo, sessionsRepo, hasher, logger)
	workspacesUC := workspacesUsecase.New(workspacesRepo)
	//boardsUC := boardsUsecase.New(boardsRepo)
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
	workspacesDel.RegisterHandlers(router, workspacesUC, logger, checkAuth)

	server := http.Server{
		Addr:    constants.ApiAddress,
		Handler: router,
	}

	logger.Info("API service started at", zap.String("address", constants.ApiAddress))
	if err = server.ListenAndServe(); err != nil {
		logger.Error("API server stopped %v", zap.Error(err))
	}
}
