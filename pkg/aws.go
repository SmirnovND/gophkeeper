package pkg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

// GeneratePreSignedURL Генерация Pre-Signed URL
func GeneratePreSignedURL(svc *s3.S3, fileName string, bucketName string) (string, error) {
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	// Даем ссылке срок жизни 15 минут
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации Pre-Signed URL: %v", err)
	}

	return urlStr, nil
}
