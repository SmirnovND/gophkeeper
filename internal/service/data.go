package service

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
)

// DataService реализует интерфейс для работы с данными пользователя
type DataService struct {
	repo     interfaces.UserDataRepo
	userRepo interfaces.UserRepo
}

// NewDataService создает новый экземпляр DataService
func NewDataService(repo interfaces.UserDataRepo, userRepo interfaces.UserRepo) interfaces.DataService {
	return &DataService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (c *DataService) SaveFileMetadata(login string, label string, fileData *domain.FileData, uploadLink string) error {
	// Получаем пользователя по логину
	user, err := c.userRepo.FindUser(login)
	if err != nil {
		return fmt.Errorf("Ошибка при поиске пользователя: %w", err)
	}

	// Создаем метаданные файла
	fileMetadata := domain.FileMetadata{
		FileName:  fileData.Name,
		Extension: fileData.Extension,
		URL:       uploadLink,
	}

	// Преобразуем метаданные в JSON
	metadataJSON, err := json.Marshal(fileMetadata)
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге метаданных файла: %w", err)
	}

	// Создаем запись в таблице user_data
	userData := &domain.UserData{
		UserID: user.Id,
		Label:  label,
		Type:   domain.UserDataTypeFile,
		Data:   metadataJSON,
	}

	// Сохраняем запись в базе данных
	err = c.repo.SaveUserData(userData)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении метаданных файла: %w", err)
	}

	return nil
}

func (c *DataService) GetFileMetadata(login string, label string) (*domain.FileMetadata, error) {
	// Получаем пользователя по логину
	user, err := c.userRepo.FindUser(login)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	// Получаем данные пользователя по метке и типу
	userData, err := c.repo.GetUserDataByLabelAndType(user.Id, label, domain.UserDataTypeFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении метаданных файла: %w", err)
	}

	// Если данные не найдены
	if userData == nil {
		return nil, fmt.Errorf("метаданные файла не найдены")
	}

	// Десериализуем метаданные из JSON
	var fileMetadata domain.FileMetadata
	err = json.Unmarshal(userData.Data, &fileMetadata)
	if err != nil {
		return nil, fmt.Errorf("ошибка при десериализации метаданных файла: %w", err)
	}

	return &fileMetadata, nil
}
