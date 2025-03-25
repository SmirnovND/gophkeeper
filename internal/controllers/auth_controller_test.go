package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Создаем мок для AuthUseCase
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Login(w http.ResponseWriter, credentials *domain.Credentials) (string, error) {
	args := m.Called(w, credentials)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUseCase) Register(w http.ResponseWriter, credentials *domain.Credentials) (string, error) {
	args := m.Called(w, credentials)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUseCase) ValidateToken(token string) (*domain.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Claims), args.Error(1)
}

// Тест для HandleRegisterJSON
func TestAuthController_HandleRegisterJSON(t *testing.T) {
	// Arrange
	mockAuthUseCase := new(MockAuthUseCase)
	controller := NewAuthController(mockAuthUseCase)
	
	// Создаем тестовые данные
	credentials := domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	
	// Создаем JSON из данных
	jsonData, _ := json.Marshal(credentials)
	
	// Создаем запрос
	req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	
	// Настраиваем поведение мока
	mockAuthUseCase.On("Register", mock.Anything, mock.MatchedBy(func(c *domain.Credentials) bool {
		return c.Login == credentials.Login && c.Password == credentials.Password
	})).Return("token123", nil)
	
	// Act
	controller.HandleRegisterJSON(rr, req)
	
	// Assert
	mockAuthUseCase.AssertExpectations(t)
}

// Тест для HandleLoginJSON
func TestAuthController_HandleLoginJSON(t *testing.T) {
	// Arrange
	mockAuthUseCase := new(MockAuthUseCase)
	controller := NewAuthController(mockAuthUseCase)
	
	// Создаем тестовые данные
	credentials := domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	
	// Создаем JSON из данных
	jsonData, _ := json.Marshal(credentials)
	
	// Создаем запрос
	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	
	// Настраиваем поведение мока
	mockAuthUseCase.On("Login", mock.Anything, mock.MatchedBy(func(c *domain.Credentials) bool {
		return c.Login == credentials.Login && c.Password == credentials.Password
	})).Return("token123", nil)
	
	// Act
	controller.HandleLoginJSON(rr, req)
	
	// Assert
	mockAuthUseCase.AssertExpectations(t)
}

// Тест для HandleRegisterJSON с некорректными данными
func TestAuthController_HandleRegisterJSON_InvalidData(t *testing.T) {
	// Arrange
	mockAuthUseCase := new(MockAuthUseCase)
	controller := NewAuthController(mockAuthUseCase)
	
	// Создаем некорректный JSON
	invalidJSON := []byte(`{"login": "testuser", "password":}`)
	
	// Создаем запрос
	req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	
	// Act
	controller.HandleRegisterJSON(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAuthUseCase.AssertNotCalled(t, "Register")
}

// Тест для HandleLoginJSON с некорректными данными
func TestAuthController_HandleLoginJSON_InvalidData(t *testing.T) {
	// Arrange
	mockAuthUseCase := new(MockAuthUseCase)
	controller := NewAuthController(mockAuthUseCase)
	
	// Создаем некорректный JSON
	invalidJSON := []byte(`{"login": "testuser", "password":}`)
	
	// Создаем запрос
	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	
	// Act
	controller.HandleLoginJSON(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAuthUseCase.AssertNotCalled(t, "Login")
}