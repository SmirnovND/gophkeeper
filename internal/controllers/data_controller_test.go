package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Создаем мок для DataUseCase
type MockDataUseCase struct {
	mock.Mock
}

func (m *MockDataUseCase) SaveCredential(w http.ResponseWriter, r *http.Request, label string, credentialData *domain.CredentialData, metadata string) {
	m.Called(w, r, label, credentialData, metadata)
}

func (m *MockDataUseCase) GetCredential(w http.ResponseWriter, r *http.Request, label string) {
	m.Called(w, r, label)
}

func (m *MockDataUseCase) DeleteCredential(w http.ResponseWriter, r *http.Request, label string) {
	m.Called(w, r, label)
}

func (m *MockDataUseCase) SaveCard(w http.ResponseWriter, r *http.Request, label string, cardData *domain.CardData, metadata string) {
	m.Called(w, r, label, cardData, metadata)
}

func (m *MockDataUseCase) GetCard(w http.ResponseWriter, r *http.Request, label string) {
	m.Called(w, r, label)
}

func (m *MockDataUseCase) DeleteCard(w http.ResponseWriter, r *http.Request, label string) {
	m.Called(w, r, label)
}

func (m *MockDataUseCase) SaveText(w http.ResponseWriter, r *http.Request, label string, textData *domain.TextData, metadata string) {
	m.Called(w, r, label, textData, metadata)
}

func (m *MockDataUseCase) GetText(w http.ResponseWriter, r *http.Request, label string) {
	m.Called(w, r, label)
}

func (m *MockDataUseCase) DeleteText(w http.ResponseWriter, r *http.Request, label string) {
	m.Called(w, r, label)
}

// Вспомогательная функция для создания запроса с параметрами URL
func createRequestWithURLParam(method, path, paramName, paramValue string, body []byte) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаем контекст с параметрами URL
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(paramName, paramValue)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	return req, rr
}

// Тесты для методов работы с учетными данными

func TestDataController_SaveCredential(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"
	credentialData := domain.CredentialData{
		Login:    "testuser",
		Password: "password123",
	}
	metadata := "test metadata"

	// Создаем JSON из данных
	requestData := struct {
		CredentialData *domain.CredentialData `json:"credential_data"`
		Metadata       string                 `json:"metadata"`
	}{
		CredentialData: &credentialData,
		Metadata:       metadata,
	}
	jsonData, _ := json.Marshal(requestData)

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("POST", "/api/data/credential/"+label, "label", label, jsonData)

	// Настраиваем поведение мока
	mockDataUseCase.On("SaveCredential", mock.Anything, mock.Anything, label, mock.MatchedBy(func(c *domain.CredentialData) bool {
		return c.Login == credentialData.Login && c.Password == credentialData.Password
	}), metadata)

	// Act
	controller.SaveCredential(rr, req)

	// Assert
	mockDataUseCase.AssertExpectations(t)
}

func TestDataController_GetCredential(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("GET", "/api/data/credential/"+label, "label", label, nil)

	// Настраиваем поведение мока
	mockDataUseCase.On("GetCredential", mock.Anything, mock.Anything, label)

	// Act
	controller.GetCredential(rr, req)

	// Assert
	mockDataUseCase.AssertExpectations(t)
}

func TestDataController_DeleteCredential(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("DELETE", "/api/data/credential/"+label, "label", label, nil)

	// Настраиваем поведение мока
	mockDataUseCase.On("DeleteCredential", mock.Anything, mock.Anything, label)

	// Act
	controller.DeleteCredential(rr, req)

	// Assert
	mockDataUseCase.AssertExpectations(t)
}

// Тесты для методов работы с данными кредитных карт

