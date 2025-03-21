package repo

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"os"
	"path/filepath"
)

// TokenStorage - структура для хранения JWT-токенов по пользователям.
type TokenStorage struct {
}

// AuthData - структура для хранения данных авторизации (токенов).
type AuthData struct {
	Token string `json:"token"`
}

// NewTokenStorage создает новый экземпляр TokenStorage.
func NewTokenStorage() interfaces.TokenStorage {
	return &TokenStorage{}
}

// SaveToken сохраняет токен для конкретного пользователя.
func (s *TokenStorage) SaveToken(name, token string) error {
	// Получаем путь к файлу с конфигурацией
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	fmt.Println(configPath)
	// Загружаем существующие данные, если они есть
	existingData, err := s.loadAllTokens(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Если данных нет, создаем пустую структуру
	if existingData == nil {
		existingData = make(map[string]AuthData)
	}

	// Сохраняем новый токен по имени пользователя
	existingData[name] = AuthData{Token: token}

	// Сохраняем обновленные данные
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем данные в файл
	return json.NewEncoder(file).Encode(existingData)
}

// LoadToken загружает токен для конкретного пользователя по имени.
func (s *TokenStorage) LoadToken(name string) (string, error) {
	// Получаем путь к файлу с конфигурацией
	configPath, err := getConfigPath()
	if err != nil {
		return "", err
	}

	// Загружаем все токены
	existingData, err := s.loadAllTokens(configPath)
	if err != nil {
		return "", err
	}

	// Ищем токен по имени пользователя
	if data, ok := existingData[name]; ok {
		return data.Token, nil
	}

	// Если токен не найден
	return "", fmt.Errorf("token for user %s not found", name)
}

// loadAllTokens загружает все токены из файла.
func (s *TokenStorage) loadAllTokens(configPath string) (map[string]AuthData, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string]AuthData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// getConfigPath возвращает путь к файлу конфигурации.
func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "passcli", "auth.json"), nil
}
