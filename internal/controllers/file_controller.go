package controllers

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/toolbox/pkg/paramsparser"
	"net/http"
)

type FileController struct {
	FileUseCase interfaces.CloudUseCase
}

func NewFileController(FileUseCase interfaces.CloudUseCase) *FileController {
	return &FileController{
		FileUseCase: FileUseCase,
	}
}

// HandleUploadFile godoc
// @Summary Загрузка файла
// @Description Генерирует ссылку для загрузки файла на сервер
// @Tags files
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен авторизации"
// @Param fileData body domain.FileData true "Информация о загружаемом файле"
// @Success 200 {object} map[string]string "Успешная генерация ссылки, возвращает URL для загрузки файла"
// @Failure 400 {object} map[string]string "Ошибка в формате запроса"
// @Failure 401 {object} map[string]string "Пользователь не авторизован"
// @Failure 413 {object} map[string]string "Превышен максимальный размер файла"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/files/upload [post]
func (f *FileController) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	fileData, err := paramsparser.JSONParse[domain.FileData](w, r)
	if err != nil {
		return
	}

	f.FileUseCase.GenerateUploadLink(w, r, fileData)
}

// HandleDownloadFile godoc
// @Summary Скачивание файла
// @Description Генерирует ссылку для скачивания файла с сервера
// @Tags files
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен авторизации"
// @Param label query string true "Метка файла"
// @Success 200 {object} map[string]string "Успешная генерация ссылки, возвращает URL для скачивания файла"
// @Failure 400 {object} map[string]string "Ошибка в формате запроса"
// @Failure 401 {object} map[string]string "Пользователь не авторизован"
// @Failure 404 {object} map[string]string "Файл не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/files/download [get]
func (f *FileController) HandleDownloadFile(w http.ResponseWriter, r *http.Request) {
	// Получаем метку файла из параметров запроса
	label := r.URL.Query().Get("label")
	if label == "" {
		http.Error(w, "Не указана метка файла", http.StatusBadRequest)
		return
	}

	// Генерируем ссылку для скачивания
	f.FileUseCase.GenerateDownloadLink(w, r, label)
}
