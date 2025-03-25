package pkg

import (
	"os"
	"path/filepath"
	"runtime"
)

// UserHomeDirFunc определяет тип функции для получения домашней директории пользователя
type UserHomeDirFunc func() (string, error)

// Переменная, содержащая функцию для получения домашней директории
// По умолчанию использует os.UserHomeDir, но может быть переопределена для тестов
var userHomeDirFunc UserHomeDirFunc = os.UserHomeDir

// GetOSFunc определяет тип функции для получения текущей операционной системы
type GetOSFunc func() string

// Переменная, содержащая функцию для получения текущей операционной системы
// По умолчанию возвращает runtime.GOOS, но может быть переопределена для тестов
var getOSFunc GetOSFunc = func() string {
	return runtime.GOOS
}

// GetDownloadsDir возвращает путь к директории загрузок пользователя
func GetDownloadsDir() string {
	// Получаем домашнюю директорию пользователя
	homeDir, err := userHomeDirFunc()
	if err != nil {
		// Если не удалось получить домашнюю директорию, используем текущую
		return "."
	}

	// В зависимости от операционной системы, выбираем путь к директории загрузок
	switch getOSFunc() {
	case "windows":
		return filepath.Join(homeDir, "Downloads")
	case "darwin": // macOS
		return filepath.Join(homeDir, "Downloads")
	default: // Linux и другие Unix-подобные системы
		// Проверяем существование директории XDG_DOWNLOAD_DIR
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			xdgConfigHome = filepath.Join(homeDir, ".config")
		}

		// Проверяем наличие файла user-dirs.dirs
		userDirsFile := filepath.Join(xdgConfigHome, "user-dirs.dirs")
		if _, err := os.Stat(userDirsFile); err == nil {
			// Если файл существует, можно было бы прочитать его и найти XDG_DOWNLOAD_DIR
			// Но для простоты используем стандартный путь
			return filepath.Join(homeDir, "Downloads")
		}

		// Если не нашли конфигурацию, используем стандартный путь
		return filepath.Join(homeDir, "Downloads")
	}
}
