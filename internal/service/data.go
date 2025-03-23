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

// SaveFileMetadata сохраняет метаданные файла
func (c *DataService) SaveFileMetadata(login string, label string, fileData *domain.FileData) error {
	// Получаем пользователя по логину
	user, err := c.userRepo.FindUser(login)
	if err != nil {
		return fmt.Errorf("Ошибка при поиске пользователя: %w", err)
	}

	// Создаем метаданные файла (сохраняем только имя и расширение, URL не сохраняем)
	fileMetadata := domain.FileMetadata{
		FileName:  fileData.Name,
		Extension: fileData.Extension,
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

// GetFileMetadata получает метаданные файла
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

// SaveCredential сохраняет учетные данные (логин/пароль)
func (c *DataService) SaveCredential(login string, label string, credentialData *domain.CredentialData) error {
	// Получаем пользователя по логину
	user, err := c.userRepo.FindUser(login)
	if err != nil {
		return fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	// Преобразуем данные в JSON
	dataJSON, err := json.Marshal(credentialData)
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге учетных данных: %w", err)
	}

	// Создаем запись в таблице user_data
	userData := &domain.UserData{
		UserID: user.Id,
		Label:  label,
		Type:   domain.UserDataTypeCredential,
		Data:   dataJSON,
	}

	// Сохраняем запись в базе данных
	err = c.repo.SaveUserData(userData)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении учетных данных: %w", err)
	}

	return nil
}

// GetCredential получает учетные данные (логин/пароль)
func (c *DataService) GetCredential(login string, label string) (*domain.CredentialData, error) {
	// Получаем пользователя по логину
	user, err := c.userRepo.FindUser(login)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	// Получаем данные пользователя по метке и типу
	userData, err := c.repo.GetUserDataByLabelAndType(user.Id, label, domain.UserDataTypeCredential)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении учетных данных: %w", err)
	}

	// Если данные не найдены
	if userData == nil {
		return nil, fmt.Errorf("учетные данные не найдены")
	}

	// Десериализуем данные из JSON
	var credentialData domain.CredentialData
	err = json.Unmarshal(userData.Data, &credentialData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при десериализации учетных данных: %w", err)
	}

	return &credentialData, nil
}

// SaveCard сохраняет данные кредитной карты
func (c *DataService) SaveCard(login string, label string, cardData *domain.CardData) error {
	// Получаем пользователя по логину
	user, err := c.userRepo.FindUser(login)
	if err != nil {
		return fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	// Преобразуем данные в JSON
	dataJSON, err := json.Marshal(cardData)
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге данных карты: %w", err)
	}

	// Создаем запись в таблице user_data
	userData := &domain.UserData{
		UserID: user.Id,
		Label:  label,
		Type:   domain.UserDataTypeCard,
		Data:   dataJSON,
	}

	// Сохраняем запись в базе данных
	err = c.repo.SaveUserData(userData)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении данных карты: %w", err)
	}

	return nil
}

// GetCard получает данные кредитной карты
func (c *DataService) GetCard(login string, label string) (*domain.CardData, error) {
	// Получаем пользователя по логину
	user, err := c.userRepo.FindUser(login)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	// Получаем данные пользователя по метке и типу
	userData, err := c.repo.GetUserDataByLabelAndType(user.Id, label, domain.UserDataTypeCard)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных карты: %w", err)
	}

	// Если данные не найдены
	if userData == nil {
		return nil, fmt.Errorf("данные карты не найдены")
	}

	// Десериализуем данные из JSON
	var cardData domain.CardData
	err = json.Unmarshal(userData.Data, &cardData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при десериализации данных карты: %w", err)
	}

	return &cardData, nil
}

// SaveText сохраняет текстовые данные
func (c *DataService) SaveText(login string, label string, textData *domain.TextData) error {
	// Получаем пользователя по логину
	user, err := c.userRepo.FindUser(login)
	if err != nil {
		return fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	// Преобразуем данные в JSON
	dataJSON, err := json.Marshal(textData)
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге текстовых данных: %w", err)
	}

	// Создаем запись в таблице user_data
	userData := &domain.UserData{
		UserID: user.Id,
		Label:  label,
		Type:   domain.UserDataTypeText,
		Data:   dataJSON,
	}

	// Сохраняем запись в базе данных
	err = c.repo.SaveUserData(userData)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении текстовых данных: %w", err)
	}

	return nil
}

// GetText получает текстовые данные
func (c *DataService) GetText(login string, label string) (*domain.TextData, error) {
	// Получаем пользователя по логину
	user, err := c.userRepo.FindUser(login)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	// Получаем данные пользователя по метке и типу
	userData, err := c.repo.GetUserDataByLabelAndType(user.Id, label, domain.UserDataTypeText)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении текстовых данных: %w", err)
	}

	// Если данные не найдены
	if userData == nil {
		return nil, fmt.Errorf("текстовые данные не найдены")
	}

	// Десериализуем данные из JSON
	var textData domain.TextData
	err = json.Unmarshal(userData.Data, &textData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при десериализации текстовых данных: %w", err)
	}

	return &textData, nil
}
