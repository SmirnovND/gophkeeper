package usecase

import (
	"encoding/json"
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockCloudService - мок для CloudService
type MockCloudService struct {
	GenerateUploadLinkFunc   func(fileName string) (string, error)
	GenerateDownloadLinkFunc func(fileName string) (string, error)
}

func (m *MockCloudService) GenerateUploadLink(fileName string) (string, error) {
	if m.GenerateUploadLinkFunc != nil {
		return m.GenerateUploadLinkFunc(fileName)
	}
	return "", nil
}

func (m *MockCloudService) GenerateDownloadLink(fileName string) (string, error) {
	if m.GenerateDownloadLinkFunc != nil {
		return m.GenerateDownloadLinkFunc(fileName)
	}
	return "", nil
}

// MockDataServiceCloud - мок для DataService
type MockDataServiceCloud struct {
	SaveFileMetadataFunc func(login string, label string, fileData *domain.FileData) error
	GetFileMetadataFunc  func(login string, label string) (*domain.FileMetadata, error)
}

func (m *MockDataServiceCloud) SaveFileMetadata(login string, label string, fileData *domain.FileData) error {
	if m.SaveFileMetadataFunc != nil {
		return m.SaveFileMetadataFunc(login, label, fileData)
	}
	return nil
}

func (m *MockDataServiceCloud) GetFileMetadata(login string, label string) (*domain.FileMetadata, error) {
	if m.GetFileMetadataFunc != nil {
		return m.GetFileMetadataFunc(login, label)
	}
	return nil, nil
}

// Заглушки для остальных методов интерфейса DataService
func (m *MockDataServiceCloud) DeleteFileMetadata(login string, label string) error {
	return nil
}

func (m *MockDataServiceCloud) SaveCredential(login string, label string, credentialData *domain.CredentialData) error {
	return nil
}

func (m *MockDataServiceCloud) GetCredential(login string, label string) (*domain.CredentialData, error) {
	return nil, nil
}

func (m *MockDataServiceCloud) DeleteCredential(login string, label string) error {
	return nil
}

func (m *MockDataServiceCloud) SaveCard(login string, label string, cardData *domain.CardData) error {
	return nil
}

func (m *MockDataServiceCloud) GetCard(login string, label string) (*domain.CardData, error) {
	return nil, nil
}

func (m *MockDataServiceCloud) DeleteCard(login string, label string) error {
	return nil
}

func (m *MockDataServiceCloud) SaveText(login string, label string, textData *domain.TextData) error {
	return nil
}

func (m *MockDataServiceCloud) GetText(login string, label string) (*domain.TextData, error) {
	return nil, nil
}

func (m *MockDataServiceCloud) DeleteText(login string, label string) error {
	return nil
}

// TestNewCloudUseCase проверяет создание нового экземпляра CloudUseCase
func TestNewCloudUseCase(t *testing.T) {
	mockCloudService := &MockCloudService{}
	MockDataServiceCloud := &MockDataServiceCloud{}
	MockJwtService := &MockJwtService{}
	cloudUseCase := NewCloudUseCase(mockCloudService, MockDataServiceCloud, MockJwtService)

	if cloudUseCase == nil {
		t.Fatal("NewCloudUseCase вернул nil")
	}
}

// TestCloudUseCase_GenerateUploadLink_Success проверяет успешную генерацию ссылки для загрузки
func TestCloudUseCase_GenerateUploadLink_Success(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{
		GenerateUploadLinkFunc: func(fileName string) (string, error) {
			expectedFileName := "testuser_test-file.txt"
			if fileName != expectedFileName {
				t.Errorf("Ожидалось имя файла '%s', получено '%s'", expectedFileName, fileName)
			}
			return "https://example.com/upload/test-file.txt", nil
		},
	}

	MockDataServiceCloud := &MockDataServiceCloud{
		SaveFileMetadataFunc: func(login string, label string, fileData *domain.FileData) error {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-file" {
				t.Errorf("Ожидалась метка 'test-file', получена '%s'", label)
			}
			if fileData.Name != "test-file" {
				t.Errorf("Ожидалось имя файла 'test-file', получено '%s'", fileData.Name)
			}
			if fileData.Extension != "txt" {
				t.Errorf("Ожидалось расширение 'txt', получено '%s'", fileData.Extension)
			}
			return nil
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/files/upload", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Создаем данные файла
	fileData := &domain.FileData{
		Name:      "test-file",
		Extension: "txt",
	}

	// Вызываем метод GenerateUploadLink
	cloudUseCase.GenerateUploadLink(w, req, fileData)

	// Проверяем статус ответа
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}

	// Проверяем тип содержимого
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Ожидался Content-Type 'application/json', получен '%s'", contentType)
	}

	// Проверяем тело ответа
	var response domain.FileDataResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	expectedURL := "https://example.com/upload/test-file.txt"
	if response.Url != expectedURL {
		t.Errorf("Ожидался URL '%s', получен '%s'", expectedURL, response.Url)
	}

	expectedDescription := "Загрузи файл по этой ссылке"
	if response.Description != expectedDescription {
		t.Errorf("Ожидалось описание '%s', получено '%s'", expectedDescription, response.Description)
	}
}

