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

// Создаем мок для CloudUseCase
type MockCloudUseCase struct {
	mock.Mock
}

func (m *MockCloudUseCase) GenerateUploadLink(w http.ResponseWriter, r *http.Request, fileData *domain.FileData) {
	m.Called(w, r, fileData)
}

func (m *MockCloudUseCase) GenerateDownloadLink(w http.ResponseWriter, r *http.Request, label string) {
	m.Called(w, r, label)
}

// Тест для HandleUploadFile
func TestFileController_HandleUploadFile(t *testing.T) {
	// Arrange
	mockCloudUseCase := new(MockCloudUseCase)
	controller := NewFileController(mockCloudUseCase)
	
	// Создаем тестовые данные
	fileData := domain.FileData{
		Name:      "test-file.txt",
		Extension: "txt",
	}
	
	// Создаем JSON из данных
	jsonData, _ := json.Marshal(fileData)
	
	// Создаем запрос
	req, _ := http.NewRequest("POST", "/api/files/upload", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	
	// Настраиваем поведение мока
	mockCloudUseCase.On("GenerateUploadLink", mock.Anything, mock.Anything, mock.MatchedBy(func(f *domain.FileData) bool {
		return f.Name == fileData.Name && f.Extension == fileData.Extension
	}))
	
	// Act
	controller.HandleUploadFile(rr, req)
	
	// Assert
	mockCloudUseCase.AssertExpectations(t)
}

// Тест для HandleUploadFile с некорректными данными
func TestFileController_HandleUploadFile_InvalidData(t *testing.T) {
	// Arrange
	mockCloudUseCase := new(MockCloudUseCase)
	controller := NewFileController(mockCloudUseCase)
	
	// Создаем некорректный JSON
	invalidJSON := []byte(`{"name": "test-file.txt", "extension":}`)
	
	// Создаем запрос
	req, _ := http.NewRequest("POST", "/api/files/upload", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	
	// Act
	controller.HandleUploadFile(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockCloudUseCase.AssertNotCalled(t, "GenerateUploadLink")
}

// Тест для HandleDownloadFile
func TestFileController_HandleDownloadFile(t *testing.T) {
	// Arrange
	mockCloudUseCase := new(MockCloudUseCase)
	controller := NewFileController(mockCloudUseCase)
	
	// Создаем тестовые данные
	label := "test-label"
	
	// Создаем запрос с параметром label
	req, _ := http.NewRequest("GET", "/api/files/download?label="+label, nil)
	
	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	
	// Настраиваем поведение мока
	mockCloudUseCase.On("GenerateDownloadLink", mock.Anything, mock.Anything, label)
	
	// Act
	controller.HandleDownloadFile(rr, req)
	
	// Assert
	mockCloudUseCase.AssertExpectations(t)
}

// Тест для HandleDownloadFile без метки
func TestFileController_HandleDownloadFile_MissingLabel(t *testing.T) {
	// Arrange
	mockCloudUseCase := new(MockCloudUseCase)
	controller := NewFileController(mockCloudUseCase)
	
	// Создаем запрос без параметра label
	req, _ := http.NewRequest("GET", "/api/files/download", nil)
	
	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	
	// Act
	controller.HandleDownloadFile(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockCloudUseCase.AssertNotCalled(t, "GenerateDownloadLink")
}