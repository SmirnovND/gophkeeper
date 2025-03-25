package service

import (
	"encoding/json"
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"testing"
	"time"
)

// TestNewDataService тестирует функцию NewDataService
func TestNewDataService(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{}
	mockUserDataRepo := &MockUserDataRepo{}

	// Вызываем функцию NewDataService
	dataService := NewDataService(mockUserDataRepo, mockUserRepo)

	// Проверяем, что возвращенный объект не nil
	if dataService == nil {
		t.Fatal("Функция NewDataService вернула nil")
	}

	// Проверяем, что возвращенный объект реализует интерфейс DataService
	_, ok := dataService.(interfaces.DataService)
	if !ok {
		t.Fatal("Возвращенный объект не реализует интерфейс DataService")
	}
}

// TestDataService_SaveFileMetadata тестирует метод SaveFileMetadata
func TestDataService_SaveFileMetadata(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	mockUserDataRepo := &MockUserDataRepo{
		SaveUserDataFunc: func(userData *domain.UserData) error {
			// Проверяем параметры
			if userData.UserID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userData.UserID)
			}
			if userData.Label != "test-file" {
				t.Errorf("Ожидалась метка 'test-file', получена '%s'", userData.Label)
			}
			if userData.Type != domain.UserDataTypeFile {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeFile, userData.Type)
			}

			// Проверяем данные
			var fileMetadata domain.FileMetadata
			if err := json.Unmarshal(userData.Data, &fileMetadata); err != nil {
				t.Fatalf("Ошибка при десериализации данных: %v", err)
			}
			if fileMetadata.FileName != "test-file" {
				t.Errorf("Ожидалось имя файла 'test-file', получено '%s'", fileMetadata.FileName)
			}
			if fileMetadata.Extension != "txt" {
				t.Errorf("Ожидалось расширение 'txt', получено '%s'", fileMetadata.Extension)
			}

			return nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод SaveFileMetadata
	fileData := &domain.FileData{
		Name:      "test-file",
		Extension: "txt",
	}
	err := dataService.SaveFileMetadata("testuser", "test-file", fileData)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveFileMetadata: %v", err)
	}
}

// TestDataService_SaveFileMetadata_UserNotFound тестирует метод SaveFileMetadata с ошибкой "пользователь не найден"
func TestDataService_SaveFileMetadata_UserNotFound(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			return nil, errors.New("пользователь не найден")
		},
	}

	mockUserDataRepo := &MockUserDataRepo{}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод SaveFileMetadata
	fileData := &domain.FileData{
		Name:      "test-file",
		Extension: "txt",
	}
	err := dataService.SaveFileMetadata("testuser", "test-file", fileData)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestDataService_GetFileMetadata тестирует метод GetFileMetadata
func TestDataService_GetFileMetadata(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	// Создаем метаданные файла
	fileMetadata := domain.FileMetadata{
		FileName:  "test-file",
		Extension: "txt",
	}
	fileMetadataJSON, _ := json.Marshal(fileMetadata)

	mockUserDataRepo := &MockUserDataRepo{
		GetUserDataByLabelAndTypeFunc: func(userID, label string, dataType string) (*domain.UserData, error) {
			// Проверяем параметры
			if userID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userID)
			}
			if label != "test-file" {
				t.Errorf("Ожидалась метка 'test-file', получена '%s'", label)
			}
			if dataType != domain.UserDataTypeFile {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeFile, dataType)
			}

			return &domain.UserData{
				ID:        "data123",
				UserID:    "user123",
				Label:     "test-file",
				Type:      domain.UserDataTypeFile,
				Data:      fileMetadataJSON,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод GetFileMetadata
	result, err := dataService.GetFileMetadata("testuser", "test-file")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове GetFileMetadata: %v", err)
	}
	if result == nil {
		t.Fatal("Результат не должен быть nil")
	}
	if result.FileName != "test-file" {
		t.Errorf("Ожидалось имя файла 'test-file', получено '%s'", result.FileName)
	}
	if result.Extension != "txt" {
		t.Errorf("Ожидалось расширение 'txt', получено '%s'", result.Extension)
	}
}

// TestDataService_GetFileMetadata_UserNotFound тестирует метод GetFileMetadata с ошибкой "пользователь не найден"
func TestDataService_GetFileMetadata_UserNotFound(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			return nil, errors.New("пользователь не найден")
		},
	}

	mockUserDataRepo := &MockUserDataRepo{}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод GetFileMetadata
	_, err := dataService.GetFileMetadata("testuser", "test-file")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestDataService_GetFileMetadata_DataNotFound тестирует метод GetFileMetadata с ошибкой "данные не найдены"
