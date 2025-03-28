package usecase

import (
	"encoding/json"
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Создаем мок для JWTService для тестов DataUseCase
type MockJWTServiceForDataUseCase struct {
	mock.Mock
}

func (m *MockJWTServiceForDataUseCase) ExtractLoginFromToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

// Создаем мок для DataService для тестов DataUseCase
type MockDataServiceForDataUseCase struct {
	mock.Mock
}

func (m *MockDataServiceForDataUseCase) SaveCard(login string, label string, cardData *domain.CardData, metadata string) error {
	args := m.Called(login, label, cardData, metadata)
	return args.Error(0)
}

func (m *MockDataServiceForDataUseCase) GetCard(login string, label string) (*domain.CardData, string, error) {
	args := m.Called(login, label)
	var cardData *domain.CardData
	if args.Get(0) != nil {
		cardData = args.Get(0).(*domain.CardData)
	}
	return cardData, args.String(1), args.Error(2)
}

func (m *MockDataServiceForDataUseCase) DeleteCard(login string, label string) error {
	args := m.Called(login, label)
	return args.Error(0)
}

func (m *MockDataServiceForDataUseCase) SaveText(login string, label string, textData *domain.TextData, metadata string) error {
	args := m.Called(login, label, textData, metadata)
	return args.Error(0)
}

func (m *MockDataServiceForDataUseCase) GetText(login string, label string) (*domain.TextData, string, error) {
	args := m.Called(login, label)
	var textData *domain.TextData
	if args.Get(0) != nil {
		textData = args.Get(0).(*domain.TextData)
	}
	return textData, args.String(1), args.Error(2)
}

func (m *MockDataServiceForDataUseCase) DeleteText(login string, label string) error {
	args := m.Called(login, label)
	return args.Error(0)
}

func (m *MockDataServiceForDataUseCase) SaveCredential(login string, label string, credentialData *domain.CredentialData, metadata string) error {
	args := m.Called(login, label, credentialData, metadata)
	return args.Error(0)
}

func (m *MockDataServiceForDataUseCase) GetCredential(login string, label string) (*domain.CredentialData, string, error) {
	args := m.Called(login, label)
	var credentialData *domain.CredentialData
	if args.Get(0) != nil {
		credentialData = args.Get(0).(*domain.CredentialData)
	}
	return credentialData, args.String(1), args.Error(2)
}

func (m *MockDataServiceForDataUseCase) DeleteCredential(login string, label string) error {
	args := m.Called(login, label)
	return args.Error(0)
}

func (m *MockDataServiceForDataUseCase) SaveFileMetadata(login string, label string, fileData *domain.FileData, metadata string) error {
	args := m.Called(login, label, fileData, metadata)
	return args.Error(0)
}

func (m *MockDataServiceForDataUseCase) GetFileMetadata(login string, label string) (*domain.FileMetadata, string, error) {
	args := m.Called(login, label)
	var fileMetadata *domain.FileMetadata
	if args.Get(0) != nil {
		fileMetadata = args.Get(0).(*domain.FileMetadata)
	}
	return fileMetadata, args.String(1), args.Error(2)
}

func (m *MockDataServiceForDataUseCase) DeleteFileMetadata(login string, label string) error {
	args := m.Called(login, label)
	return args.Error(0)
}

// TestDataUseCase_SaveCard тестирует метод SaveCard
func TestDataUseCase_SaveCard(t *testing.T) {
	// Тест успешного сохранения карты
	t.Run("Success", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("SaveCard", "testuser", "test-card", mock.AnythingOfType("*domain.CardData"), "test metadata").Return(nil)

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/data/card/test-card", nil)
		r.Header.Set("Authorization", "Bearer token123")
		cardData := &domain.CardData{
			Number:     "1234567890123456",
			Holder:     "Test User",
			ExpiryDate: "12/25",
			CVV:        "123",
		}

		// Вызываем метод SaveCard
		dataUseCase.SaveCard(w, r, "test-card", cardData, "test metadata")

		// Проверяем результаты
		assert.Equal(t, http.StatusOK, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)
	})

	// Тест ошибки при извлечении логина из токена
	t.Run("ExtractLoginError", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("", errors.New("ошибка извлечения логина"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/data/card/test-card", nil)
		r.Header.Set("Authorization", "Bearer token123")
		cardData := &domain.CardData{
			Number:     "1234567890123456",
			Holder:     "Test User",
			ExpiryDate: "12/25",
			CVV:        "123",
		}

		// Вызываем метод SaveCard
		dataUseCase.SaveCard(w, r, "test-card", cardData, "test metadata")

		// Проверяем результаты
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertNotCalled(t, "SaveCard")
	})

	// Тест ошибки при сохранении карты
	t.Run("SaveCardError", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("SaveCard", "testuser", "test-card", mock.AnythingOfType("*domain.CardData"), "test metadata").Return(errors.New("ошибка сохранения карты"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/data/card/test-card", nil)
		r.Header.Set("Authorization", "Bearer token123")
		cardData := &domain.CardData{
			Number:     "1234567890123456",
			Holder:     "Test User",
			ExpiryDate: "12/25",
			CVV:        "123",
		}

		// Вызываем метод SaveCard
		dataUseCase.SaveCard(w, r, "test-card", cardData, "test metadata")

		// Проверяем результаты
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)
	})
}

// TestDataUseCase_GetCard тестирует метод GetCard
func TestDataUseCase_GetCard(t *testing.T) {
	// Тест успешного получения карты
	t.Run("Success", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Создаем тестовые данные
		cardData := &domain.CardData{
			Number:     "1234567890123456",
			Holder:     "Test User",
			ExpiryDate: "12/25",
			CVV:        "123",
		}

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("GetCard", "testuser", "test-card").Return(cardData, "test metadata", nil)

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/data/card/test-card", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод GetCard
		dataUseCase.GetCard(w, r, "test-card")

		// Проверяем результаты
		assert.Equal(t, http.StatusOK, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)

		// Проверяем содержимое ответа
		var response struct {
			CardData *domain.CardData `json:"card_data"`
			Metadata string           `json:"metadata"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, cardData, response.CardData)
		assert.Equal(t, "test metadata", response.Metadata)
	})

	// Тест ошибки при извлечении логина из токена
	t.Run("ExtractLoginError", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("", errors.New("ошибка извлечения логина"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/data/card/test-card", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод GetCard
		dataUseCase.GetCard(w, r, "test-card")

		// Проверяем результаты
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertNotCalled(t, "GetCard")
	})

	// Тест ошибки "данные карты не найдены"
	t.Run("CardNotFound", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("GetCard", "testuser", "test-card").Return(nil, "", errors.New("данные карты не найдены"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/data/card/test-card", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод GetCard
		dataUseCase.GetCard(w, r, "test-card")

		// Проверяем результаты
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)
	})

	// Тест другой ошибки при получении карты
	t.Run("OtherError", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("GetCard", "testuser", "test-card").Return(nil, "", errors.New("другая ошибка"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/data/card/test-card", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод GetCard
		dataUseCase.GetCard(w, r, "test-card")

		// Проверяем результаты
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)
	})
}

// TestDataUseCase_GetText тестирует метод GetText
func TestDataUseCase_GetText(t *testing.T) {
	// Тест успешного получения текста
	t.Run("Success", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Создаем тестовые данные
		textData := &domain.TextData{
			Content: "Test content",
		}

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("GetText", "testuser", "test-text").Return(textData, "test metadata", nil)

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/data/text/test-text", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод GetText
		dataUseCase.GetText(w, r, "test-text")

		// Проверяем результаты
		assert.Equal(t, http.StatusOK, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)

		// Проверяем содержимое ответа
		var response struct {
			TextData *domain.TextData `json:"text_data"`
			Metadata string           `json:"metadata"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, textData, response.TextData)
		assert.Equal(t, "test metadata", response.Metadata)
	})

	// Тест ошибки при извлечении логина из токена
	t.Run("ExtractLoginError", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("", errors.New("ошибка извлечения логина"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/data/text/test-text", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод GetText
		dataUseCase.GetText(w, r, "test-text")

		// Проверяем результаты
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertNotCalled(t, "GetText")
	})

	// Тест ошибки "текстовые данные не найдены"
	t.Run("TextNotFound", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("GetText", "testuser", "test-text").Return(nil, "", errors.New("текстовые данные не найдены"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/data/text/test-text", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод GetText
		dataUseCase.GetText(w, r, "test-text")

		// Проверяем результаты
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)
	})

	// Тест другой ошибки при получении текста
	t.Run("OtherError", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("GetText", "testuser", "test-text").Return(nil, "", errors.New("другая ошибка"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/data/text/test-text", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод GetText
		dataUseCase.GetText(w, r, "test-text")

		// Проверяем результаты
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)
	})
}

// TestDataUseCase_DeleteCredential тестирует метод DeleteCredential
func TestDataUseCase_DeleteCredential(t *testing.T) {
	// Тест успешного удаления учетных данных
	t.Run("Success", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("DeleteCredential", "testuser", "test-credential").Return(nil)

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/api/data/credential/test-credential", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод DeleteCredential
		dataUseCase.DeleteCredential(w, r, "test-credential")

		// Проверяем результаты
		assert.Equal(t, http.StatusOK, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)
	})

	// Тест ошибки при извлечении логина из токена
	t.Run("ExtractLoginError", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("", errors.New("ошибка извлечения логина"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/api/data/credential/test-credential", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод DeleteCredential
		dataUseCase.DeleteCredential(w, r, "test-credential")

		// Проверяем результаты
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertNotCalled(t, "DeleteCredential")
	})

	// Тест ошибки "учетные данные не найдены"
	t.Run("CredentialNotFound", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("DeleteCredential", "testuser", "test-credential").Return(errors.New("учетные данные не найдены"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/api/data/credential/test-credential", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод DeleteCredential
		dataUseCase.DeleteCredential(w, r, "test-credential")

		// Проверяем результаты
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)
	})

	// Тест другой ошибки при удалении учетных данных
	t.Run("OtherError", func(t *testing.T) {
		// Создаем моки
		mockJWTService := new(MockJWTServiceForDataUseCase)
		mockDataService := new(MockDataServiceForDataUseCase)

		// Настраиваем поведение моков
		mockJWTService.On("ExtractLoginFromToken", "Bearer token123").Return("testuser", nil)
		mockDataService.On("DeleteCredential", "testuser", "test-credential").Return(errors.New("другая ошибка"))

		// Создаем экземпляр DataUseCase
		dataUseCase := &DataUseCase{
			jwtService:  mockJWTService,
			dataService: mockDataService,
		}

		// Создаем тестовые данные
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/api/data/credential/test-credential", nil)
		r.Header.Set("Authorization", "Bearer token123")

		// Вызываем метод DeleteCredential
		dataUseCase.DeleteCredential(w, r, "test-credential")

		// Проверяем результаты
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockJWTService.AssertExpectations(t)
		mockDataService.AssertExpectations(t)
	})
}