package pkg

import (
	"context"
	"github.com/minio/minio-go/v7"
	"net/url"
	"time"
)

// GeneratePreSignedURL Генерация Pre-Signed URL для загрузки файла
func GeneratePreSignedURL(client *minio.Client, bucketName, fileName string) (string, error) {
	ctx := context.Background()
	presignedURL, err := client.PresignedPutObject(ctx, bucketName, fileName, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// GenerateDownloadURL Генерация Pre-Signed URL для скачивания файла
func GenerateDownloadURL(client *minio.Client, bucketName, fileName string) (string, error) {
	ctx := context.Background()
	// Устанавливаем срок действия ссылки на 15 минут
	reqParams := make(url.Values)
	presignedURL, err := client.PresignedGetObject(ctx, bucketName, fileName, 15*time.Minute, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
