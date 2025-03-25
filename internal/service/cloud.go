package service

import (
	"context"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"net/url"
	"time"
)

type Cloud struct {
	minio      interfaces.MinioClientInterface
	bucketName string
}

func NewCloud(minio interfaces.MinioClientInterface, bucketName string) interfaces.CloudService {
	return &Cloud{
		minio:      minio,
		bucketName: bucketName,
	}
}

func (c *Cloud) GenerateUploadLink(fileName string) (string, error) {
	ctx := context.Background()
	presignedURL, err := c.minio.PresignedPutObject(ctx, c.bucketName, fileName, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (c *Cloud) GenerateDownloadLink(fileName string) (string, error) {
	ctx := context.Background()
	// Устанавливаем срок действия ссылки на 15 минут
	reqParams := make(url.Values)
	presignedURL, err := c.minio.PresignedGetObject(ctx, c.bucketName, fileName, 15*time.Minute, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
