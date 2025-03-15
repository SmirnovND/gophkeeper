package container

import (
	config "github.com/SmirnovND/gophkeeper/internal/config/server"
	"github.com/SmirnovND/toolbox/pkg/db"
	"github.com/SmirnovND/toolbox/pkg/http"
	"github.com/SmirnovND/toolbox/pkg/rabbitmq"
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
	c.provideRabbitMQ()
	return c
}

// provideDependencies - функция, регистрирующая зависимости
func (c *Container) provideDependencies() {
	// Регистрируем конфигурацию
	c.container.Provide(config.NewConfig())
	c.container.Provide(func(cfg *config.Config) *sqlx.DB {
		return db.NewDB(cfg.GetDBDsn())
	})
	c.container.Provide(db.NewTransactionManager)
	c.container.Provide(http.NewAPIClient)
}

func (c *Container) provideUsecase() {
	c.container.Provide(usecase.NewAuthUseCase)
	c.container.Provide(usecase.NewOrderUseCase)
	c.container.Provide(usecase.NewUserUseCase)
	c.container.Provide(usecase.NewProcessingUseCase)
}

func (c *Container) provideRepo() {
	c.container.Provide(repo.NewUserRepo)
	c.container.Provide(repo.NewOrderRepo)
	c.container.Provide(repo.NewBalanceRepo)
	c.container.Provide(repo.NewTransactionRepo)
}

func (c *Container) provideService() {
	c.container.Provide(service.NewAuthService)
	c.container.Provide(service.NewUserService)
	c.container.Provide(service.NewOrderService)
	c.container.Provide(service.NewBalanceService)
	c.container.Provide(service.NewRabbitMqService)
	c.container.Provide(func(cfg *config.Config, client *http.APIClient) *service.ProcessingService {
		return service.NewProcessingService(cfg.AccrualSystemAddress, client)
	})
}

func (c *Container) provideController() {
	c.container.Provide(controllers.NewAuthController)
	c.container.Provide(controllers.NewOrderController)
	c.container.Provide(controllers.NewUserController)
}

func (c *Container) provideRabbitMQ() {
	c.container.Provide(func(cfg *config.Config) *rabbitmq.RabbitMQConnection {
		return rabbitmq.NewRabbitMQConnection(cfg.GetRabbitURL())
	})
	c.container.Provide(func(conn *rabbitmq.RabbitMQConnection, cfg *config.Config) *rabbitmq.RabbitMQProducer {
		return rabbitmq.NewRabbitMQProducer(conn.Conn)
	})
	c.container.Provide(func(conn *rabbitmq.RabbitMQConnection, cfg *config.Config) *rabbitmq.RabbitMQConsumer {
		return rabbitmq.NewRabbitMQConsumer(conn.Conn, cfg.GetRabbitQueue())
	})
}

// Invoke - функция для вызова и инжекта зависимостей
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}

func ProvideDBDsn() string {
	return "postgresql://developer:developer@localhost:5432/postgres?sslmode=disable"
}
