package router

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/container"
	"github.com/SmirnovND/gophkeeper/internal/controllers"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func Handler(diContainer *container.Container) http.Handler {
	var HealthcheckController *controllers.HealthcheckController
	var AuthController *controllers.AuthController
	var cf interfaces.ConfigServer
	err := diContainer.Invoke(func(
		c interfaces.ConfigServer,
		healthcheckControl *controllers.HealthcheckController,
		authControl *controllers.AuthController,
	) {
		HealthcheckController = healthcheckControl
		AuthController = authControl
		cf = c
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", cf.GetRunAddr())),
	))

	r.Post("/api/user/register", AuthController.HandleRegisterJSON)
	r.Post("/api/user/login", AuthController.HandleLoginJSON)

	r.Get("/ping", HealthcheckController.HandlePing)

	// Обработчик для неподходящего метода (405 Method Not Allowed)
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Обработчик для несуществующих маршрутов (404 Not Found)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Route not found", http.StatusNotFound)
	})

	return r
}
