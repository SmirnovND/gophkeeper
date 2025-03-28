package usecase

import (
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
)

// MockJwtService - мок для интерфейса JwtService
type MockJwtService struct {
	ExtractLoginFromTokenFunc func(tokenString string) (string, error)
}

func (m *MockJwtService) ExtractLoginFromToken(tokenString string) (string, error) {
	if m.ExtractLoginFromTokenFunc != nil {
		return m.ExtractLoginFromTokenFunc(tokenString)
	}
	if tokenString == "Bearer valid-token" {
		return "testuser", nil
	}
	if tokenString == "Bearer error-token" {
		return "", errors.New("ошибка извлечения логина из токена")
	}
	return "", errors.New("неверный токен")
}

// MockDataService - мок для интерфейса DataService
type MockDataService struct {
	GetCredentialFunc    func(login string, label string) (*domain.CredentialData, string, error)
	SaveCredentialFunc   func(login string, label string, credentialData *domain.CredentialData, metadata string) error
	DeleteCredentialFunc func(login string, label string) error

	GetCardFunc    func(login string, label string) (*domain.CardData, string, error)
	SaveCardFunc   func(login string, label string, cardData *domain.CardData, metadata string) error
	DeleteCardFunc func(login string, label string) error

	GetTextFunc    func(login string, label string) (*domain.TextData, string, error)
	SaveTextFunc   func(login string, label string, textData *domain.TextData, metadata string) error
	DeleteTextFunc func(login string, label string) error

	GetFileMetadataFunc    func(login string, label string) (*domain.FileMetadata, string, error)
	SaveFileMetadataFunc   func(login string, label string, fileData *domain.FileData, metadata string) error
	DeleteFileMetadataFunc func(login string, label string) error
}

// Реализация методов интерфейса DataService для мока
func (m *MockDataService) GetCredential(login string, label string) (*domain.CredentialData, string, error) {
	if m.GetCredentialFunc != nil {
		return m.GetCredentialFunc(login, label)
	}
	return nil, "", nil
}

func (m *MockDataService) SaveCredential(login string, label string, credentialData *domain.CredentialData, metadata string) error {
	if m.SaveCredentialFunc != nil {
		return m.SaveCredentialFunc(login, label, credentialData, metadata)
	}
	return nil
}

func (m *MockDataService) DeleteCredential(login string, label string) error {
	if m.DeleteCredentialFunc != nil {
		return m.DeleteCredentialFunc(login, label)
	}
	return nil
}

func (m *MockDataService) GetCard(login string, label string) (*domain.CardData, string, error) {
	if m.GetCardFunc != nil {
		return m.GetCardFunc(login, label)
	}
	return nil, "", nil
}

func (m *MockDataService) SaveCard(login string, label string, cardData *domain.CardData, metadata string) error {
	if m.SaveCardFunc != nil {
		return m.SaveCardFunc(login, label, cardData, metadata)
	}
	return nil
}

func (m *MockDataService) DeleteCard(login string, label string) error {
	if m.DeleteCardFunc != nil {
		return m.DeleteCardFunc(login, label)
	}
	return nil
}

func (m *MockDataService) GetText(login string, label string) (*domain.TextData, string, error) {
	if m.GetTextFunc != nil {
		return m.GetTextFunc(login, label)
	}
	return nil, "", nil
}

func (m *MockDataService) SaveText(login string, label string, textData *domain.TextData, metadata string) error {
	if m.SaveTextFunc != nil {
		return m.SaveTextFunc(login, label, textData, metadata)
	}
	return nil
}

func (m *MockDataService) DeleteText(login string, label string) error {
	if m.DeleteTextFunc != nil {
		return m.DeleteTextFunc(login, label)
	}
	return nil
}

func (m *MockDataService) GetFileMetadata(login string, label string) (*domain.FileMetadata, string, error) {
	if m.GetFileMetadataFunc != nil {
		return m.GetFileMetadataFunc(login, label)
	}
	return nil, "", nil
}

func (m *MockDataService) SaveFileMetadata(login string, label string, fileData *domain.FileData, metadata string) error {
	if m.SaveFileMetadataFunc != nil {
		return m.SaveFileMetadataFunc(login, label, fileData, metadata)
	}
	return nil
}

func (m *MockDataService) DeleteFileMetadata(login string, label string) error {
	if m.DeleteFileMetadataFunc != nil {
		return m.DeleteFileMetadataFunc(login, label)
	}
	return nil
}