func TestDataService_GetFileMetadata_DataNotFound(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	mockUserDataRepo := &MockUserDataRepo{
		GetUserDataByLabelAndTypeFunc: func(userID, label string, dataType string) (*domain.UserData, error) {
			return nil, domain.ErrNotFound
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод GetFileMetadata
	_, err := dataService.GetFileMetadata("testuser", "test-file")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestDataService_DeleteFileMetadata тестирует метод DeleteFileMetadata
func TestDataService_DeleteFileMetadata(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	mockUserDataRepo := &MockUserDataRepo{
		GetUserDataByLabelAndTypeFunc: func(userID, label string, dataType string) (*domain.UserData, error) {
			// Проверяем параметры
			if userID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userID)
			}
			if label != "test-file" {
				t.Errorf("Ожидалась метка 'test-file', получена '%s'", label)
			}
			if dataType != domain.UserDataTypeFile {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeFile, dataType)
			}

			return &domain.UserData{
				ID:     "data123",
				UserID: "user123",
				Label:  "test-file",
				Type:   domain.UserDataTypeFile,
			}, nil
		},
		DeleteUserDataFunc: func(id string) error {
			// Проверяем параметры
			if id != "data123" {
				t.Errorf("Ожидался ID 'data123', получен '%s'", id)
			}
			return nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод DeleteFileMetadata
	err := dataService.DeleteFileMetadata("testuser", "test-file")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове DeleteFileMetadata: %v", err)
	}
}

// TestDataService_SaveCredential тестирует метод SaveCredential
func TestDataService_SaveCredential(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	mockUserDataRepo := &MockUserDataRepo{
		SaveUserDataFunc: func(userData *domain.UserData) error {
			// Проверяем параметры
			if userData.UserID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userData.UserID)
			}
			if userData.Label != "test-credential" {
				t.Errorf("Ожидалась метка 'test-credential', получена '%s'", userData.Label)
			}
			if userData.Type != domain.UserDataTypeCredential {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeCredential, userData.Type)
			}

			// Проверяем данные
			var credentialData domain.CredentialData
			if err := json.Unmarshal(userData.Data, &credentialData); err != nil {
				t.Fatalf("Ошибка при десериализации данных: %v", err)
			}
			if credentialData.Login != "service-login" {
				t.Errorf("Ожидался логин 'service-login', получен '%s'", credentialData.Login)
			}
			if credentialData.Password != "service-password" {
				t.Errorf("Ожидался пароль 'service-password', получен '%s'", credentialData.Password)
			}

			return nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод SaveCredential
	credentialData := &domain.CredentialData{
		Login:    "service-login",
		Password: "service-password",
	}
	err := dataService.SaveCredential("testuser", "test-credential", credentialData)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveCredential: %v", err)
	}
}

// TestDataService_GetCredential тестирует метод GetCredential
func TestDataService_GetCredential(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	// Создаем данные учетной записи
	credentialData := domain.CredentialData{
		Login:    "service-login",
		Password: "service-password",
	}
	credentialDataJSON, _ := json.Marshal(credentialData)

	mockUserDataRepo := &MockUserDataRepo{
		GetUserDataByLabelAndTypeFunc: func(userID, label string, dataType string) (*domain.UserData, error) {
			// Проверяем параметры
			if userID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userID)
			}
			if label != "test-credential" {
				t.Errorf("Ожидалась метка 'test-credential', получена '%s'", label)
			}
			if dataType != domain.UserDataTypeCredential {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeCredential, dataType)
			}

			return &domain.UserData{
				ID:        "data123",
				UserID:    "user123",
				Label:     "test-credential",
				Type:      domain.UserDataTypeCredential,
				Data:      credentialDataJSON,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод GetCredential
	result, err := dataService.GetCredential("testuser", "test-credential")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове GetCredential: %v", err)
	}
	if result == nil {
		t.Fatal("Результат не должен быть nil")
	}
	if result.Login != "service-login" {
		t.Errorf("Ожидался логин 'service-login', получен '%s'", result.Login)
	}
	if result.Password != "service-password" {
		t.Errorf("Ожидался пароль 'service-password', получен '%s'", result.Password)
	}
}

