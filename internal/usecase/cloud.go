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
}

func NewCloudUseCase(
	cloudService interfaces.CloudService,
	dataService interfaces.DataService,
) interfaces.CloudUseCase {
	return &CloudUseCase{
		cloudService: cloudService,
		dataService:  dataService,
	}
}

func (c *CloudUseCase) GenerateUploadLink(w http.ResponseWriter, fileData *domain.FileData, login string) {
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
	err = c.dataService.SaveFileMetadata(login, fileData.Name, fileData, uploadLink)
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
