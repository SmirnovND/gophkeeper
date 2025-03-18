package main

import (
	"github.com/SmirnovND/gophkeeper/internal/container"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/gophkeeper/internal/router"
	"github.com/SmirnovND/toolbox/pkg/logger"
	"github.com/SmirnovND/toolbox/pkg/middleware"
	"github.com/SmirnovND/toolbox/pkg/migrations"
	"github.com/SmirnovND/toolbox/pkg/rabbitmq"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {
	diContainer := container.NewContainer()

	var cf interfaces.ConfigServer
	diContainer.Invoke(func(c interfaces.ConfigServer) {
		cf = c
	})

	var dbx *sqlx.DB
	diContainer.Invoke(func(db *sqlx.DB) {
		dbx = db
	})

	dbBase := dbx.DB
	migrations.StartMigrations(dbBase)

	// Подключение к RabbitMQ через контейнер
	var rabbitConnection *rabbitmq.RabbitMQConnection
	diContainer.Invoke(func(rc *rabbitmq.RabbitMQConnection) {
		rabbitConnection = rc
	})
	defer rabbitConnection.Close()

	var (
		rabbitProducer *rabbitmq.RabbitMQProducer
		rabbitConsumer *rabbitmq.RabbitMQConsumer
	)

	err := diContainer.Invoke(func(p *rabbitmq.RabbitMQProducer, c *rabbitmq.RabbitMQConsumer) {
		rabbitProducer = p
		rabbitConsumer = c
	})
	if err != nil {
		log.Fatalf("Failed to resolve dependencies: %s", err)
	}

	defer rabbitProducer.Close()
	defer rabbitConsumer.Close()

	return http.ListenAndServe(cf.GetRunAddr(), middleware.ChainMiddleware(
		router.Handler(diContainer),
		logger.WithLogging,
	))
}
