package pkg

import (
	"context"
	"github.com/minio/minio-go/v7"
	"time"
)

// GeneratePreSignedURL Генерация Pre-Signed URL
func GeneratePreSignedURL(client *minio.Client, bucketName, fileName string) (string, error) {
	ctx := context.Background()
	presignedURL, err := client.PresignedPutObject(ctx, bucketName, fileName, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
