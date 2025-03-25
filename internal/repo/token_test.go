package repo

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// Переопределяем функцию getConfigPath для тестов
func mockGetConfigPath() func() (string, error) {

	// Создаем временную директорию для тестов
	tempDir, _ := os.MkdirTemp("", "token_test")

	// Возвращаем функцию, которая будет использоваться вместо оригинальной
	return func() (string, error) {
		return filepath.Join(tempDir, "auth.json"), nil
	}
}

// Восстанавливаем оригинальную функцию после тестов
func restoreGetConfigPath(original func() (string, error)) {
	getConfigPath = original
}

// Тест для SaveToken и LoadToken
func TestTokenStorage_SaveAndLoadToken(t *testing.T) {
	// Сохраняем оригинальную функцию и подменяем ее на тестовую
	original := getConfigPath
	getConfigPath = mockGetConfigPath()
	defer restoreGetConfigPath(original)

	// Создаем экземпляр TokenStorage
	storage := NewTokenStorage()

	// Тестовый токен
	testToken := "test-jwt-token"

	// Сохраняем токен
	err := storage.SaveToken(testToken)
	assert.NoError(t, err)

	// Загружаем токен
	loadedToken, err := storage.LoadToken()
	assert.NoError(t, err)
	assert.Equal(t, testToken, loadedToken)
}

// Тест для LoadToken - файл не существует
func TestTokenStorage_LoadToken_FileNotExists(t *testing.T) {
	// Сохраняем оригинальную функцию и подменяем ее на тестовую
	original := getConfigPath
	getConfigPath = mockGetConfigPath()
	defer restoreGetConfigPath(original)

	// Создаем экземпляр TokenStorage
	storage := NewTokenStorage()

	// Пытаемся загрузить токен из несуществующего файла
	_, err := storage.LoadToken()
	assert.Error(t, err)
}

// Тест для LoadToken - некорректный формат файла
func TestTokenStorage_LoadToken_InvalidFormat(t *testing.T) {
	// Сохраняем оригинальную функцию и подменяем ее на тестовую
	original := getConfigPath
	getConfigPath = mockGetConfigPath()
	defer func() {
		restoreGetConfigPath(original)
	}()

	// Получаем путь к файлу
	configPath, _ := getConfigPath()

	// Создаем директорию, если она не существует
	dir := filepath.Dir(configPath)
	_ = os.MkdirAll(dir, 0755)

	// Создаем файл с некорректным содержимым
	file, _ := os.Create(configPath)
	file.WriteString("invalid json")
	file.Close()

	// Создаем экземпляр TokenStorage
	storage := NewTokenStorage()

	// Пытаемся загрузить токен из файла с некорректным форматом
	_, err := storage.LoadToken()
	assert.Error(t, err)
}

// Тест для LoadToken - пустой токен
func TestTokenStorage_LoadToken_EmptyToken(t *testing.T) {
	// Сохраняем оригинальную функцию и подменяем ее на тестовую
	original := getConfigPath
	getConfigPath = mockGetConfigPath()
	defer func() {
		restoreGetConfigPath(original)
	}()

	// Получаем путь к файлу
	configPath, _ := getConfigPath()

	// Создаем директорию, если она не существует
	dir := filepath.Dir(configPath)
	_ = os.MkdirAll(dir, 0755)

	// Создаем файл с пустым токеном
	file, _ := os.Create(configPath)
	json.NewEncoder(file).Encode(AuthData{Token: ""})
	file.Close()

	// Создаем экземпляр TokenStorage
	storage := NewTokenStorage()

	// Пытаемся загрузить пустой токен
	_, err := storage.LoadToken()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token not found")
}

// Тест для SaveToken - ошибка создания директории
func TestTokenStorage_SaveToken_DirectoryError(t *testing.T) {
	// Подменяем функцию getConfigPath на функцию, которая возвращает путь к несуществующей директории
	original := getConfigPath
	getConfigPath = func() (string, error) {
		// Используем путь, который не может быть создан (например, корневой каталог устройства)
		return "/dev/null/impossible/path", nil
	}
	defer restoreGetConfigPath(original)

	// Создаем экземпляр TokenStorage
	storage := NewTokenStorage()

	// Пытаемся сохранить токен
	err := storage.SaveToken("test-token")
	assert.Error(t, err)
}
