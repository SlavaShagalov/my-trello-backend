package storages

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
)

func NewS3(log *zap.Logger) (*s3.Client, error) {
	log.Info("Connecting to S3...")

	cfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithDefaultRegion(viper.GetString(config.S3DefaultRegion)),
	)
	if err != nil {
		log.Error("Failed to create S3 connection", zap.Error(err))
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.EndpointResolver = s3.EndpointResolverFromURL(viper.GetString(config.S3Endpoint))
	})

	log.Info("S3 bucket connection created successfully")
	return client, nil
}
