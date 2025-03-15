package interfaces

type Config interface {
	GetJwtSecret() string
	GetDBDsn() string
}
