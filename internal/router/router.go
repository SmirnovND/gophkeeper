package router

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/container/server"
	"github.com/SmirnovND/gophkeeper/internal/controllers"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/toolbox/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func Handler(diContainer *server.Container) http.Handler {
	var HealthcheckController *controllers.HealthcheckController
	var AuthController *controllers.AuthController
	var FileController *controllers.FileController
	var DataController *controllers.DataController
	var cf interfaces.ConfigServer
	err := diContainer.Invoke(func(
		c interfaces.ConfigServer,
		healthcheckControl *controllers.HealthcheckController,
		authControl *controllers.AuthController,
		fileControl *controllers.FileController,
		dataControl *controllers.DataController,
	) {
		HealthcheckController = healthcheckControl
		AuthController = authControl
		FileController = fileControl
		DataController = dataControl
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

	r.Post("/api/file/upload", func(w http.ResponseWriter, r *http.Request) {
		auth.AuthMiddleware(cf.GetJwtSecret(), http.HandlerFunc(FileController.HandleUploadFile)).ServeHTTP(w, r)
	})

	r.Get("/api/file/download", func(w http.ResponseWriter, r *http.Request) {
		auth.AuthMiddleware(cf.GetJwtSecret(), http.HandlerFunc(FileController.HandleDownloadFile)).ServeHTTP(w, r)
	})

	// Маршруты для работы с данными пользователя
	r.Route("/api/data", func(r chi.Router) {
		// Применяем middleware аутентификации ко всем маршрутам данных
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				auth.AuthMiddleware(cf.GetJwtSecret(), next).ServeHTTP(w, r)
			})
		})

		// Маршруты для работы с учетными данными (логин/пароль)
		r.Route("/credential/{label}", func(r chi.Router) {
			r.Post("/", DataController.SaveCredential)
			r.Get("/", DataController.GetCredential)
		})

		// Маршруты для работы с данными кредитных карт
		r.Route("/card/{label}", func(r chi.Router) {
			r.Post("/", DataController.SaveCard)
			r.Get("/", DataController.GetCard)
		})

		// Маршруты для работы с текстовыми данными
		r.Route("/text/{label}", func(r chi.Router) {
			r.Post("/", DataController.SaveText)
			r.Get("/", DataController.GetText)
		})
	})

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
