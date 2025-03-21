package service

import (
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/gophkeeper/pkg"
	"github.com/minio/minio-go/v7"
)

type Cloud struct {
	minio      *minio.Client
	bucketName string
}

func NewCloud(minio *minio.Client, bucketName string) interfaces.CloudService {
	return &Cloud{
		minio:      minio,
		bucketName: bucketName,
	}
}

func (awc *Cloud) GenerateUploadLink(fileName string) (string, error) {
	return pkg.GeneratePreSignedURL(awc.minio, awc.bucketName, fileName)
}
