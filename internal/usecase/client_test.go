package usecase

import (
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"os"
	"path/filepath"
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

// TestClientUseCase_Login тестирует метод Login
func TestClientUseCase_Login(t *testing.T) {
	// Тест успешного входа
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			SaveTokenFunc: func(token string) {
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
			},
		}

		mockClientService := &MockClientServiceFixed{
			LoginFunc: func(login string, password string) (string, error) {
				if login != "testuser" || password != "testpass" {
					t.Errorf("Ожидались логин 'testuser' и пароль 'testpass', получены '%s' и '%s'", login, password)
				}
				return "test-token", nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Login("testuser", "testpass")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
	})

	// Тест ошибки при входе
	t.Run("Error", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{
			LoginFunc: func(login string, password string) (string, error) {
				return "", errors.New("ошибка аутентификации")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Login("testuser", "testpass")
		if err == nil {
			t.Error("Ожидалась ошибка, но ее не было")
		}
	})
}

// TestClientUseCase_Register тестирует метод Register
func TestClientUseCase_Register(t *testing.T) {
	// Тест успешной регистрации
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			SaveTokenFunc: func(token string) {
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
			},
		}

		mockClientService := &MockClientServiceFixed{
			RegisterFunc: func(login string, password string) (string, error) {
				if login != "testuser" || password != "testpass" {
					t.Errorf("Ожидались логин 'testuser' и пароль 'testpass', получены '%s' и '%s'", login, password)
				}
				return "test-token", nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Register("testuser", "testpass", "testpass")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
	})

	// Тест ошибки при несовпадении паролей
	t.Run("PasswordMismatch", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Register("testuser", "testpass", "wrongpass")
		if err == nil {
			t.Error("Ожидалась ошибка несовпадения паролей, но ее не было")
		}
	})

	// Тест ошибки при регистрации
	t.Run("RegistrationError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{
			RegisterFunc: func(login string, password string) (string, error) {
				return "", errors.New("ошибка регистрации")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Register("testuser", "testpass", "testpass")
		if err == nil {
			t.Error("Ожидалась ошибка регистрации, но ее не было")
		}
	})
}

// TestClientUseCase_Download тестирует метод Download
func TestClientUseCase_Download(t *testing.T) {
	// Тест успешного скачивания
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			GetDownloadLinkFunc: func(label string, token string) (string, *domain.FileMetadata, string, error) {
				if label != "test-file" {
					t.Errorf("Ожидалась метка 'test-file', получена '%s'", label)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return "http://example.com/download", &domain.FileMetadata{
					FileName:  "test-file",
					Extension: "txt",
				}, "test metadata", nil
			},
			DownloadFileFromServerFunc: func(url string, outputPath string) error {
				if url != "http://example.com/download" {
					t.Errorf("Ожидался URL 'http://example.com/download', получен '%s'", url)
				}
				// Проверяем, что путь содержит правильное имя файла
				expectedFilename := "test-file.txt"
				if filepath.Base(outputPath) != expectedFilename {
					t.Errorf("Ожидалось имя файла '%s', получено '%s'", expectedFilename, filepath.Base(outputPath))
				}
				return nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Download("test-file")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Download("")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Download("test-file")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при получении ссылки на скачивание
	t.Run("GetDownloadLinkError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			GetDownloadLinkFunc: func(label string, token string) (string, *domain.FileMetadata, string, error) {
				return "", nil, "", errors.New("ошибка получения ссылки")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Download("test-file")
		if err == nil {
			t.Error("Ожидалась ошибка получения ссылки, но ее не было")
		}
	})

	// Тест ошибки при скачивании файла
	t.Run("DownloadFileError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			GetDownloadLinkFunc: func(label string, token string) (string, *domain.FileMetadata, string, error) {
				return "http://example.com/download", &domain.FileMetadata{
					FileName:  "test-file",
					Extension: "txt",
				}, "", nil
			},
			DownloadFileFromServerFunc: func(url string, outputPath string) error {
				return errors.New("ошибка скачивания файла")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.Download("test-file")
		if err == nil {
			t.Error("Ожидалась ошибка скачивания файла, но ее не было")
		}
	})
}

