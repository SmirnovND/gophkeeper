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
	Upload(filePath string, label string) (string, error)
	Download(label string) error
	
	// Методы для работы с текстовыми данными
	SaveText(label string, textData *domain.TextData, metadata string) error
	GetText(label string) (*domain.TextData, string, error)
	DeleteText(label string) error
	
	// Методы для работы с данными кредитных карт
	SaveCard(label string, cardData *domain.CardData, metadata string) error
	GetCard(label string) (*domain.CardData, string, error)
	DeleteCard(label string) error
	
	// Методы для работы с учетными данными
	SaveCredential(label string, credentialData *domain.CredentialData, metadata string) error
	GetCredential(label string) (*domain.CredentialData, string, error)
	DeleteCredential(label string) error
}

type CloudUseCase interface {
	GenerateUploadLink(w http.ResponseWriter, r *http.Request, fileData *domain.FileData)
	GenerateDownloadLink(w http.ResponseWriter, r *http.Request, label string)
}

type DataUseCase interface {
	SaveCredential(w http.ResponseWriter, r *http.Request, label string, credentialData *domain.CredentialData, metadata string)
	GetCredential(w http.ResponseWriter, r *http.Request, label string)
	DeleteCredential(w http.ResponseWriter, r *http.Request, label string)

	SaveCard(w http.ResponseWriter, r *http.Request, label string, cardData *domain.CardData, metadata string)
	GetCard(w http.ResponseWriter, r *http.Request, label string)
	DeleteCard(w http.ResponseWriter, r *http.Request, label string)

	SaveText(w http.ResponseWriter, r *http.Request, label string, textData *domain.TextData, metadata string)
	GetText(w http.ResponseWriter, r *http.Request, label string)
	DeleteText(w http.ResponseWriter, r *http.Request, label string)
}
