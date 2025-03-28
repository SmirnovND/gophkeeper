package usecase

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"os"
	"testing"
)

// MockTokenService - мок для интерфейса TokenService
type MockTokenServiceFixed struct {
	SaveTokenFunc func(token string)
	LoadTokenFunc func() (string, error)
}

func (m *MockTokenServiceFixed) SaveToken(token string) {
	if m.SaveTokenFunc != nil {
		m.SaveTokenFunc(token)
	}
}

func (m *MockTokenServiceFixed) LoadToken() (string, error) {
	if m.LoadTokenFunc != nil {
		return m.LoadTokenFunc()
	}
	return "", nil
}

// MockClientService - мок для интерфейса ClientService
type MockClientServiceFixed struct {
	LoginFunc                  func(login string, password string) (string, error)
	RegisterFunc               func(login string, password string) (string, error)
	GetUploadLinkFunc          func(label string, extension string, metadata string, token string) (string, error)
	GetDownloadLinkFunc        func(label string, token string) (string, *domain.FileMetadata, string, error)
	SendFileToServerFunc       func(url string, file *os.File) (string, error)
	DownloadFileFromServerFunc func(url string, outputPath string) error
	SaveTextFunc               func(label string, textData *domain.TextData, metadata string, token string) error
	GetTextFunc                func(label string, token string) (*domain.TextData, string, error)
	DeleteTextFunc             func(label string, token string) error
	SaveCardFunc               func(label string, cardData *domain.CardData, metadata string, token string) error
	GetCardFunc                func(label string, token string) (*domain.CardData, string, error)
	DeleteCardFunc             func(label string, token string) error
	SaveCredentialFunc         func(label string, credentialData *domain.CredentialData, metadata string, token string) error
	GetCredentialFunc          func(label string, token string) (*domain.CredentialData, string, error)
	DeleteCredentialFunc       func(label string, token string) error
}

func (m *MockClientServiceFixed) Login(login string, password string) (string, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(login, password)
	}
	return "", nil
}

func (m *MockClientServiceFixed) Register(login string, password string) (string, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(login, password)
	}
	return "", nil
}

func (m *MockClientServiceFixed) GetUploadLink(label string, extension string, metadata string, token string) (string, error) {
	if m.GetUploadLinkFunc != nil {
		return m.GetUploadLinkFunc(label, extension, metadata, token)
	}
	return "", nil
}

func (m *MockClientServiceFixed) GetDownloadLink(label string, token string) (string, *domain.FileMetadata, string, error) {
	if m.GetDownloadLinkFunc != nil {
		return m.GetDownloadLinkFunc(label, token)
	}
	return "", nil, "", nil
}

func (m *MockClientServiceFixed) SendFileToServer(url string, file *os.File) (string, error) {
	if m.SendFileToServerFunc != nil {
		return m.SendFileToServerFunc(url, file)
	}
	return "", nil
}

func (m *MockClientServiceFixed) DownloadFileFromServer(url string, outputPath string) error {
	if m.DownloadFileFromServerFunc != nil {
		return m.DownloadFileFromServerFunc(url, outputPath)
	}
	return nil
}

func (m *MockClientServiceFixed) SaveText(label string, textData *domain.TextData, metadata string, token string) error {
	if m.SaveTextFunc != nil {
		return m.SaveTextFunc(label, textData, metadata, token)
	}
	return nil
}

func (m *MockClientServiceFixed) GetText(label string, token string) (*domain.TextData, string, error) {
	if m.GetTextFunc != nil {
		return m.GetTextFunc(label, token)
	}
	return nil, "", nil
}

func (m *MockClientServiceFixed) DeleteText(label string, token string) error {
	if m.DeleteTextFunc != nil {
		return m.DeleteTextFunc(label, token)
	}
	return nil
}

func (m *MockClientServiceFixed) SaveCard(label string, cardData *domain.CardData, metadata string, token string) error {
	if m.SaveCardFunc != nil {
		return m.SaveCardFunc(label, cardData, metadata, token)
	}
	return nil
}

func (m *MockClientServiceFixed) GetCard(label string, token string) (*domain.CardData, string, error) {
	if m.GetCardFunc != nil {
		return m.GetCardFunc(label, token)
	}
	return nil, "", nil
}

func (m *MockClientServiceFixed) DeleteCard(label string, token string) error {
	if m.DeleteCardFunc != nil {
		return m.DeleteCardFunc(label, token)
	}
	return nil
}

func (m *MockClientServiceFixed) SaveCredential(label string, credentialData *domain.CredentialData, metadata string, token string) error {
	if m.SaveCredentialFunc != nil {
		return m.SaveCredentialFunc(label, credentialData, metadata, token)
	}
	return nil
}

func (m *MockClientServiceFixed) GetCredential(label string, token string) (*domain.CredentialData, string, error) {
	if m.GetCredentialFunc != nil {
		return m.GetCredentialFunc(label, token)
	}
	return nil, "", nil
}

func (m *MockClientServiceFixed) DeleteCredential(label string, token string) error {
	if m.DeleteCredentialFunc != nil {
		return m.DeleteCredentialFunc(label, token)
	}
	return nil
}

// TestClientUseCase_Upload_Success_Fixed тестирует успешную загрузку файла
func TestClientUseCase_Upload_Success_Fixed(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Записываем данные во временный файл
	_, err = tempFile.WriteString("test content")
	if err != nil {
		t.Fatalf("Ошибка при записи во временный файл: %v", err)
	}

	// Создаем моки
	mockTokenService := &MockTokenServiceFixed{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientServiceFixed{
		GetUploadLinkFunc: func(label string, extension string, metadata string, token string) (string, error) {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if extension != "txt" {
				t.Errorf("Ожидалось расширение 'txt', получено '%s'", extension)
			}
			// Не проверяем metadata, так как она вводится пользователем
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return "http://example.com/upload", nil
		},
		SendFileToServerFunc: func(url string, file *os.File) (string, error) {
			// Проверяем параметры
			if url != "http://example.com/upload" {
				t.Errorf("Ожидался URL 'http://example.com/upload', получен '%s'", url)
			}
			return "success", nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Upload
	result, err := clientUseCase.Upload(tempFile.Name(), "test_label")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове Upload: %v", err)
	}
	if result != "success" {
		t.Errorf("Ожидался результат 'success', получен '%s'", result)
	}
}