// TestCloudUseCase_GenerateUploadLink_TokenError проверяет обработку ошибки при извлечении логина из токена
func TestCloudUseCase_GenerateUploadLink_TokenError(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{}
	MockDataServiceCloud := &MockDataServiceCloud{}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "", errors.New("ошибка извлечения логина из токена")
		},
	}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/files/upload", nil)
	req.Header.Set("Authorization", "Bearer error-token")
	w := httptest.NewRecorder()

	// Создаем данные файла
	fileData := &domain.FileData{
		Name:      "test-file",
		Extension: "txt",
	}

	// Вызываем метод GenerateUploadLink
	cloudUseCase.GenerateUploadLink(w, req, fileData)

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка получения логина") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка получения логина', получено '%s'", w.Body.String())
	}
}

// TestCloudUseCase_GenerateUploadLink_InvalidFileData проверяет обработку невалидных данных файла
func TestCloudUseCase_GenerateUploadLink_InvalidFileData(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{}
	MockDataServiceCloud := &MockDataServiceCloud{}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/files/upload", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Тест с nil fileData
	cloudUseCase.GenerateUploadLink(w, req, nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверные данные файла")

	// Сбрасываем recorder
	w = httptest.NewRecorder()

	// Тест с пустым именем файла
	fileData := &domain.FileData{
		Name:      "",
		Extension: "txt",
	}
	cloudUseCase.GenerateUploadLink(w, req, fileData)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверные данные файла")

	// Сбрасываем recorder
	w = httptest.NewRecorder()

	// Тест с пустым расширением файла
	fileData = &domain.FileData{
		Name:      "test-file",
		Extension: "",
	}
	cloudUseCase.GenerateUploadLink(w, req, fileData)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Неверные данные файла")
}

// TestCloudUseCase_GenerateUploadLink_CloudServiceError проверяет обработку ошибки при генерации ссылки
func TestCloudUseCase_GenerateUploadLink_CloudServiceError(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{
		GenerateUploadLinkFunc: func(fileName string) (string, error) {
			return "", errors.New("ошибка генерации ссылки")
		},
	}
	MockDataServiceCloud := &MockDataServiceCloud{}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/files/upload", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Создаем данные файла
	fileData := &domain.FileData{
		Name:      "test-file",
		Extension: "txt",
	}

	// Вызываем метод GenerateUploadLink
	cloudUseCase.GenerateUploadLink(w, req, fileData)

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка при генерации ссылки") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка при генерации ссылки', получено '%s'", w.Body.String())
	}
}

// TestCloudUseCase_GenerateUploadLink_DataServiceError проверяет обработку ошибки при сохранении метаданных
func TestCloudUseCase_GenerateUploadLink_DataServiceError(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{
		GenerateUploadLinkFunc: func(fileName string) (string, error) {
			return "https://example.com/upload/test-file.txt", nil
		},
	}
	MockDataServiceCloud := &MockDataServiceCloud{
		SaveFileMetadataFunc: func(login string, label string, fileData *domain.FileData) error {
			return errors.New("ошибка сохранения метаданных")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/files/upload", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Создаем данные файла
	fileData := &domain.FileData{
		Name:      "test-file",
		Extension: "txt",
	}

	// Вызываем метод GenerateUploadLink
	cloudUseCase.GenerateUploadLink(w, req, fileData)

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка при сохранении метаданных файла") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка при сохранении метаданных файла', получено '%s'", w.Body.String())
	}
}

