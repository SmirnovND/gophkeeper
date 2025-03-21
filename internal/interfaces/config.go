package interfaces

type ConfigServer interface {
	GetJwtSecret() string
	GetDBDsn() string
	GetRunAddr() string
	GetS3BucketName() string
	GetS3Region() string
	GetS3AccessKey() string
	GetS3SecretKey() string
}
