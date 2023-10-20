package s3

import (
	"bytes"
	"context"
	pImages "github.com/SlavaShagalov/my-trello-backend/internal/images"
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"go.uber.org/zap"
)

type repository struct {
	uploader *manager.Uploader
	log      *zap.Logger
}

func New(s3Client *s3.Client, log *zap.Logger) pImages.Repository {
	return &repository{
		uploader: manager.NewUploader(s3Client),
		log:      log,
	}
}

func (repo *repository) Create(imgName string, imgData []byte) (location string, err error) {
	repo.log.Debug("Start image creating...")

	bucketName := viper.GetString(config.S3BucketName)
	output, err := repo.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &imgName,
		Body:   bytes.NewReader(imgData),
	})
	if err != nil {
		repo.log.Error("Failed to create image", zap.Error(err))
		return "", err
	}

	repo.log.Debug("Image created", zap.String("location", output.Location))
	return output.Location, nil
}

func (repo *repository) Get(location string) (imgData []byte, err error) {
	return nil, nil
}

func (repo *repository) Update(location string, imgData []byte) (err error) {
	repo.log.Debug("Start image updating...")

	bucketName := viper.GetString(config.S3BucketName)

	prefixS := "https://" + bucketName + ".hb.vkcs.cloud/"
	prefix := "http://" + bucketName + ".hb.vkcs.cloud/"
	imgName := strings.TrimPrefix(location, prefixS)
	imgName = strings.TrimPrefix(imgName, prefix)

	output, err := repo.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &imgName,
		Body:   bytes.NewReader(imgData),
	})
	if err != nil {
		repo.log.Error("Failed to update image", zap.Error(err))
		return err
	}

	repo.log.Debug("Image updated", zap.String("location", output.Location))
	return nil
}

func (repo *repository) Delete(location string) (err error) {
	return nil
}
