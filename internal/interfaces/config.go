package interfaces

type ConfigServer interface {
	GetJwtSecret() string
	GetDBDsn() string
	GetRunAddr() string
	GetMinioBucketName() string
	GetMinioAccessKey() string
	GetMinioSecretKey() string
	GetMinioHost() string
}
