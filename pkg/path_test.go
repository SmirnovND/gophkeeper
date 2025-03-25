package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDownloadsDir(t *testing.T) {
	// Сохраняем оригинальные функции и переменные окружения
	originalUserHomeDirFunc := userHomeDirFunc
	originalGetOSFunc := getOSFunc
	originalXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")

	// Восстанавливаем оригинальные значения после завершения теста
	defer func() {
		userHomeDirFunc = originalUserHomeDirFunc
		getOSFunc = originalGetOSFunc
		if originalXDGConfigHome != "" {
			os.Setenv("XDG_CONFIG_HOME", originalXDGConfigHome)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	tests := []struct {
		name           string
		mockHomeDir    string
		mockGOOS       string
		mockXDGConfig  string
		expectedSuffix string
	}{
		{
			name:           "Windows",
			mockHomeDir:    "C:\\Users\\testuser",
			mockGOOS:       "windows",
			mockXDGConfig:  "",
			expectedSuffix: filepath.Join("Users", "testuser", "Downloads"),
		},
		{
			name:           "macOS",
			mockHomeDir:    "/Users/testuser",
			mockGOOS:       "darwin",
			mockXDGConfig:  "",
			expectedSuffix: filepath.Join("Users", "testuser", "Downloads"),
		},
		{
			name:           "Linux без XDG_CONFIG_HOME",
			mockHomeDir:    "/home/testuser",
			mockGOOS:       "linux",
			mockXDGConfig:  "",
			expectedSuffix: filepath.Join("home", "testuser", "Downloads"),
		},
		{
			name:           "Linux с XDG_CONFIG_HOME",
			mockHomeDir:    "/home/testuser",
			mockGOOS:       "linux",
			mockXDGConfig:  "/home/testuser/.config",
			expectedSuffix: filepath.Join("home", "testuser", "Downloads"),
		},
		{
			name:           "Ошибка получения домашней директории",
			mockHomeDir:    "",
			mockGOOS:       "linux",
			mockXDGConfig:  "",
			expectedSuffix: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем моки
			// Мокаем функцию получения ОС
			getOSFunc = func() string {
				return tt.mockGOOS
			}

			// Устанавливаем переменную окружения XDG_CONFIG_HOME
			if tt.mockXDGConfig != "" {
				os.Setenv("XDG_CONFIG_HOME", tt.mockXDGConfig)
			} else {
				os.Unsetenv("XDG_CONFIG_HOME")
			}

			// Мокаем функцию получения домашней директории
			if tt.mockHomeDir == "" {
				userHomeDirFunc = func() (string, error) {
					return "", os.ErrNotExist
				}
			} else {
				userHomeDirFunc = func() (string, error) {
					return tt.mockHomeDir, nil
				}
			}

			// Вызываем тестируемую функцию
			result := GetDownloadsDir()

			// Проверяем результат
			if tt.mockHomeDir == "" {
				assert.Equal(t, ".", result)
			} else if tt.mockGOOS == "windows" {
				// Для Windows используем специальную проверку из-за разных разделителей путей
				assert.Equal(t, filepath.Join(tt.mockHomeDir, "Downloads"), result)
			} else {
				assert.Contains(t, result, tt.expectedSuffix)
			}
		})
	}
}