func TestDataController_SaveCard(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"
	cardData := domain.CardData{
		Number:     "1234567890123456",
		Holder:     "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}
	metadata := "card metadata"

	// Создаем JSON из данных
	requestData := struct {
		CardData *domain.CardData `json:"card_data"`
		Metadata string           `json:"metadata"`
	}{
		CardData: &cardData,
		Metadata: metadata,
	}
	jsonData, _ := json.Marshal(requestData)

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("POST", "/api/data/card/"+label, "label", label, jsonData)

	// Настраиваем поведение мока
	mockDataUseCase.On("SaveCard", mock.Anything, mock.Anything, label, mock.MatchedBy(func(c *domain.CardData) bool {
		return c.Number == cardData.Number && c.Holder == cardData.Holder &&
			c.ExpiryDate == cardData.ExpiryDate && c.CVV == cardData.CVV
	}), metadata)

	// Act
	controller.SaveCard(rr, req)

	// Assert
	mockDataUseCase.AssertExpectations(t)
}

func TestDataController_GetCard(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("GET", "/api/data/card/"+label, "label", label, nil)

	// Настраиваем поведение мока
	mockDataUseCase.On("GetCard", mock.Anything, mock.Anything, label)

	// Act
	controller.GetCard(rr, req)

	// Assert
	mockDataUseCase.AssertExpectations(t)
}

func TestDataController_DeleteCard(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("DELETE", "/api/data/card/"+label, "label", label, nil)

	// Настраиваем поведение мока
	mockDataUseCase.On("DeleteCard", mock.Anything, mock.Anything, label)

	// Act
	controller.DeleteCard(rr, req)

	// Assert
	mockDataUseCase.AssertExpectations(t)
}

// Тесты для методов работы с текстовыми данными

func TestDataController_SaveText(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"
	textData := domain.TextData{
		Content: "This is a test text content",
	}
	metadata := "text metadata"

	// Создаем JSON из данных
	requestData := struct {
		TextData *domain.TextData `json:"text_data"`
		Metadata string           `json:"metadata"`
	}{
		TextData: &textData,
		Metadata: metadata,
	}
	jsonData, _ := json.Marshal(requestData)

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("POST", "/api/data/text/"+label, "label", label, jsonData)

	// Настраиваем поведение мока
	mockDataUseCase.On("SaveText", mock.Anything, mock.Anything, label, mock.MatchedBy(func(t *domain.TextData) bool {
		return t.Content == textData.Content
	}), metadata)

	// Act
	controller.SaveText(rr, req)

	// Assert
	mockDataUseCase.AssertExpectations(t)
}

func TestDataController_GetText(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("GET", "/api/data/text/"+label, "label", label, nil)

	// Настраиваем поведение мока
	mockDataUseCase.On("GetText", mock.Anything, mock.Anything, label)

	// Act
	controller.GetText(rr, req)

	// Assert
	mockDataUseCase.AssertExpectations(t)
}

func TestDataController_DeleteText(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("DELETE", "/api/data/text/"+label, "label", label, nil)

	// Настраиваем поведение мока
	mockDataUseCase.On("DeleteText", mock.Anything, mock.Anything, label)

	// Act
	controller.DeleteText(rr, req)

	// Assert
	mockDataUseCase.AssertExpectations(t)
}

// Тесты для проверки обработки ошибок

func TestDataController_SaveCredential_MissingLabel(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	credentialData := domain.CredentialData{
		Login:    "testuser",
		Password: "password123",
	}
	metadata := "test metadata"

	// Создаем JSON из данных
	requestData := struct {
		CredentialData *domain.CredentialData `json:"credential_data"`
		Metadata       string                 `json:"metadata"`
	}{
		CredentialData: &credentialData,
		Metadata:       metadata,
	}
	jsonData, _ := json.Marshal(requestData)

	// Создаем запрос без параметра label
	req, _ := http.NewRequest("POST", "/api/data/credential/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Act
	controller.SaveCredential(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockDataUseCase.AssertNotCalled(t, "SaveCredential")
}

func TestDataController_SaveCredential_InvalidJSON(t *testing.T) {
	// Arrange
	mockDataUseCase := new(MockDataUseCase)
	controller := NewDataController(mockDataUseCase)

	// Создаем тестовые данные
	label := "test-label"

	// Создаем некорректный JSON
	invalidJSON := []byte(`{"credential_data": {"login": "testuser", "password":}, "metadata": "test"}`)

	// Создаем запрос с параметрами URL
	req, rr := createRequestWithURLParam("POST", "/api/data/credential/"+label, "label", label, invalidJSON)

	// Act
	controller.SaveCredential(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockDataUseCase.AssertNotCalled(t, "SaveCredential")
}
