package server

import (
	"fmt"
	config "github.com/SmirnovND/gophkeeper/internal/config/server"
	"github.com/SmirnovND/gophkeeper/internal/controllers"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/gophkeeper/internal/repo"
	"github.com/SmirnovND/gophkeeper/internal/service"
	"github.com/SmirnovND/gophkeeper/internal/usecase"
	"github.com/SmirnovND/toolbox/pkg/db"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	c.container.Provide(func(configServer interfaces.ConfigServer) *minio.Client {
		client, err := minio.New(configServer.GetMinioHost(), &minio.Options{
			Creds:  credentials.NewStaticV4(configServer.GetMinioAccessKey(), configServer.GetMinioSecretKey(), ""),
			Secure: false, // Без HTTPS для локальной установки
		})
		if err != nil {
			// Исправление: добавляем сообщение об ошибке в вызов panic
			panic(fmt.Sprintf("Ошибка создания MinIO клиента: %v", err))
		}

		return client
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
	c.container.Provide(usecase.NewDataUseCase)
}

func (c *Container) provideRepo() {
	c.container.Provide(repo.NewUserRepo)
	c.container.Provide(repo.NewUserDataRepo)
}

func (c *Container) provideService() {
	c.container.Provide(service.NewAuthService)
	c.container.Provide(service.NewUserService)
	c.container.Provide(service.NewDataService)

	c.container.Provide(func(minio *minio.Client, configServer interfaces.ConfigServer) interfaces.CloudService {
		return service.NewCloud(minio, configServer.GetMinioBucketName())
	})

}

func (c *Container) provideController() {
	c.container.Provide(controllers.NewAuthController)
	c.container.Provide(controllers.NewHealthcheckController)
	c.container.Provide(controllers.NewFileController)
	c.container.Provide(controllers.NewDataController)
}

// Invoke - функция для вызова и инжекта зависимостей
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}
