package interfaces

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"net/http"
	"os"
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

// TokenService определяет интерфейс для работы с токеном
type TokenService interface {
	// SaveToken сохраняет новый токен в хранилище
	SaveToken(token string)
	// LoadToken загружает токен из хранилища
	LoadToken() (string, error)
}

// ClientService определяет интерфейс для клиентского сервиса
type ClientService interface {
	// Login выполняет запрос к API сервера для аутентификации пользователя и получения токена
	Login(login string, password string) (string, error)

	// Register выполняет запрос к API сервера для регистрации пользователя и получения токена
	Register(login string, password string) (string, error)

	GetUploadLink(label string, extension string, token string) (string, error)

	GetDownloadLink(label string, token string) (string, *domain.FileMetadata, error)

	SendFileToServer(url string, file *os.File) (string, error)

	DownloadFileFromServer(url string, outputPath string) error

	// Методы для работы с текстовыми данными
	SaveText(label string, textData *domain.TextData, token string) error
	GetText(label string, token string) (*domain.TextData, error)
	DeleteText(label string, token string) error

	// Методы для работы с данными кредитных карт
	SaveCard(label string, cardData *domain.CardData, token string) error
	GetCard(label string, token string) (*domain.CardData, error)
	DeleteCard(label string, token string) error

	// Методы для работы с учетными данными
	SaveCredential(label string, credentialData *domain.CredentialData, token string) error
	GetCredential(label string, token string) (*domain.CredentialData, error)
	DeleteCredential(label string, token string) error
}

type CloudService interface {
	GenerateUploadLink(fileName string) (string, error)
	GenerateDownloadLink(fileName string) (string, error)
}

// DataService определяет интерфейс для работы с данными пользователя
type DataService interface {
	// Методы для работы с файлами
	SaveFileMetadata(login string, label string, fileData *domain.FileData) error
	GetFileMetadata(login string, label string) (*domain.FileMetadata, error)
	DeleteFileMetadata(login string, label string) error

	// Методы для работы с учетными данными (логин/пароль)
	SaveCredential(login string, label string, credentialData *domain.CredentialData) error
	GetCredential(login string, label string) (*domain.CredentialData, error)
	DeleteCredential(login string, label string) error

	// Методы для работы с данными кредитных карт
	SaveCard(login string, label string, cardData *domain.CardData) error
	GetCard(login string, label string) (*domain.CardData, error)
	DeleteCard(login string, label string) error

	// Методы для работы с текстовыми данными
	SaveText(login string, label string, textData *domain.TextData) error
	GetText(login string, label string) (*domain.TextData, error)
	DeleteText(login string, label string) error
}

type JwtService interface {
	ExtractLoginFromToken(tokenString string) (string, error)
}
