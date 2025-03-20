package interfaces

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"net/http"
)

// UserService определяет интерфейс для работы с пользователями
type UserService interface {
	// FindUser находит пользователя по логину
	FindUser(login string) (*domain.User, error)

	// SaveUser сохраняет нового пользователя с указанным логином и паролем
	SaveUser(login string, password string) (*domain.User, error)
}

// AuthService определяет интерфейс для аутентификации и авторизации
type AuthService interface {
	// GenerateToken генерирует JWT токен для указанного логина
	GenerateToken(login string) (string, error)

	// ValidateToken проверяет валидность JWT токена и возвращает claims
	ValidateToken(tokenString string) (*domain.Claims, error)

	// HashPassword хеширует пароль
	HashPassword(password string) (string, error)

	// CheckPasswordHash проверяет соответствие пароля его хешу
	CheckPasswordHash(password, hash string) bool

	// SetResponseAuthData устанавливает данные авторизации в HTTP-ответе
	SetResponseAuthData(w http.ResponseWriter, token string)
}

// TokenService определяет интерфейс для работы с токенами
type TokenService interface {
	// SaveToken сохраняет новый токен в хранилище
	SaveToken(token string)
}

// ClientService определяет интерфейс для клиентского сервиса
type ClientService interface {
	// Login выполняет запрос к API сервера для аутентификации пользователя и получения токена
	Login(login string, password string) (string, error)

	// Register выполняет запрос к API сервера для регистрации пользователя и получения токена
	Register(login string, password string) (string, error)

	// Методы UserService, необходимые для реализации интерфейса
	FindUser(login string) (*domain.User, error)
	SaveUser(login string, password string) (*domain.User, error)
}
