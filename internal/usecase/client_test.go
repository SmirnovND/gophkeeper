package usecase

import (
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"os"
	"path/filepath"
	"testing"
)

// MockTokenService - мок для интерфейса TokenService
type MockTokenService struct {
	SaveTokenFunc func(token string)
	LoadTokenFunc func() (string, error)
}

func (m *MockTokenService) SaveToken(token string) {
	m.SaveTokenFunc(token)
}

func (m *MockTokenService) LoadToken() (string, error) {
	return m.LoadTokenFunc()
}

// MockClientService - мок для интерфейса ClientService
type MockClientService struct {
	LoginFunc                  func(login string, password string) (string, error)
	RegisterFunc               func(login string, password string) (string, error)
	GetUploadLinkFunc          func(label string, extension string, token string) (string, error)
	GetDownloadLinkFunc        func(label string, token string) (string, *domain.FileMetadata, error)
	SendFileToServerFunc       func(url string, file *os.File) (string, error)
	DownloadFileFromServerFunc func(url string, outputPath string) error
	SaveTextFunc               func(label string, textData *domain.TextData, token string) error
	GetTextFunc                func(label string, token string) (*domain.TextData, error)
	DeleteTextFunc             func(label string, token string) error
	SaveCardFunc               func(label string, cardData *domain.CardData, token string) error
	GetCardFunc                func(label string, token string) (*domain.CardData, error)
	DeleteCardFunc             func(label string, token string) error
	SaveCredentialFunc         func(label string, credentialData *domain.CredentialData, token string) error
	GetCredentialFunc          func(label string, token string) (*domain.CredentialData, error)
	DeleteCredentialFunc       func(label string, token string) error
}

func (m *MockClientService) Login(login string, password string) (string, error) {
	return m.LoginFunc(login, password)
}

func (m *MockClientService) Register(login string, password string) (string, error) {
	return m.RegisterFunc(login, password)
}

func (m *MockClientService) GetUploadLink(label string, extension string, token string) (string, error) {
	return m.GetUploadLinkFunc(label, extension, token)
}

func (m *MockClientService) GetDownloadLink(label string, token string) (string, *domain.FileMetadata, error) {
	return m.GetDownloadLinkFunc(label, token)
}

func (m *MockClientService) SendFileToServer(url string, file *os.File) (string, error) {
	return m.SendFileToServerFunc(url, file)
}

func (m *MockClientService) DownloadFileFromServer(url string, outputPath string) error {
	return m.DownloadFileFromServerFunc(url, outputPath)
}

func (m *MockClientService) SaveText(label string, textData *domain.TextData, token string) error {
	return m.SaveTextFunc(label, textData, token)
}

func (m *MockClientService) GetText(label string, token string) (*domain.TextData, error) {
	return m.GetTextFunc(label, token)
}

func (m *MockClientService) DeleteText(label string, token string) error {
	return m.DeleteTextFunc(label, token)
}

func (m *MockClientService) SaveCard(label string, cardData *domain.CardData, token string) error {
	return m.SaveCardFunc(label, cardData, token)
}

func (m *MockClientService) GetCard(label string, token string) (*domain.CardData, error) {
	return m.GetCardFunc(label, token)
}

func (m *MockClientService) DeleteCard(label string, token string) error {
	return m.DeleteCardFunc(label, token)
}

func (m *MockClientService) SaveCredential(label string, credentialData *domain.CredentialData, token string) error {
	return m.SaveCredentialFunc(label, credentialData, token)
}

func (m *MockClientService) GetCredential(label string, token string) (*domain.CredentialData, error) {
	return m.GetCredentialFunc(label, token)
}

func (m *MockClientService) DeleteCredential(label string, token string) error {
	return m.DeleteCredentialFunc(label, token)
}

// TestClientUseCase_Login_Success тестирует успешный вход пользователя
func TestClientUseCase_Login_Success(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		SaveTokenFunc: func(token string) {
			// Проверяем, что токен сохраняется
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
		},
	}

	mockClientService := &MockClientService{
		LoginFunc: func(login string, password string) (string, error) {
			// Проверяем параметры
			if login != "testuser" || password != "testpassword" {
				t.Errorf("Ожидались логин 'testuser' и пароль 'testpassword', получены '%s' и '%s'", login, password)
			}
			return "test_token", nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Login
	err := clientUseCase.Login("testuser", "testpassword")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при входе: %v", err)
	}
}

