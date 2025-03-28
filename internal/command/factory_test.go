package command

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockClientUseCaseForFactory - мок для интерфейса ClientUseCase для тестов фабрики
type MockClientUseCaseForFactory struct {
	mock.Mock
}

func (m *MockClientUseCaseForFactory) Login(username string, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *MockClientUseCaseForFactory) Register(username string, password string, passwordCheck string) error {
	args := m.Called(username, password, passwordCheck)
	return args.Error(0)
}

func (m *MockClientUseCaseForFactory) Upload(filePath string, label string) (string, error) {
	args := m.Called(filePath, label)
	return args.String(0), args.Error(1)
}

func (m *MockClientUseCaseForFactory) Download(label string) error {
	args := m.Called(label)
	return args.Error(0)
}

func (m *MockClientUseCaseForFactory) SaveText(label string, textData *domain.TextData, metadata string) error {
	args := m.Called(label, textData, metadata)
	return args.Error(0)
}

func (m *MockClientUseCaseForFactory) GetText(label string) (*domain.TextData, string, error) {
	args := m.Called(label)
	var textData *domain.TextData
	if args.Get(0) != nil {
		textData = args.Get(0).(*domain.TextData)
	}
	return textData, args.String(1), args.Error(2)
}

func (m *MockClientUseCaseForFactory) DeleteText(label string) error {
	args := m.Called(label)
	return args.Error(0)
}

func (m *MockClientUseCaseForFactory) SaveCard(label string, cardData *domain.CardData, metadata string) error {
	args := m.Called(label, cardData, metadata)
	return args.Error(0)
}

func (m *MockClientUseCaseForFactory) GetCard(label string) (*domain.CardData, string, error) {
	args := m.Called(label)
	var cardData *domain.CardData
	if args.Get(0) != nil {
		cardData = args.Get(0).(*domain.CardData)
	}
	return cardData, args.String(1), args.Error(2)
}

func (m *MockClientUseCaseForFactory) DeleteCard(label string) error {
	args := m.Called(label)
	return args.Error(0)
}

func (m *MockClientUseCaseForFactory) SaveCredential(label string, credentialData *domain.CredentialData, metadata string) error {
	args := m.Called(label, credentialData, metadata)
	return args.Error(0)
}

func (m *MockClientUseCaseForFactory) GetCredential(label string) (*domain.CredentialData, string, error) {
	args := m.Called(label)
	var credentialData *domain.CredentialData
	if args.Get(0) != nil {
		credentialData = args.Get(0).(*domain.CredentialData)
	}
	return credentialData, args.String(1), args.Error(2)
}

func (m *MockClientUseCaseForFactory) DeleteCredential(label string) error {
	args := m.Called(label)
	return args.Error(0)
}

// Тест для функции NewCommand
func TestNewCommand(t *testing.T) {
	// Arrange
	mockClientUseCase := new(MockClientUseCaseForFactory)
	
	// Act
	cmd := NewCommand(mockClientUseCase)
	
	// Assert
	assert.NotNil(t, cmd)
	assert.Implements(t, (*interfaces.Command)(nil), cmd)
	
	// Проверяем, что команда содержит переданный ClientUseCase
	command, ok := cmd.(*Command)
	assert.True(t, ok)
	assert.Equal(t, mockClientUseCase, command.clientUseCase)
}

// Тест для проверки, что все методы интерфейса Command реализованы
func TestCommand_ImplementsCommandInterface(t *testing.T) {
	// Arrange
	mockClientUseCase := new(MockClientUseCaseForFactory)
	cmd := NewCommand(mockClientUseCase)
	
	// Act & Assert
	// Проверяем, что все методы интерфейса Command возвращают не nil
	assert.NotNil(t, cmd.Login())
	assert.NotNil(t, cmd.RegisterCmd())
	assert.NotNil(t, cmd.UploadCmd())
	assert.NotNil(t, cmd.DownloadCmd())
	
	assert.NotNil(t, cmd.SaveTextCmd())
	assert.NotNil(t, cmd.GetTextCmd())
	assert.NotNil(t, cmd.DeleteTextCmd())
	
	assert.NotNil(t, cmd.SaveCardCmd())
	assert.NotNil(t, cmd.GetCardCmd())
	assert.NotNil(t, cmd.DeleteCardCmd())
	
	assert.NotNil(t, cmd.SaveCredentialCmd())
	assert.NotNil(t, cmd.GetCredentialCmd())
	assert.NotNil(t, cmd.DeleteCredentialCmd())
}