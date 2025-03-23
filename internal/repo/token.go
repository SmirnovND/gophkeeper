package repo

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"os"
	"path/filepath"
)

// TokenStorage - структура для хранения JWT-токена.
type TokenStorage struct {
}

// AuthData - структура для хранения данных авторизации (токена).
type AuthData struct {
	Token string `json:"token"`
}

// NewTokenStorage создает новый экземпляр TokenStorage.
func NewTokenStorage() interfaces.TokenStorage {
	return &TokenStorage{}
}

// SaveToken сохраняет токен.
func (s *TokenStorage) SaveToken(token string) error {
	// Получаем путь к файлу с конфигурацией
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Создаем директорию, если она не существует
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Создаем структуру с токеном
	authData := AuthData{Token: token}

	// Сохраняем данные
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем данные в файл
	return json.NewEncoder(file).Encode(authData)
}

// LoadToken загружает токен.
func (s *TokenStorage) LoadToken() (string, error) {
	// Получаем путь к файлу с конфигурацией
	configPath, err := getConfigPath()
	if err != nil {
		return "", err
	}

	// Открываем файл
	file, err := os.Open(configPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Декодируем данные
	var authData AuthData
	if err := json.NewDecoder(file).Decode(&authData); err != nil {
		return "", err
	}

	// Возвращаем токен
	if authData.Token == "" {
		return "", fmt.Errorf("token not found")
	}
	return authData.Token, nil
}

// getConfigPath возвращает путь к файлу конфигурации.
func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "passcli", "auth.json"), nil
}
