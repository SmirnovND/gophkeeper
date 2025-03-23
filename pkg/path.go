package pkg

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDownloadsDir возвращает путь к директории загрузок пользователя
func GetDownloadsDir() string {
	// Получаем домашнюю директорию пользователя
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Если не удалось получить домашнюю директорию, используем текущую
		return "."
	}

	// В зависимости от операционной системы, выбираем путь к директории загрузок
	switch runtime.GOOS {
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
