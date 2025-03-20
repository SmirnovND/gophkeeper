package interfaces

type ConfigServer interface {
	GetJwtSecret() string
	GetDBDsn() string
	GetRunAddr() string
}

type ConfigClient interface {
	GetServerAddr() string
}