// TestCloudUseCase_GenerateDownloadLink_Success проверяет успешную генерацию ссылки для скачивания
func TestCloudUseCase_GenerateDownloadLink_Success(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{
		GenerateDownloadLinkFunc: func(fileName string) (string, error) {
			expectedFileName := "testuser_test-file.txt"
			if fileName != expectedFileName {
				t.Errorf("Ожидалось имя файла '%s', получено '%s'", expectedFileName, fileName)
			}
			return "https://example.com/download/test-file.txt", nil
		},
	}

	MockDataServiceCloud := &MockDataServiceCloud{
		GetFileMetadataFunc: func(login string, label string) (*domain.FileMetadata, error) {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-file" {
				t.Errorf("Ожидалась метка 'test-file', получена '%s'", label)
			}
			return &domain.FileMetadata{
				FileName:  "test-file",
				Extension: "txt",
			}, nil
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/files/download?label=test-file", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GenerateDownloadLink
	cloudUseCase.GenerateDownloadLink(w, req, "test-file")

	// Проверяем статус ответа
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}

	// Проверяем тип содержимого
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Ожидался Content-Type 'application/json', получен '%s'", contentType)
	}

	// Проверяем тело ответа
	var response struct {
		URL         string              `json:"url"`
		Description string              `json:"description"`
		Metadata    domain.FileMetadata `json:"metadata"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	expectedURL := "https://example.com/download/test-file.txt"
	if response.URL != expectedURL {
		t.Errorf("Ожидался URL '%s', получен '%s'", expectedURL, response.URL)
	}

	expectedDescription := "Скачай файл по этой ссылке"
	if response.Description != expectedDescription {
		t.Errorf("Ожидалось описание '%s', получено '%s'", expectedDescription, response.Description)
	}

	if response.Metadata.FileName != "test-file" {
		t.Errorf("Ожидалось имя файла 'test-file', получено '%s'", response.Metadata.FileName)
	}

	if response.Metadata.Extension != "txt" {
		t.Errorf("Ожидалось расширение 'txt', получено '%s'", response.Metadata.Extension)
	}
}

// TestCloudUseCase_GenerateDownloadLink_TokenError проверяет обработку ошибки при извлечении логина из токена
func TestCloudUseCase_GenerateDownloadLink_TokenError(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{}
	MockDataServiceCloud := &MockDataServiceCloud{}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "", errors.New("ошибка извлечения логина из токена")
		},
	}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/files/download?label=test-file", nil)
	req.Header.Set("Authorization", "Bearer error-token")
	w := httptest.NewRecorder()

	// Вызываем метод GenerateDownloadLink
	cloudUseCase.GenerateDownloadLink(w, req, "test-file")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка получения логина") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка получения логина', получено '%s'", w.Body.String())
	}
}

// TestCloudUseCase_GenerateDownloadLink_EmptyLabel проверяет обработку пустой метки файла
func TestCloudUseCase_GenerateDownloadLink_EmptyLabel(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{}
	MockDataServiceCloud := &MockDataServiceCloud{}
	MockJwtService := &MockJwtService{}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   MockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/files/download", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GenerateDownloadLink с пустой меткой
	cloudUseCase.GenerateDownloadLink(w, req, "")

	// Проверяем статус ответа
	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Не указана метка файла") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Не указана метка файла', получено '%s'", w.Body.String())
	}
}

// TestCloudUseCase_GenerateDownloadLink_GetFileMetadataError проверяет обработку ошибки при получении метаданных файла
func TestCloudUseCase_GenerateDownloadLink_GetFileMetadataError(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{}
	mockJwtService := &MockJwtService{}
	MockDataServiceCloud := &MockDataServiceCloud{
		GetFileMetadataFunc: func(login string, label string) (*domain.FileMetadata, error) {
			return nil, errors.New("ошибка получения метаданных файла")
		},
	}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/files/download?label=test-file", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GenerateDownloadLink
	cloudUseCase.GenerateDownloadLink(w, req, "test-file")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка при получении метаданных файла") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка при получении метаданных файла', получено '%s'", w.Body.String())
	}
}

// TestCloudUseCase_GenerateDownloadLink_GenerateDownloadLinkError проверяет обработку ошибки при генерации ссылки для скачивания
func TestCloudUseCase_GenerateDownloadLink_GenerateDownloadLinkError(t *testing.T) {
	// Создаем моки для сервисов
	mockCloudService := &MockCloudService{
		GenerateDownloadLinkFunc: func(fileName string) (string, error) {
			return "", errors.New("ошибка генерации ссылки для скачивания")
		},
	}
	mockJwtService := &MockJwtService{}
	MockDataServiceCloud := &MockDataServiceCloud{
		GetFileMetadataFunc: func(login string, label string) (*domain.FileMetadata, error) {
			return &domain.FileMetadata{
				FileName:  "test-file",
				Extension: "txt",
			}, nil
		},
	}

	// Создаем экземпляр CloudUseCase
	cloudUseCase := &CloudUseCase{
		cloudService: mockCloudService,
		dataService:  MockDataServiceCloud,
		jwtService:   mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/files/download?label=test-file", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GenerateDownloadLink
	cloudUseCase.GenerateDownloadLink(w, req, "test-file")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка при генерации ссылки для скачивания") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка при генерации ссылки для скачивания', получено '%s'", w.Body.String())
	}
}
