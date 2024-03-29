package main

import (
	"context"
	"github.com/SlavaShagalov/my-trello-backend/internal/boards"
	boardsRepositoryPgx "github.com/SlavaShagalov/my-trello-backend/internal/boards/repository/pgx"
	boardsRepository "github.com/SlavaShagalov/my-trello-backend/internal/boards/repository/std"
	"github.com/SlavaShagalov/my-trello-backend/internal/cards"
	cardsRepository "github.com/SlavaShagalov/my-trello-backend/internal/cards/repository/postgres"
	imagesRepository "github.com/SlavaShagalov/my-trello-backend/internal/images/repository/s3"
	"github.com/SlavaShagalov/my-trello-backend/internal/lists"
	listsRepository "github.com/SlavaShagalov/my-trello-backend/internal/lists/repository/postgres"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/constants"
	pHasher "github.com/SlavaShagalov/my-trello-backend/internal/pkg/hasher/bcrypt"
	pLog "github.com/SlavaShagalov/my-trello-backend/internal/pkg/log/zap"
	pMetrics "github.com/SlavaShagalov/my-trello-backend/internal/pkg/metrics"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/opentel"
	pStorages "github.com/SlavaShagalov/my-trello-backend/internal/pkg/storages"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/storages/postgres"
	sessionsRepository "github.com/SlavaShagalov/my-trello-backend/internal/sessions/repository/redis"
	"github.com/SlavaShagalov/my-trello-backend/internal/users"
	usersRepository "github.com/SlavaShagalov/my-trello-backend/internal/users/repository/postgres"
	"github.com/SlavaShagalov/my-trello-backend/internal/workspaces"
	workspacesRepository "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/repository/postgres"
	"log"
	"net/http"
	"os"

	authUsecase "github.com/SlavaShagalov/my-trello-backend/internal/auth/usecase"
	boardsUsecase "github.com/SlavaShagalov/my-trello-backend/internal/boards/usecase"
	cardsUsecase "github.com/SlavaShagalov/my-trello-backend/internal/cards/usecase"
	listsUsecase "github.com/SlavaShagalov/my-trello-backend/internal/lists/usecase"
	usersUsecase "github.com/SlavaShagalov/my-trello-backend/internal/users/usecase"
	workspacesUsecase "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/usecase"

	authDel "github.com/SlavaShagalov/my-trello-backend/internal/auth/delivery/http"
	boardsDel "github.com/SlavaShagalov/my-trello-backend/internal/boards/delivery/http"
	cardsDel "github.com/SlavaShagalov/my-trello-backend/internal/cards/delivery/http"
	listsDel "github.com/SlavaShagalov/my-trello-backend/internal/lists/delivery/http"
	mw "github.com/SlavaShagalov/my-trello-backend/internal/middleware"
	usersDel "github.com/SlavaShagalov/my-trello-backend/internal/users/delivery/http"
	workspacesDel "github.com/SlavaShagalov/my-trello-backend/internal/workspaces/delivery/http"

	_ "github.com/SlavaShagalov/my-trello-backend/docs"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

// main godoc
//
//	@title						MyTrello API
//
//	@version					1.0
//	@description				MyTrello API documentation.
//	@termsOfService				http://127.0.0.1/terms
//
//	@contact.name				MyTrello API Support
//	@contact.url				http://127.0.0.1/support
//	@contact.email				my-trello-support@vk.com
//
//	@host						127.0.0.1
//	@BasePath					/api/v1
//	@securitydefinitions.apikey	cookieAuth
//
//	@in							cookie
//	@name						JSESSIONID
func main() {
	ctx := context.Background()

	// ===== Configuration =====
	config.SetDefaultPostgresConfig()
	config.SetDefaultRedisConfig()
	config.SetDefaultS3Config()
	config.SetDefaultValidationConfig()
	viper.SetConfigName("api")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/configs")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Failed to read configuration: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Configuration read successfully")

	// ===== Logger =====
	logger, logfile, err := pLog.NewProdLogger("/logs/" + viper.GetString(config.ServerName) + ".log")
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

	// ===== OpenTelemetry =====
	serviceName := viper.GetString(config.ServerName)
	logger.Info("", zap.String("server_name", serviceName))
	serviceVersion := "0.1.0"
	tp, mp, err := opentel.SetupOTelSDK(context.Background(), logger, serviceName, serviceVersion)
	if err != nil {
		return
	}
	defer func() { _ = tp.Shutdown(ctx) }()
	defer func() { _ = mp.Shutdown(ctx) }()

	// ===== Data Storage =====
	db, err := postgres.NewStd(logger)
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

	// ===== Sessions Storage =====
	redisClient, err := pStorages.NewRedis(logger, ctx)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		err = redisClient.Close()
		if err != nil {
			logger.Error("Failed to close Redis client", zap.Error(err))
		}
		logger.Info("Redis client closed")
	}()

	// ===== S3 =====
	s3Client, err := pStorages.NewS3(logger)
	if err != nil {
		os.Exit(1)
	}

	// ===== Prometheus =====
	mt := pMetrics.NewPrometheusMetrics("api")
	err = mt.SetupMetrics()
	if err != nil {
		logger.Error("failed to setup prometheus", zap.Error(err))
		os.Exit(1)
	}

	// ===== Hasher =====
	hasher := pHasher.New()

	// ===== Repositories =====
	var usersRepo users.Repository
	var workspacesRepo workspaces.Repository
	var boardsRepo boards.Repository
	var listsRepo lists.Repository
	var cardsRepo cards.Repository
	usersRepo = usersRepository.New(db, logger)
	workspacesRepo = workspacesRepository.New(db, logger)
	listsRepo = listsRepository.New(db, logger)
	cardsRepo = cardsRepository.New(db, logger)

	serverType := viper.GetString(config.ServerType)

	mode := "std"
	if mode == "std" {
		boardsRepo = boardsRepository.New(db, logger)
	} else if mode == "pgx" {
		pgxPool, err := postgres.NewPgx(logger)
		if err != nil {
			os.Exit(1)
		}

		boardsRepo = boardsRepositoryPgx.New(pgxPool, logger)
	}

	imagesRepo := imagesRepository.New(s3Client, logger)
	sessionsRepo := sessionsRepository.New(redisClient, context.Background(), logger)

	// ===== Usecases =====
	authUC := authUsecase.New(usersRepo, sessionsRepo, hasher, logger)
	usersUC := usersUsecase.New(usersRepo, imagesRepo)
	workspacesUC := workspacesUsecase.New(workspacesRepo)
	boardsUC := boardsUsecase.New(boardsRepo, imagesRepo)
	listsUC := listsUsecase.New(listsRepo)
	cardsUC := cardsUsecase.New(cardsRepo)

	// ===== Middleware =====
	checkAuth := mw.NewCheckAuth(authUC, logger)
	accessLog := mw.NewAccessLog(serverType, logger)
	cors := mw.NewCors()
	metrics := mw.NewMetrics(mt)

	router := mux.NewRouter()

	// ===== Delivery =====
	authDel.RegisterHandlers(router, authUC, usersUC, logger, checkAuth, metrics)
	usersDel.RegisterHandlers(router, usersUC, logger, checkAuth, metrics)
	workspacesDel.RegisterHandlers(router, workspacesUC, boardsUC, logger, checkAuth, metrics)
	boardsDel.RegisterHandlers(router, boardsUC, logger, checkAuth, metrics)
	listsDel.RegisterHandlers(router, listsUC, cardsUC, logger, checkAuth, metrics)
	cardsDel.RegisterHandlers(router, cardsUC, logger, checkAuth, metrics)

	// ===== Swagger =====
	router.PathPrefix(constants.ApiPrefix + "/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	// ===== Router =====
	server := http.Server{
		Addr:    ":" + viper.GetString(config.ServerPort),
		Handler: accessLog(cors(router)),
	}

	logger.Info("Starting metrics...", zap.String("address", "0.0.0.0:9001"))
	go pMetrics.ServePrometheusHTTP("0.0.0.0:9001")
	logger.Info("Metrics started")

	// ===== Start =====
	logger.Info("API service started", zap.String("port", viper.GetString(config.ServerPort)))
	if err = server.ListenAndServe(); err != nil {
		logger.Error("API server stopped", zap.Error(err))
	}
}