// TestClientUseCase_Login_Error тестирует ошибку при входе пользователя
func TestClientUseCase_Login_Error(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{}

	mockClientService := &MockClientService{
		LoginFunc: func(login string, password string) (string, error) {
			return "", errors.New("ошибка аутентификации")
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Login
	err := clientUseCase.Login("testuser", "testpassword")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Register_Success тестирует успешную регистрацию пользователя
func TestClientUseCase_Register_Success(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		SaveTokenFunc: func(token string) {
			// Проверяем, что токен сохраняется
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
		},
	}

	mockClientService := &MockClientService{
		RegisterFunc: func(login string, password string) (string, error) {
			// Проверяем параметры
			if login != "testuser" || password != "testpassword" {
				t.Errorf("Ожидались логин 'testuser' и пароль 'testpassword', получены '%s' и '%s'", login, password)
			}
			return "test_token", nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Register
	err := clientUseCase.Register("testuser", "testpassword", "testpassword")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при регистрации: %v", err)
	}
}

// TestClientUseCase_Register_PasswordMismatch тестирует ошибку при несовпадении паролей
func TestClientUseCase_Register_PasswordMismatch(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{}
	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Register с разными паролями
	err := clientUseCase.Register("testuser", "testpassword", "differentpassword")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Register_Error тестирует ошибку при регистрации пользователя
func TestClientUseCase_Register_Error(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{}

	mockClientService := &MockClientService{
		RegisterFunc: func(login string, password string) (string, error) {
			return "", errors.New("ошибка регистрации")
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Register
	err := clientUseCase.Register("testuser", "testpassword", "testpassword")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_SaveText_Success тестирует успешное сохранение текстовых данных
func TestClientUseCase_SaveText_Success(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		SaveTextFunc: func(label string, textData *domain.TextData, token string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if textData.Content != "test content" {
				t.Errorf("Ожидалось содержимое 'test content', получено '%s'", textData.Content)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Создаем тестовые текстовые данные
	textData := &domain.TextData{
		Content: "test content",
	}

	// Вызываем метод SaveText
	err := clientUseCase.SaveText("test_label", textData)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при сохранении текстовых данных: %v", err)
	}
}

// TestClientUseCase_SaveText_EmptyLabel тестирует ошибку при пустой метке
func TestClientUseCase_SaveText_EmptyLabel(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{}
	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Создаем тестовые текстовые данные
	textData := &domain.TextData{
		Content: "test content",
	}

	// Вызываем метод SaveText с пустой меткой
	err := clientUseCase.SaveText("", textData)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_SaveText_LoadTokenError тестирует ошибку при загрузке токена
func TestClientUseCase_SaveText_LoadTokenError(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "", errors.New("ошибка загрузки токена")
		},
	}

	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Создаем тестовые текстовые данные
	textData := &domain.TextData{
		Content: "test content",
	}

	// Вызываем метод SaveText
	err := clientUseCase.SaveText("test_label", textData)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_SaveText_SaveError тестирует ошибку при сохранении текстовых данных
func TestClientUseCase_SaveText_SaveError(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		SaveTextFunc: func(label string, textData *domain.TextData, token string) error {
			return errors.New("ошибка сохранения текстовых данных")
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Создаем тестовые текстовые данные
	textData := &domain.TextData{
		Content: "test content",
	}

	// Вызываем метод SaveText
	err := clientUseCase.SaveText("test_label", textData)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_GetText_Success тестирует успешное получение текстовых данных
func TestClientUseCase_GetText_Success(t *testing.T) {
	// Создаем ожидаемые текстовые данные
	expectedTextData := &domain.TextData{
		Content: "test content",
	}

	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetTextFunc: func(label string, token string) (*domain.TextData, error) {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return expectedTextData, nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод GetText
	textData, err := clientUseCase.GetText("test_label")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении текстовых данных: %v", err)
	}

	if textData != expectedTextData {
		t.Errorf("Ожидались текстовые данные %v, получены %v", expectedTextData, textData)
	}
}

// TestClientUseCase_GetText_EmptyLabel тестирует ошибку при пустой метке
func TestClientUseCase_GetText_EmptyLabel(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{}
	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод GetText с пустой меткой
	_, err := clientUseCase.GetText("")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_GetText_LoadTokenError тестирует ошибку при загрузке токена
func TestClientUseCase_GetText_LoadTokenError(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "", errors.New("ошибка загрузки токена")
		},
	}

	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод GetText
	_, err := clientUseCase.GetText("test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_GetText_GetError тестирует ошибку при получении текстовых данных
func TestClientUseCase_GetText_GetError(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetTextFunc: func(label string, token string) (*domain.TextData, error) {
			return nil, errors.New("ошибка получения текстовых данных")
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод GetText
	_, err := clientUseCase.GetText("test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_DeleteText_Success тестирует успешное удаление текстовых данных
func TestClientUseCase_DeleteText_Success(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		DeleteTextFunc: func(label string, token string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод DeleteText
	err := clientUseCase.DeleteText("test_label")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении текстовых данных: %v", err)
	}
}

// TestClientUseCase_DeleteText_EmptyLabel тестирует ошибку при пустой метке
func TestClientUseCase_DeleteText_EmptyLabel(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{}
	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод DeleteText с пустой меткой
	err := clientUseCase.DeleteText("")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_DeleteText_LoadTokenError тестирует ошибку при загрузке токена
func TestClientUseCase_DeleteText_LoadTokenError(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "", errors.New("ошибка загрузки токена")
		},
	}

	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод DeleteText
	err := clientUseCase.DeleteText("test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_DeleteText_DeleteError тестирует ошибку при удалении текстовых данных
func TestClientUseCase_DeleteText_DeleteError(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		DeleteTextFunc: func(label string, token string) error {
			return errors.New("ошибка удаления текстовых данных")
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод DeleteText
	err := clientUseCase.DeleteText("test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_SaveCard_Success тестирует успешное сохранение данных карты
func TestClientUseCase_SaveCard_Success(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		SaveCardFunc: func(label string, cardData *domain.CardData, token string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if cardData.Number != "1234567890123456" {
				t.Errorf("Ожидался номер карты '1234567890123456', получен '%s'", cardData.Number)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Создаем тестовые данные карты
	cardData := &domain.CardData{
		Number:     "1234567890123456",
		Holder:     "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}

	// Вызываем метод SaveCard
	err := clientUseCase.SaveCard("test_label", cardData)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при сохранении данных карты: %v", err)
	}
}

// TestClientUseCase_GetCard_Success тестирует успешное получение данных карты
func TestClientUseCase_GetCard_Success(t *testing.T) {
	// Создаем ожидаемые данные карты
	expectedCardData := &domain.CardData{
		Number:     "1234567890123456",
		Holder:     "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}

	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetCardFunc: func(label string, token string) (*domain.CardData, error) {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return expectedCardData, nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод GetCard
	cardData, err := clientUseCase.GetCard("test_label")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении данных карты: %v", err)
	}

	if cardData != expectedCardData {
		t.Errorf("Ожидались данные карты %v, получены %v", expectedCardData, cardData)
	}
}

// TestClientUseCase_DeleteCard_Success тестирует успешное удаление данных карты
func TestClientUseCase_DeleteCard_Success(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		DeleteCardFunc: func(label string, token string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод DeleteCard
	err := clientUseCase.DeleteCard("test_label")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении данных карты: %v", err)
	}
}

// TestClientUseCase_SaveCredential_Success тестирует успешное сохранение учетных данных
func TestClientUseCase_SaveCredential_Success(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		SaveCredentialFunc: func(label string, credentialData *domain.CredentialData, token string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if credentialData.Login != "testlogin" {
				t.Errorf("Ожидался логин 'testlogin', получен '%s'", credentialData.Login)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Создаем тестовые учетные данные
	credentialData := &domain.CredentialData{
		Login:    "testlogin",
		Password: "testpassword",
	}

	// Вызываем метод SaveCredential
	err := clientUseCase.SaveCredential("test_label", credentialData)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при сохранении учетных данных: %v", err)
	}
}

// TestClientUseCase_GetCredential_Success тестирует успешное получение учетных данных
func TestClientUseCase_GetCredential_Success(t *testing.T) {
	// Создаем ожидаемые учетные данные
	expectedCredentialData := &domain.CredentialData{
		Login:    "testlogin",
		Password: "testpassword",
	}

	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetCredentialFunc: func(label string, token string) (*domain.CredentialData, error) {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return expectedCredentialData, nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод GetCredential
	credentialData, err := clientUseCase.GetCredential("test_label")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении учетных данных: %v", err)
	}

	if credentialData != expectedCredentialData {
		t.Errorf("Ожидались учетные данные %v, получены %v", expectedCredentialData, credentialData)
	}
}

// TestClientUseCase_DeleteCredential_Success тестирует успешное удаление учетных данных
func TestClientUseCase_DeleteCredential_Success(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		DeleteCredentialFunc: func(label string, token string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод DeleteCredential
	err := clientUseCase.DeleteCredential("test_label")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении учетных данных: %v", err)
	}
}

// TestClientUseCase_Download_Success тестирует успешное скачивание файла
func TestClientUseCase_Download_Success(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetDownloadLinkFunc: func(label string, token string) (string, *domain.FileMetadata, error) {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if token != "test_token" {
				t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
			}
			return "http://example.com/download", &domain.FileMetadata{
				FileName:  "test_file",
				Extension: "txt",
			}, nil
		},
		DownloadFileFromServerFunc: func(url string, outputPath string) error {
			// Проверяем параметры
			if url != "http://example.com/download" {
				t.Errorf("Ожидался URL 'http://example.com/download', получен '%s'", url)
			}
			// Проверяем, что путь к файлу содержит имя файла
			expectedFileName := filepath.Join(filepath.Dir(outputPath), "test_label.txt")
			if outputPath != expectedFileName {
				t.Errorf("Ожидался путь к файлу '%s', получен '%s'", expectedFileName, outputPath)
			}
			return nil
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Download
	err := clientUseCase.Download("test_label")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при скачивании файла: %v", err)
	}
}

// TestClientUseCase_Download_EmptyLabel тестирует ошибку при пустой метке
func TestClientUseCase_Download_EmptyLabel(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{}
	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Download с пустой меткой
	err := clientUseCase.Download("")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Download_LoadTokenError тестирует ошибку при загрузке токена
func TestClientUseCase_Download_LoadTokenError(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "", errors.New("ошибка загрузки токена")
		},
	}

	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Download
	err := clientUseCase.Download("test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Download_GetDownloadLinkError тестирует ошибку при получении ссылки на скачивание
func TestClientUseCase_Download_GetDownloadLinkError(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetDownloadLinkFunc: func(label string, token string) (string, *domain.FileMetadata, error) {
			return "", nil, errors.New("ошибка получения ссылки на скачивание")
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Download
	err := clientUseCase.Download("test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Download_DownloadFileError тестирует ошибку при скачивании файла
func TestClientUseCase_Download_DownloadFileError(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetDownloadLinkFunc: func(label string, token string) (string, *domain.FileMetadata, error) {
			return "http://example.com/download", &domain.FileMetadata{
				FileName:  "test_file",
				Extension: "txt",
			}, nil
		},
		DownloadFileFromServerFunc: func(url string, outputPath string) error {
			return errors.New("ошибка скачивания файла")
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Download
	err := clientUseCase.Download("test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Upload_Success тестирует успешную загрузку файла
func TestClientUseCase_Upload_Success(t *testing.T) {
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
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetUploadLinkFunc: func(label string, extension string, token string) (string, error) {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if extension != "txt" {
				t.Errorf("Ожидалось расширение 'txt', получено '%s'", extension)
			}
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
		t.Fatalf("Ошибка при загрузке файла: %v", err)
	}

	if result != "success" {
		t.Errorf("Ожидался результат 'success', получен '%s'", result)
	}
}

// TestClientUseCase_Upload_FileNotFound тестирует ошибку при отсутствии файла
func TestClientUseCase_Upload_FileNotFound(t *testing.T) {
	// Создаем моки
	mockTokenService := &MockTokenService{}
	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Upload с несуществующим файлом
	_, err := clientUseCase.Upload("/non/existent/file.txt", "test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Upload_IsDirectory тестирует ошибку при указании директории вместо файла
func TestClientUseCase_Upload_IsDirectory(t *testing.T) {
	// Создаем временную директорию для тестирования
	tempDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatalf("Ошибка при создании временной директории: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем моки
	mockTokenService := &MockTokenService{}
	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Upload с директорией вместо файла
	_, err = clientUseCase.Upload(tempDir, "test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Upload_UnsupportedFileType тестирует ошибку при неподдерживаемом типе файла
func TestClientUseCase_Upload_UnsupportedFileType(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "test_*.dat")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Создаем моки
	mockTokenService := &MockTokenService{}
	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Upload с файлом неподдерживаемого типа
	_, err = clientUseCase.Upload(tempFile.Name(), "test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Upload_LoadTokenError тестирует ошибку при загрузке токена
func TestClientUseCase_Upload_LoadTokenError(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "", errors.New("ошибка загрузки токена")
		},
	}

	mockClientService := &MockClientService{}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Upload
	_, err = clientUseCase.Upload(tempFile.Name(), "test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Upload_GetUploadLinkError тестирует ошибку при получении ссылки на загрузку
func TestClientUseCase_Upload_GetUploadLinkError(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetUploadLinkFunc: func(label string, extension string, token string) (string, error) {
			return "", errors.New("ошибка получения ссылки на загрузку")
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Upload
	_, err = clientUseCase.Upload(tempFile.Name(), "test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestClientUseCase_Upload_SendFileError тестирует ошибку при отправке файла на сервер
func TestClientUseCase_Upload_SendFileError(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Создаем моки
	mockTokenService := &MockTokenService{
		LoadTokenFunc: func() (string, error) {
			return "test_token", nil
		},
	}

	mockClientService := &MockClientService{
		GetUploadLinkFunc: func(label string, extension string, token string) (string, error) {
			return "http://example.com/upload", nil
		},
		SendFileToServerFunc: func(url string, file *os.File) (string, error) {
			return "", errors.New("ошибка отправки файла на сервер")
		},
	}

	// Создаем экземпляр ClientUseCase
	clientUseCase := NewClientUseCase(mockTokenService, mockClientService)

	// Вызываем метод Upload
	_, err = clientUseCase.Upload(tempFile.Name(), "test_label")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}
