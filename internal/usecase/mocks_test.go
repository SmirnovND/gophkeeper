package usecase

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
)

// MockDataService - мок для интерфейса DataService
type MockDataService struct {
	GetCredentialFunc    func(login string, label string) (*domain.CredentialData, error)
	SaveCredentialFunc   func(login string, label string, credentialData *domain.CredentialData) error
	DeleteCredentialFunc func(login string, label string) error

	GetCardFunc    func(login string, label string) (*domain.CardData, error)
	SaveCardFunc   func(login string, label string, cardData *domain.CardData) error
	DeleteCardFunc func(login string, label string) error

	GetTextFunc    func(login string, label string) (*domain.TextData, error)
	SaveTextFunc   func(login string, label string, textData *domain.TextData) error
	DeleteTextFunc func(login string, label string) error

	GetFileMetadataFunc    func(login string, label string) (*domain.FileMetadata, error)
	SaveFileMetadataFunc   func(login string, label string, fileData *domain.FileData) error
	DeleteFileMetadataFunc func(login string, label string) error
}

// Реализация методов интерфейса DataService для мока
func (m *MockDataService) GetCredential(login string, label string) (*domain.CredentialData, error) {
	if m.GetCredentialFunc != nil {
		return m.GetCredentialFunc(login, label)
	}
	return nil, nil
}

func (m *MockDataService) SaveCredential(login string, label string, credentialData *domain.CredentialData) error {
	if m.SaveCredentialFunc != nil {
		return m.SaveCredentialFunc(login, label, credentialData)
	}
	return nil
}

func (m *MockDataService) DeleteCredential(login string, label string) error {
	if m.DeleteCredentialFunc != nil {
		return m.DeleteCredentialFunc(login, label)
	}
	return nil
}

func (m *MockDataService) GetCard(login string, label string) (*domain.CardData, error) {
	if m.GetCardFunc != nil {
		return m.GetCardFunc(login, label)
	}
	return nil, nil
}

func (m *MockDataService) SaveCard(login string, label string, cardData *domain.CardData) error {
	if m.SaveCardFunc != nil {
		return m.SaveCardFunc(login, label, cardData)
	}
	return nil
}

func (m *MockDataService) DeleteCard(login string, label string) error {
	if m.DeleteCardFunc != nil {
		return m.DeleteCardFunc(login, label)
	}
	return nil
}

func (m *MockDataService) GetText(login string, label string) (*domain.TextData, error) {
	if m.GetTextFunc != nil {
		return m.GetTextFunc(login, label)
	}
	return nil, nil
}

func (m *MockDataService) SaveText(login string, label string, textData *domain.TextData) error {
	if m.SaveTextFunc != nil {
		return m.SaveTextFunc(login, label, textData)
	}
	return nil
}

func (m *MockDataService) DeleteText(login string, label string) error {
	if m.DeleteTextFunc != nil {
		return m.DeleteTextFunc(login, label)
	}
	return nil
}

func (m *MockDataService) GetFileMetadata(login string, label string) (*domain.FileMetadata, error) {
	if m.GetFileMetadataFunc != nil {
		return m.GetFileMetadataFunc(login, label)
	}
	return nil, nil
}

func (m *MockDataService) SaveFileMetadata(login string, label string, fileData *domain.FileData) error {
	if m.SaveFileMetadataFunc != nil {
		return m.SaveFileMetadataFunc(login, label, fileData)
	}
	return nil
}

func (m *MockDataService) DeleteFileMetadata(login string, label string) error {
	if m.DeleteFileMetadataFunc != nil {
		return m.DeleteFileMetadataFunc(login, label)
	}
	return nil
}