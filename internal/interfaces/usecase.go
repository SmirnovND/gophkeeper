package interfaces

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"net/http"
)

// AuthUseCase определяет интерфейс для использования аутентификации
type AuthUseCase interface {
	// Login выполняет вход пользователя и возвращает JWT токен
	Login(w http.ResponseWriter, credentials *domain.Credentials) (string, error)

	// Register регистрирует нового пользователя и возвращает JWT токен
	Register(w http.ResponseWriter, credentials *domain.Credentials) (string, error)

	// ValidateToken проверяет валидность JWT токена и возвращает claims
	ValidateToken(token string) (*domain.Claims, error)
}

type ClientUseCase interface {
	Login(username string, password string) error
	Register(username string, password string, passwordCheck string) error
}

type FileUseCase interface {
	UploadFile(w http.ResponseWriter, fileData *domain.FileData)
}
