package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"net/http"
)

type FileUseCase struct {
	awsService interfaces.CloudService
}

func NewCloudUseCase(
	awsService interfaces.CloudService,
) interfaces.FileUseCase {
	return &FileUseCase{
		awsService: awsService,
	}
}

func (c *FileUseCase) UploadFile(w http.ResponseWriter, fileData *domain.FileData) {
	// Валидация входящего объекта FileData
	if fileData == nil || fileData.Name == "" || fileData.Extension == "" {
		http.Error(w, "Неверные данные файла", http.StatusBadRequest)
		return
	}

	// Формируем имя файла
	fileName := fmt.Sprintf("%s.%s", fileData.Name, fileData.Extension)

	// Получаем ссылку для загрузки
	uploadLink, err := c.awsService.GenerateUploadLink(fileName)
	if err != nil {
		http.Error(w, "Ошибка при генерации ссылки: "+err.Error(), http.StatusInternalServerError)
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
