package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SmirnovND/gophkeeper/internal/controllers"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUseCase - мок для интерфейса AuthUseCase
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
	return args.Get(0).(*domain.Claims), args.Error(1)
}

func TestAuthController_HandleRegisterJSON(t *testing.T) {
	mockAuthUseCase := new(MockAuthUseCase)
	authController := controllers.NewAuthController(mockAuthUseCase)

	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}

	// Настраиваем поведение мока
	mockAuthUseCase.On("Register", mock.Anything, credentials).Return("jwt-token", nil)

	// Создаем HTTP-запрос
	body, _ := json.Marshal(credentials)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	w := httptest.NewRecorder()

	// Вызываем тестируемый метод
	authController.HandleRegisterJSON(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что все ожидаемые методы были вызваны
	mockAuthUseCase.AssertExpectations(t)
}

func TestAuthController_HandleLoginJSON(t *testing.T) {
	mockAuthUseCase := new(MockAuthUseCase)
	authController := controllers.NewAuthController(mockAuthUseCase)

	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}

	// Настраиваем поведение мока
	mockAuthUseCase.On("Login", mock.Anything, credentials).Return("jwt-token", nil)

	// Создаем HTTP-запрос
	body, _ := json.Marshal(credentials)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	w := httptest.NewRecorder()

	// Вызываем тестируемый метод
	authController.HandleLoginJSON(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что все ожидаемые методы были вызваны
	mockAuthUseCase.AssertExpectations(t)
}
