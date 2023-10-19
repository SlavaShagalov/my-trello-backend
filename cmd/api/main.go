package main

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pkgDb "github.com/SlavaShagalov/my-trello-backend/internal/pkg/db"
	pkgHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher/bcrypt"
	pkgLog "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	"log"
	"net/http"
	"os"

	sessionsRepository "github.com/SlavaShagalov/my-trello-backend/internal/sessions/repository/redis"
	usersRepository "github.com/SlavaShagalov/my-trello-backend/internal/users/repository/postgres"

	authUsecase "github.com/SlavaShagalov/my-trello-backend/internal/auth/usecase"

	authDel "github.com/SlavaShagalov/my-trello-backend/internal/auth/delivery/http"
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Logger
	logger, logfile, err := pkgLog.NewProdLogger("/logs/api.log")
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
	db, err := pkgDb.NewPostgres(logger)
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
	rdb, err := pkgDb.NewRedis(logger, context.Background())
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
	hasher := pkgHasher.NewHasher()

	// Repo
	sessionsRepo := sessionsRepository.NewRepository(rdb, context.Background(), logger)
	usersRepo := usersRepository.NewRepository(db, logger)
	//workspacesRepo := workspacesRepository.NewRepository(db, logger)
	//boardsRepo := boardsRepository.NewRepository(db, logger)
	//listsRepo := listsRepository.NewRepository(db, logger)
	//cardsRepo := cardsRepository.NewRepository(db, logger)

	// Use cases
	authUC := authUsecase.NewUsecase(usersRepo, sessionsRepo, hasher, logger)
	//boardsUC := boardsUsecase.NewUsecase(boardsRepo)
	//listsUC := listsUsecase.NewUsecase(listsRepo)
	//cardsUC := cardsUsecase.NewUsecase(cardsRepo)

	// Middleware
	checkAuth := mw.NewCheckAuth(authUC, logger)

	// Delivery
	router := mux.NewRouter()
	router.Use(
		mw.NewAccessLog(logger),
	)

	authDel.RegisterHandlers(router, authUC, logger, checkAuth)

	server := http.Server{
		Addr:    constants.ApiAddress,
		Handler: router,
	}

	logger.Info("API service started at", zap.String("address", constants.ApiAddress))
	if err = server.ListenAndServe(); err != nil {
		logger.Error("API server stopped %v", zap.Error(err))
	}
}