// TestDataService_DeleteCredential тестирует метод DeleteCredential
func TestDataService_DeleteCredential(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	mockUserDataRepo := &MockUserDataRepo{
		GetUserDataByLabelAndTypeFunc: func(userID, label string, dataType string) (*domain.UserData, error) {
			// Проверяем параметры
			if userID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userID)
			}
			if label != "test-credential" {
				t.Errorf("Ожидалась метка 'test-credential', получена '%s'", label)
			}
			if dataType != domain.UserDataTypeCredential {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeCredential, dataType)
			}

			return &domain.UserData{
				ID:     "data123",
				UserID: "user123",
				Label:  "test-credential",
				Type:   domain.UserDataTypeCredential,
			}, nil
		},
		DeleteUserDataFunc: func(id string) error {
			// Проверяем параметры
			if id != "data123" {
				t.Errorf("Ожидался ID 'data123', получен '%s'", id)
			}
			return nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод DeleteCredential
	err := dataService.DeleteCredential("testuser", "test-credential")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове DeleteCredential: %v", err)
	}
}

// TestDataService_SaveCard тестирует метод SaveCard
func TestDataService_SaveCard(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	mockUserDataRepo := &MockUserDataRepo{
		SaveUserDataFunc: func(userData *domain.UserData) error {
			// Проверяем параметры
			if userData.UserID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userData.UserID)
			}
			if userData.Label != "test-card" {
				t.Errorf("Ожидалась метка 'test-card', получена '%s'", userData.Label)
			}
			if userData.Type != domain.UserDataTypeCard {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeCard, userData.Type)
			}

			// Проверяем данные
			var cardData domain.CardData
			if err := json.Unmarshal(userData.Data, &cardData); err != nil {
				t.Fatalf("Ошибка при десериализации данных: %v", err)
			}
			if cardData.Number != "1234567890123456" {
				t.Errorf("Ожидался номер карты '1234567890123456', получен '%s'", cardData.Number)
			}
			if cardData.Holder != "Test User" {
				t.Errorf("Ожидался держатель карты 'Test User', получен '%s'", cardData.Holder)
			}
			if cardData.ExpiryDate != "12/25" {
				t.Errorf("Ожидался срок действия '12/25', получен '%s'", cardData.ExpiryDate)
			}
			if cardData.CVV != "123" {
				t.Errorf("Ожидался CVV '123', получен '%s'", cardData.CVV)
			}

			return nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод SaveCard
	cardData := &domain.CardData{
		Number:     "1234567890123456",
		Holder:     "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}
	err := dataService.SaveCard("testuser", "test-card", cardData)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveCard: %v", err)
	}
}

// TestDataService_GetCard тестирует метод GetCard
func TestDataService_GetCard(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	// Создаем данные карты
	cardData := domain.CardData{
		Number:     "1234567890123456",
		Holder:     "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}
	cardDataJSON, _ := json.Marshal(cardData)

	mockUserDataRepo := &MockUserDataRepo{
		GetUserDataByLabelAndTypeFunc: func(userID, label string, dataType string) (*domain.UserData, error) {
			// Проверяем параметры
			if userID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userID)
			}
			if label != "test-card" {
				t.Errorf("Ожидалась метка 'test-card', получена '%s'", label)
			}
			if dataType != domain.UserDataTypeCard {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeCard, dataType)
			}

			return &domain.UserData{
				ID:        "data123",
				UserID:    "user123",
				Label:     "test-card",
				Type:      domain.UserDataTypeCard,
				Data:      cardDataJSON,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод GetCard
	result, err := dataService.GetCard("testuser", "test-card")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове GetCard: %v", err)
	}
	if result == nil {
		t.Fatal("Результат не должен быть nil")
	}
	if result.Number != "1234567890123456" {
		t.Errorf("Ожидался номер карты '1234567890123456', получен '%s'", result.Number)
	}
	if result.Holder != "Test User" {
		t.Errorf("Ожидался держатель карты 'Test User', получен '%s'", result.Holder)
	}
	if result.ExpiryDate != "12/25" {
		t.Errorf("Ожидался срок действия '12/25', получен '%s'", result.ExpiryDate)
	}
	if result.CVV != "123" {
		t.Errorf("Ожидался CVV '123', получен '%s'", result.CVV)
	}
}

