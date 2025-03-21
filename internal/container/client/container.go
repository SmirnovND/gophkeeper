package client

import (
	"github.com/SmirnovND/gophkeeper/internal/command"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
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

func NewContainer(serverAddress string) *Container {
	c := &Container{container: dig.New()}
	c.provideDependencies()
	c.provideRepo()
	c.provideService(serverAddress)
	c.provideUsecase()
	c.provideCommand()
	return c
}

// provideDependencies - функция, регистрирующая зависимости
func (c *Container) provideDependencies() {
}

func (c *Container) provideUsecase() {
	c.container.Provide(usecase.NewClientUseCase)
}

func (c *Container) provideRepo() {
	c.container.Provide(repo.NewTokenStorage)
}

func (c *Container) provideService(serverAddress string) {
	c.container.Provide(service.NewTokenService)

	c.container.Provide(func() interfaces.ClientService {
		return service.NewClientService(serverAddress)
	})
}

func (c *Container) provideCommand() {
	c.container.Provide(command.NewCommand)
}

// Invoke - функция для вызова и инжекта зависимостей
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}
