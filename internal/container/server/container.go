package server

import (
	config "github.com/SmirnovND/gophkeeper/internal/config/server"
	"github.com/SmirnovND/gophkeeper/internal/controllers"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/gophkeeper/internal/repo"
	"github.com/SmirnovND/gophkeeper/internal/service"
	"github.com/SmirnovND/gophkeeper/internal/usecase"
	"github.com/SmirnovND/toolbox/pkg/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/dig"
)

// Container - структура контейнера, обертывающая dig-контейнер
type Container struct {
	container *dig.Container
}

func NewContainer() *Container {
	c := &Container{container: dig.New()}
	c.provideDependencies()
	c.provideRepo()
	c.provideService()
	c.provideUsecase()
	c.provideController()
	return c
}

// provideDependencies - функция, регистрирующая зависимости
func (c *Container) provideDependencies() {
	// Регистрируем конфигурацию
	c.container.Provide(config.NewConfig)
	c.container.Provide(func(configServer interfaces.ConfigServer) *sqlx.DB {
		return db.NewDB(configServer.GetDBDsn())
	})
	// Регистрируем DB интерфейс
	// При регистрации:
	c.container.Provide(func(db *sqlx.DB) interfaces.DB {
		return NewDBAdapter(db)
	})

	c.container.Provide(func(configServer interfaces.ConfigServer) *s3.S3 {
		sess := session.Must(session.NewSession(&aws.Config{
			Region:                        aws.String(configServer.GetS3Region()),
			CredentialsChainVerboseErrors: aws.Bool(true),
			Credentials: credentials.NewStaticCredentials(
				configServer.GetS3AccessKey(),
				configServer.GetS3SecretKey(),
				"",
			),
		}))

		return s3.New(sess)
	})

}

type DBAdapter struct {
	*sqlx.DB
}

func NewDBAdapter(db *sqlx.DB) *DBAdapter {
	return &DBAdapter{db}
}

func (d *DBAdapter) QueryRow(query string, args ...interface{}) *sqlx.Row {
	return d.DB.QueryRowx(query, args...) // Используем QueryRowx вместо QueryRow
}

func (c *Container) provideUsecase() {
	c.container.Provide(usecase.NewAuthUseCase)
	c.container.Provide(usecase.NewCloudUseCase)
}

func (c *Container) provideRepo() {
	c.container.Provide(repo.NewUserRepo)
}

func (c *Container) provideService() {
	c.container.Provide(service.NewAuthService)
	c.container.Provide(service.NewUserService)

	c.container.Provide(func(svc *s3.S3, configServer interfaces.ConfigServer) interfaces.AwsService {
		return service.NewAws(svc, configServer.GetS3BucketName())
	})

}

func (c *Container) provideController() {
	c.container.Provide(controllers.NewAuthController)
	c.container.Provide(controllers.NewHealthcheckController)
	c.container.Provide(controllers.NewFileController)
}

// Invoke - функция для вызова и инжекта зависимостей
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}
