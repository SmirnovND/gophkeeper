package container

import (
	config "github.com/SmirnovND/gophkeeper/internal/config/server"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/toolbox/pkg/db"
	"github.com/SmirnovND/toolbox/pkg/http"
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
	c.container.Provide(config.NewConfig())
	c.container.Provide(func(configServer interfaces.ConfigServer) *sqlx.DB {
		return db.NewDB(configServer.GetDBDsn())
	})
	c.container.Provide(db.NewTransactionManager)
	c.container.Provide(http.NewAPIClient)
}

func (c *Container) provideUsecase() {
}

func (c *Container) provideRepo() {
}

func (c *Container) provideService() {
}

func (c *Container) provideController() {
}

// Invoke - функция для вызова и инжекта зависимостей
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}
