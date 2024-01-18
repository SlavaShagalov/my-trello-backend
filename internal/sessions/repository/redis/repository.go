package redis

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	pkgErrors "github.com/SlavaShagalov/my-trello-backend/internal/pkg/errors"
	pkgSessions "github.com/SlavaShagalov/my-trello-backend/internal/sessions"
)

const (
	componentName = "Sessions Repository"
)

type repository struct {
	rdb    *redis.Client
	ctx    context.Context
	log    *zap.Logger
	tracer trace.Tracer
}

func New(rdb *redis.Client, ctx context.Context, log *zap.Logger, tracer trace.Tracer) pkgSessions.Repository {
	return &repository{
		rdb:    rdb,
		ctx:    ctx,
		log:    log,
		tracer: tracer,
	}
}

func (repo *repository) Create(ctx context.Context, userID int) (string, error) {
	_, span := repo.tracer.Start(ctx, componentName+" "+"Create")
	defer span.End()

	authToken := strconv.Itoa(userID) + "$" + uuid.New().String()

	err := repo.rdb.HSet(repo.ctx, strconv.Itoa(userID), authToken, []byte{}).Err()
	if err != nil {
		repo.log.Error("Failed to set key-value in Redis", zap.Error(err), zap.Int("user_id", userID),
			zap.String("auth_token", authToken))
		repo.rdb.Expire(repo.ctx, strconv.Itoa(userID), 5*time.Hour)
		return "", err
	}

	return authToken, nil
}

func (repo *repository) Get(ctx context.Context, userID int, authToken string) (int, error) {
	_, span := repo.tracer.Start(ctx, componentName+" "+"Get")
	defer span.End()

	err := repo.rdb.HGet(repo.ctx, strconv.Itoa(userID), authToken).Err()
	if err != nil {
		repo.log.Info("Failed to get session", zap.Error(err), zap.Int("user_id", userID),
			zap.String("token", authToken))
		return 0, pkgErrors.ErrSessionNotFound
	}

	return userID, nil
}

func (repo *repository) Delete(ctx context.Context, userID int, authToken string) error {
	_, span := repo.tracer.Start(ctx, componentName+" "+"Delete")
	defer span.End()

	if err := repo.rdb.HGet(repo.ctx, strconv.Itoa(userID), authToken).Err(); err != nil {
		repo.log.Info("Failed to delete session", zap.Error(err), zap.Int("user_id", userID),
			zap.String("token", authToken))
		return pkgErrors.ErrSessionNotFound
	}

	repo.rdb.HDel(repo.ctx, strconv.Itoa(userID), authToken)
	return nil
}
