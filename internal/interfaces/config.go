package interfaces

type ConfigServer interface {
	GetJwtSecret() string
	GetDBDsn() string
	GetRabbitMQURI() string
	GetRunAddr() string
}
