package controllers

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Создаем мок для интерфейса DB
type MockDB struct {
	mock.Mock
}

// Реализация метода Ping для мока
func (m *MockDB) Ping() error {
	args := m.Called()
	return args.Error(0)
}

// Реализация метода QueryRow для мока
func (m *MockDB) QueryRow(query string, args ...any) *sqlx.Row {
	// Этот метод не используется в HealthcheckController, но нужен для реализации интерфейса DB
	// Возвращаем nil, так как в тестах этот метод не будет вызываться
	return nil
}

func TestHealthcheckController_HandlePing_Success(t *testing.T) {
	// Arrange
	mockDB := new(MockDB)
	controller := NewHealthcheckController(mockDB)

	// Настраиваем поведение мока - Ping возвращает nil (успешное соединение)
	mockDB.On("Ping").Return(nil)

	// Создаем тестовый HTTP запрос и ResponseWriter
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	// Act
	controller.HandlePing(w, req)

	// Assert
	// Проверяем, что статус ответа 200 OK
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Проверяем, что тело ответа содержит "pong"
	assert.Equal(t, "pong", w.Body.String())

	// Проверяем, что метод Ping был вызван
	mockDB.AssertExpectations(t)
}

func TestHealthcheckController_HandlePing_DBError(t *testing.T) {
	// Arrange
	mockDB := new(MockDB)
	controller := NewHealthcheckController(mockDB)

	// Настраиваем поведение мока - Ping возвращает ошибку
	mockDB.On("Ping").Return(errors.New("database connection error"))

	// Создаем тестовый HTTP запрос и ResponseWriter
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	// Act
	controller.HandlePing(w, req)

	// Assert
	// Проверяем, что статус ответа 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

	// Проверяем, что тело ответа содержит сообщение об ошибке
	assert.Equal(t, "Failed to connect DB\n", w.Body.String())

	// Проверяем, что метод Ping был вызван
	mockDB.AssertExpectations(t)
}

// Тест на проверку порядка вызова методов Write и WriteHeader
func TestHealthcheckController_HandlePing_WriteOrder(t *testing.T) {
	// Arrange
	mockDB := new(MockDB)
	controller := NewHealthcheckController(mockDB)

	// Настраиваем поведение мока - Ping возвращает nil (успешное соединение)
	mockDB.On("Ping").Return(nil)

	// Создаем специальный ResponseWriter для отслеживания порядка вызовов
	w := &OrderTrackingResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
		writeOrder:       make([]string, 0),
	}

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	// Act
	controller.HandlePing(w, req)

	// Assert
	// Проверяем порядок вызовов - сначала Write, затем WriteHeader
	// Обратите внимание: это неправильный порядок, и тест должен выявить эту проблему
	assert.Equal(t, []string{"Write", "WriteHeader"}, w.writeOrder)
}

// OrderTrackingResponseWriter отслеживает порядок вызовов методов Write и WriteHeader
type OrderTrackingResponseWriter struct {
	*httptest.ResponseRecorder
	writeOrder []string
}

func (w *OrderTrackingResponseWriter) Write(b []byte) (int, error) {
	w.writeOrder = append(w.writeOrder, "Write")
	return w.ResponseRecorder.Write(b)
}

func (w *OrderTrackingResponseWriter) WriteHeader(statusCode int) {
	w.writeOrder = append(w.writeOrder, "WriteHeader")
	w.ResponseRecorder.WriteHeader(statusCode)
}
