package client

import (
	"github.com/SmirnovND/gophkeeper/internal/command"
	config "github.com/SmirnovND/gophkeeper/internal/config/client"
	"github.com/SmirnovND/gophkeeper/internal/repo"
	"github.com/SmirnovND/gophkeeper/internal/service"
	"github.com/SmirnovND/gophkeeper/internal/usecase"
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
	c.provideCommand()
	return c
}

// provideDependencies - функция, регистрирующая зависимости
func (c *Container) provideDependencies() {
	// Регистрируем конфигурацию
	c.container.Provide(config.NewConfig)
}

func (c *Container) provideUsecase() {
	c.container.Provide(usecase.NewClientUseCase)
}

func (c *Container) provideRepo() {
	c.container.Provide(repo.NewTokenStorage)
}

func (c *Container) provideService() {
	c.container.Provide(service.NewTokenService)
	c.container.Provide(service.NewClientService)
}

func (c *Container) provideCommand() {
	c.container.Provide(command.NewCommand)
}

// Invoke - функция для вызова и инжекта зависимостей
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}
