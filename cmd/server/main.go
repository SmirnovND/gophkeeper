package main

import (
	_ "github.com/SmirnovND/gophkeeper/docs"
	"github.com/SmirnovND/gophkeeper/internal/container"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/gophkeeper/internal/router"
	"github.com/SmirnovND/toolbox/pkg/logger"
	"github.com/SmirnovND/toolbox/pkg/middleware"
	"github.com/SmirnovND/toolbox/pkg/migrations"
	"github.com/jmoiron/sqlx"
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

	return http.ListenAndServe(cf.GetRunAddr(), middleware.ChainMiddleware(
		router.Handler(diContainer),
		logger.WithLogging,
	))
}
