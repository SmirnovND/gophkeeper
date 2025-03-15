package interfaces

type ConfigServer interface {
	GetJwtSecret() string
	GetDBDsn() string
}
