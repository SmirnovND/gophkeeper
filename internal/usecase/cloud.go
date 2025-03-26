package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"net/http"
)

type CloudUseCase struct {
	cloudService interfaces.CloudService
	dataService  interfaces.DataService
	jwtService   interfaces.JwtService
}

func NewCloudUseCase(
	cloudService interfaces.CloudService,
	dataService interfaces.DataService,
	jwtService interfaces.JwtService,
) interfaces.CloudUseCase {
	return &CloudUseCase{
		cloudService: cloudService,
		dataService:  dataService,
		jwtService:   jwtService,
	}
}

func (c *CloudUseCase) GenerateUploadLink(w http.ResponseWriter, r *http.Request, fileData *domain.FileData) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Валидация входящего объекта FileData
	if fileData == nil || fileData.Name == "" || fileData.Extension == "" {
		http.Error(w, "Неверные данные файла", http.StatusBadRequest)
		return
	}

	// Формируем имя файла
	fileName := fmt.Sprintf("%s_%s.%s", login, fileData.Name, fileData.Extension)

	// Получаем ссылку для загрузки
	uploadLink, err := c.cloudService.GenerateUploadLink(fileName)

	if err != nil {
		http.Error(w, "Ошибка при генерации ссылки: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем метаданные файла в таблице user_data
	err = c.dataService.SaveFileMetadata(login, fileData.Name, fileData, fileData.Metadata)
	if err != nil {
		http.Error(w, "Ошибка при сохранении метаданных файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := domain.FileDataResponse{
		Url:         uploadLink,
		Description: "Загрузи файл по этой ссылке",
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *CloudUseCase) GenerateDownloadLink(w http.ResponseWriter, r *http.Request, label string) {
	// Получаем логин пользователя из токена
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Валидация входящих данных
	if label == "" {
		http.Error(w, "Не указана метка файла", http.StatusBadRequest)
		return
	}

	// Получаем метаданные файла из базы данных
	fileMetadata, metadata, err := c.dataService.GetFileMetadata(login, label)
	if err != nil {
		http.Error(w, "Ошибка при получении метаданных файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Формируем имя файла
	fileName := fmt.Sprintf("%s_%s.%s", login, fileMetadata.FileName, fileMetadata.Extension)

	// Получаем ссылку для скачивания
	downloadLink, err := c.cloudService.GenerateDownloadLink(fileName)
	if err != nil {
		http.Error(w, "Ошибка при генерации ссылки для скачивания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем расширенный ответ с метаданными и метаинформацией
	response := struct {
		URL         string              `json:"url"`
		Description string              `json:"description"`
		Metadata    domain.FileMetadata `json:"metadata"`
		MetaInfo    string              `json:"meta_info"`
	}{
		URL:         downloadLink,
		Description: "Скачай файл по этой ссылке",
		Metadata:    *fileMetadata,
		MetaInfo:    metadata,
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