// TestClientUseCase_SaveText тестирует метод SaveText
func TestClientUseCase_SaveText(t *testing.T) {
	// Тест успешного сохранения текста
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			SaveTextFunc: func(label string, textData *domain.TextData, metadata string, token string) error {
				if label != "test-text" {
					t.Errorf("Ожидалась метка 'test-text', получена '%s'", label)
				}
				if textData.Content != "test content" {
					t.Errorf("Ожидался текст 'test content', получен '%s'", textData.Content)
				}
				if metadata != "test metadata" {
					t.Errorf("Ожидались метаданные 'test metadata', получены '%s'", metadata)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		textData := &domain.TextData{Content: "test content"}
		err := clientUseCase.SaveText("test-text", textData, "test metadata")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		textData := &domain.TextData{Content: "test content"}
		err := clientUseCase.SaveText("", textData, "test metadata")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		textData := &domain.TextData{Content: "test content"}
		err := clientUseCase.SaveText("test-text", textData, "test metadata")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при сохранении текста
	t.Run("SaveTextError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			SaveTextFunc: func(label string, textData *domain.TextData, metadata string, token string) error {
				return errors.New("ошибка сохранения текста")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		textData := &domain.TextData{Content: "test content"}
		err := clientUseCase.SaveText("test-text", textData, "test metadata")
		if err == nil {
			t.Error("Ожидалась ошибка сохранения текста, но ее не было")
		}
	})
}

// TestClientUseCase_GetText тестирует метод GetText
func TestClientUseCase_GetText(t *testing.T) {
	// Тест успешного получения текста
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			GetTextFunc: func(label string, token string) (*domain.TextData, string, error) {
				if label != "test-text" {
					t.Errorf("Ожидалась метка 'test-text', получена '%s'", label)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return &domain.TextData{Content: "test content"}, "test metadata", nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		textData, metadata, err := clientUseCase.GetText("test-text")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
		if textData.Content != "test content" {
			t.Errorf("Ожидался текст 'test content', получен '%s'", textData.Content)
		}
		if metadata != "test metadata" {
			t.Errorf("Ожидались метаданные 'test metadata', получены '%s'", metadata)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		_, _, err := clientUseCase.GetText("")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		_, _, err := clientUseCase.GetText("test-text")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при получении текста
	t.Run("GetTextError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			GetTextFunc: func(label string, token string) (*domain.TextData, string, error) {
				return nil, "", errors.New("ошибка получения текста")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		_, _, err := clientUseCase.GetText("test-text")
		if err == nil {
			t.Error("Ожидалась ошибка получения текста, но ее не было")
		}
	})
}

// TestClientUseCase_DeleteText тестирует метод DeleteText
func TestClientUseCase_DeleteText(t *testing.T) {
	// Тест успешного удаления текста
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			DeleteTextFunc: func(label string, token string) error {
				if label != "test-text" {
					t.Errorf("Ожидалась метка 'test-text', получена '%s'", label)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteText("test-text")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteText("")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteText("test-text")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при удалении текста
	t.Run("DeleteTextError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			DeleteTextFunc: func(label string, token string) error {
				return errors.New("ошибка удаления текста")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteText("test-text")
		if err == nil {
			t.Error("Ожидалась ошибка удаления текста, но ее не было")
		}
	})
}

// TestClientUseCase_SaveCard тестирует метод SaveCard
func TestClientUseCase_SaveCard(t *testing.T) {
	// Тест успешного сохранения карты
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			SaveCardFunc: func(label string, cardData *domain.CardData, metadata string, token string) error {
				if label != "test-card" {
					t.Errorf("Ожидалась метка 'test-card', получена '%s'", label)
				}
				if cardData.Number != "1234567890123456" {
					t.Errorf("Ожидался номер карты '1234567890123456', получен '%s'", cardData.Number)
				}
				if metadata != "test metadata" {
					t.Errorf("Ожидались метаданные 'test metadata', получены '%s'", metadata)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		cardData := &domain.CardData{
			Number:     "1234567890123456",
			Holder:     "Test User",
			ExpiryDate: "12/25",
			CVV:        "123",
		}
		err := clientUseCase.SaveCard("test-card", cardData, "test metadata")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		cardData := &domain.CardData{
			Number:     "1234567890123456",
			Holder:     "Test User",
			ExpiryDate: "12/25",
			CVV:        "123",
		}
		err := clientUseCase.SaveCard("", cardData, "test metadata")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		cardData := &domain.CardData{
			Number:     "1234567890123456",
			Holder:     "Test User",
			ExpiryDate: "12/25",
			CVV:        "123",
		}
		err := clientUseCase.SaveCard("test-card", cardData, "test metadata")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при сохранении карты
	t.Run("SaveCardError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			SaveCardFunc: func(label string, cardData *domain.CardData, metadata string, token string) error {
				return errors.New("ошибка сохранения карты")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		cardData := &domain.CardData{
			Number:     "1234567890123456",
			Holder:     "Test User",
			ExpiryDate: "12/25",
			CVV:        "123",
		}
		err := clientUseCase.SaveCard("test-card", cardData, "test metadata")
		if err == nil {
			t.Error("Ожидалась ошибка сохранения карты, но ее не было")
		}
	})
}

// TestClientUseCase_GetCard тестирует метод GetCard
func TestClientUseCase_GetCard(t *testing.T) {
	// Тест успешного получения карты
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			GetCardFunc: func(label string, token string) (*domain.CardData, string, error) {
				if label != "test-card" {
					t.Errorf("Ожидалась метка 'test-card', получена '%s'", label)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return &domain.CardData{
					Number:     "1234567890123456",
					Holder:     "Test User",
					ExpiryDate: "12/25",
					CVV:        "123",
				}, "test metadata", nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		cardData, metadata, err := clientUseCase.GetCard("test-card")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
		if cardData.Number != "1234567890123456" {
			t.Errorf("Ожидался номер карты '1234567890123456', получен '%s'", cardData.Number)
		}
		if metadata != "test metadata" {
			t.Errorf("Ожидались метаданные 'test metadata', получены '%s'", metadata)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		_, _, err := clientUseCase.GetCard("")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		_, _, err := clientUseCase.GetCard("test-card")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при получении карты
	t.Run("GetCardError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			GetCardFunc: func(label string, token string) (*domain.CardData, string, error) {
				return nil, "", errors.New("ошибка получения карты")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		_, _, err := clientUseCase.GetCard("test-card")
		if err == nil {
			t.Error("Ожидалась ошибка получения карты, но ее не было")
		}
	})
}

// TestClientUseCase_DeleteCard тестирует метод DeleteCard
func TestClientUseCase_DeleteCard(t *testing.T) {
	// Тест успешного удаления карты
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			DeleteCardFunc: func(label string, token string) error {
				if label != "test-card" {
					t.Errorf("Ожидалась метка 'test-card', получена '%s'", label)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteCard("test-card")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteCard("")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteCard("test-card")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при удалении карты
	t.Run("DeleteCardError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			DeleteCardFunc: func(label string, token string) error {
				return errors.New("ошибка удаления карты")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteCard("test-card")
		if err == nil {
			t.Error("Ожидалась ошибка удаления карты, но ее не было")
		}
	})
}

// TestClientUseCase_SaveCredential тестирует метод SaveCredential
func TestClientUseCase_SaveCredential(t *testing.T) {
	// Тест успешного сохранения учетных данных
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			SaveCredentialFunc: func(label string, credentialData *domain.CredentialData, metadata string, token string) error {
				if label != "test-credential" {
					t.Errorf("Ожидалась метка 'test-credential', получена '%s'", label)
				}
				if credentialData.Login != "testuser" {
					t.Errorf("Ожидался логин 'testuser', получен '%s'", credentialData.Login)
				}
				if credentialData.Password != "testpass" {
					t.Errorf("Ожидался пароль 'testpass', получен '%s'", credentialData.Password)
				}
				if metadata != "test metadata" {
					t.Errorf("Ожидались метаданные 'test metadata', получены '%s'", metadata)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		credentialData := &domain.CredentialData{
			Login:    "testuser",
			Password: "testpass",
		}
		err := clientUseCase.SaveCredential("test-credential", credentialData, "test metadata")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		credentialData := &domain.CredentialData{
			Login:    "testuser",
			Password: "testpass",
		}
		err := clientUseCase.SaveCredential("", credentialData, "test metadata")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		credentialData := &domain.CredentialData{
			Login:    "testuser",
			Password: "testpass",
		}
		err := clientUseCase.SaveCredential("test-credential", credentialData, "test metadata")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при сохранении учетных данных
	t.Run("SaveCredentialError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			SaveCredentialFunc: func(label string, credentialData *domain.CredentialData, metadata string, token string) error {
				return errors.New("ошибка сохранения учетных данных")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		credentialData := &domain.CredentialData{
			Login:    "testuser",
			Password: "testpass",
		}
		err := clientUseCase.SaveCredential("test-credential", credentialData, "test metadata")
		if err == nil {
			t.Error("Ожидалась ошибка сохранения учетных данных, но ее не было")
		}
	})
}

// TestClientUseCase_GetCredential тестирует метод GetCredential
func TestClientUseCase_GetCredential(t *testing.T) {
	// Тест успешного получения учетных данных
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			GetCredentialFunc: func(label string, token string) (*domain.CredentialData, string, error) {
				if label != "test-credential" {
					t.Errorf("Ожидалась метка 'test-credential', получена '%s'", label)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return &domain.CredentialData{
					Login:    "testuser",
					Password: "testpass",
				}, "test metadata", nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		credentialData, metadata, err := clientUseCase.GetCredential("test-credential")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
		if credentialData.Login != "testuser" {
			t.Errorf("Ожидался логин 'testuser', получен '%s'", credentialData.Login)
		}
		if credentialData.Password != "testpass" {
			t.Errorf("Ожидался пароль 'testpass', получен '%s'", credentialData.Password)
		}
		if metadata != "test metadata" {
			t.Errorf("Ожидались метаданные 'test metadata', получены '%s'", metadata)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		_, _, err := clientUseCase.GetCredential("")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		_, _, err := clientUseCase.GetCredential("test-credential")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при получении учетных данных
	t.Run("GetCredentialError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			GetCredentialFunc: func(label string, token string) (*domain.CredentialData, string, error) {
				return nil, "", errors.New("ошибка получения учетных данных")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		_, _, err := clientUseCase.GetCredential("test-credential")
		if err == nil {
			t.Error("Ожидалась ошибка получения учетных данных, но ее не было")
		}
	})
}

// TestClientUseCase_DeleteCredential тестирует метод DeleteCredential
func TestClientUseCase_DeleteCredential(t *testing.T) {
	// Тест успешного удаления учетных данных
	t.Run("Success", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}

		mockClientService := &MockClientServiceFixed{
			DeleteCredentialFunc: func(label string, token string) error {
				if label != "test-credential" {
					t.Errorf("Ожидалась метка 'test-credential', получена '%s'", label)
				}
				if token != "test-token" {
					t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
				}
				return nil
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteCredential("test-credential")
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
	})

	// Тест ошибки при пустой метке
	t.Run("EmptyLabel", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteCredential("")
		if err == nil {
			t.Error("Ожидалась ошибка пустой метки, но ее не было")
		}
	})

	// Тест ошибки при загрузке токена
	t.Run("TokenLoadError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "", errors.New("ошибка загрузки токена")
			},
		}
		mockClientService := &MockClientServiceFixed{}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteCredential("test-credential")
		if err == nil {
			t.Error("Ожидалась ошибка загрузки токена, но ее не было")
		}
	})

	// Тест ошибки при удалении учетных данных
	t.Run("DeleteCredentialError", func(t *testing.T) {
		mockTokenService := &MockTokenServiceFixed{
			LoadTokenFunc: func() (string, error) {
				return "test-token", nil
			},
		}
		mockClientService := &MockClientServiceFixed{
			DeleteCredentialFunc: func(label string, token string) error {
				return errors.New("ошибка удаления учетных данных")
			},
		}

		clientUseCase := NewClientUseCase(mockTokenService, mockClientService)
		err := clientUseCase.DeleteCredential("test-credential")
		if err == nil {
			t.Error("Ожидалась ошибка удаления учетных данных, но ее не было")
		}
	})
}
