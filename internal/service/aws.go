package service

import (
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/gophkeeper/pkg"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Aws struct {
	svc        *s3.S3
	bucketName string
}

func NewAws(svc *s3.S3, bucketName string) interfaces.AwsService {
	return &Aws{
		svc:        svc,
		bucketName: bucketName,
	}
}

func (awc *Aws) GenerateUploadLink(fileName string) (string, error) {
	return pkg.GeneratePreSignedURL(awc.svc, fileName, awc.bucketName)
}