// TestDataService_DeleteCard тестирует метод DeleteCard
func TestDataService_DeleteCard(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	mockUserDataRepo := &MockUserDataRepo{
		GetUserDataByLabelAndTypeFunc: func(userID, label string, dataType string) (*domain.UserData, error) {
			// Проверяем параметры
			if userID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userID)
			}
			if label != "test-card" {
				t.Errorf("Ожидалась метка 'test-card', получена '%s'", label)
			}
			if dataType != domain.UserDataTypeCard {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeCard, dataType)
			}

			return &domain.UserData{
				ID:     "data123",
				UserID: "user123",
				Label:  "test-card",
				Type:   domain.UserDataTypeCard,
			}, nil
		},
		DeleteUserDataFunc: func(id string) error {
			// Проверяем параметры
			if id != "data123" {
				t.Errorf("Ожидался ID 'data123', получен '%s'", id)
			}
			return nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод DeleteCard
	err := dataService.DeleteCard("testuser", "test-card")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове DeleteCard: %v", err)
	}
}

// TestDataService_SaveText тестирует метод SaveText
func TestDataService_SaveText(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	mockUserDataRepo := &MockUserDataRepo{
		SaveUserDataFunc: func(userData *domain.UserData) error {
			// Проверяем параметры
			if userData.UserID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userData.UserID)
			}
			if userData.Label != "test-text" {
				t.Errorf("Ожидалась метка 'test-text', получена '%s'", userData.Label)
			}
			if userData.Type != domain.UserDataTypeText {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeText, userData.Type)
			}

			// Проверяем данные
			var textData domain.TextData
			if err := json.Unmarshal(userData.Data, &textData); err != nil {
				t.Fatalf("Ошибка при десериализации данных: %v", err)
			}
			if textData.Content != "test text content" {
				t.Errorf("Ожидался текст 'test text content', получен '%s'", textData.Content)
			}

			return nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод SaveText
	textData := &domain.TextData{
		Content: "test text content",
	}
	err := dataService.SaveText("testuser", "test-text", textData)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveText: %v", err)
	}
}

// TestDataService_GetText тестирует метод GetText
func TestDataService_GetText(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	// Создаем текстовые данные
	textData := domain.TextData{
		Content: "test text content",
	}
	textDataJSON, _ := json.Marshal(textData)

	mockUserDataRepo := &MockUserDataRepo{
		GetUserDataByLabelAndTypeFunc: func(userID, label string, dataType string) (*domain.UserData, error) {
			// Проверяем параметры
			if userID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userID)
			}
			if label != "test-text" {
				t.Errorf("Ожидалась метка 'test-text', получена '%s'", label)
			}
			if dataType != domain.UserDataTypeText {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeText, dataType)
			}

			return &domain.UserData{
				ID:        "data123",
				UserID:    "user123",
				Label:     "test-text",
				Type:      domain.UserDataTypeText,
				Data:      textDataJSON,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод GetText
	result, err := dataService.GetText("testuser", "test-text")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове GetText: %v", err)
	}
	if result == nil {
		t.Fatal("Результат не должен быть nil")
	}
	if result.Content != "test text content" {
		t.Errorf("Ожидался текст 'test text content', получен '%s'", result.Content)
	}
}

// TestDataService_DeleteText тестирует метод DeleteText
func TestDataService_DeleteText(t *testing.T) {
	// Создаем моки для репозиториев
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	mockUserDataRepo := &MockUserDataRepo{
		GetUserDataByLabelAndTypeFunc: func(userID, label string, dataType string) (*domain.UserData, error) {
			// Проверяем параметры
			if userID != "user123" {
				t.Errorf("Ожидался UserID 'user123', получен '%s'", userID)
			}
			if label != "test-text" {
				t.Errorf("Ожидалась метка 'test-text', получена '%s'", label)
			}
			if dataType != domain.UserDataTypeText {
				t.Errorf("Ожидался тип '%s', получен '%s'", domain.UserDataTypeText, dataType)
			}

			return &domain.UserData{
				ID:     "data123",
				UserID: "user123",
				Label:  "test-text",
				Type:   domain.UserDataTypeText,
			}, nil
		},
		DeleteUserDataFunc: func(id string) error {
			// Проверяем параметры
			if id != "data123" {
				t.Errorf("Ожидался ID 'data123', получен '%s'", id)
			}
			return nil
		},
	}

	// Создаем экземпляр DataService
	dataService := &DataService{
		repo:     mockUserDataRepo,
		userRepo: mockUserRepo,
	}

	// Вызываем метод DeleteText
	err := dataService.DeleteText("testuser", "test-text")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове DeleteText: %v", err)
	}
}
