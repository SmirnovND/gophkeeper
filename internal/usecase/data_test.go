package usecase

import (
	"encoding/json"
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestNewDataUseCase проверяет создание нового экземпляра DataUseCase
func TestNewDataUseCase(t *testing.T) {
	mockDataService := &MockDataService{}
	mockJwtService := &MockJwtService{}
	dataUseCase := NewDataUseCase(mockDataService, mockJwtService)

	if dataUseCase == nil {
		t.Fatal("NewDataUseCase вернул nil")
	}
}

// TestDataUseCase_GetCredential_Success проверяет успешное получение учетных данных
func TestDataUseCase_GetCredential_Success(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		GetCredentialFunc: func(login string, label string) (*domain.CredentialData, string, error) {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-credential" {
				t.Errorf("Ожидалась метка 'test-credential', получена '%s'", label)
			}
			return &domain.CredentialData{
				Login:    "service-login",
				Password: "service-password",
			}, "", nil
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			if tokenString == "Bearer valid-token" {
				return "testuser", nil
			}
			return "", errors.New("неверный токен")
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/data/credential/test-credential", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GetCredential
	dataUseCase.GetCredential(w, req, "test-credential")

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
		CredentialData *domain.CredentialData `json:"credential_data"`
		Metadata       string                 `json:"metadata"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	if response.CredentialData.Login != "service-login" {
		t.Errorf("Ожидался логин 'service-login', получен '%s'", response.CredentialData.Login)
	}
	if response.CredentialData.Password != "service-password" {
		t.Errorf("Ожидался пароль 'service-password', получен '%s'", response.CredentialData.Password)
	}
}

// TestDataUseCase_GetCredential_TokenError проверяет обработку ошибки при извлечении логина из токена
func TestDataUseCase_GetCredential_TokenError(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "", errors.New("ошибка извлечения логина из токена")
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/data/credential/test-credential", nil)
	req.Header.Set("Authorization", "Bearer error-token")
	w := httptest.NewRecorder()

	// Вызываем метод GetCredential
	dataUseCase.GetCredential(w, req, "test-credential")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка получения логина") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка получения логина', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_GetCredential_NotFound проверяет обработку ошибки "учетные данные не найдены"
func TestDataUseCase_GetCredential_NotFound(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		GetCredentialFunc: func(login string, label string) (*domain.CredentialData, string, error) {
			return nil, "", errors.New("учетные данные не найдены")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/data/credential/test-credential", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GetCredential
	dataUseCase.GetCredential(w, req, "test-credential")

	// Проверяем статус ответа
	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "учетные данные не найдены") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'учетные данные не найдены', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_GetCredential_OtherError проверяет обработку других ошибок
func TestDataUseCase_GetCredential_OtherError(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		GetCredentialFunc: func(login string, label string) (*domain.CredentialData, string, error) {
			return nil, "", errors.New("внутренняя ошибка сервера")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/data/credential/test-credential", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GetCredential
	dataUseCase.GetCredential(w, req, "test-credential")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "внутренняя ошибка сервера") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'внутренняя ошибка сервера', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_SaveCredential_Success проверяет успешное сохранение учетных данных
func TestDataUseCase_SaveCredential_Success(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		SaveCredentialFunc: func(login string, label string, credentialData *domain.CredentialData, metadata string) error {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-credential" {
				t.Errorf("Ожидалась метка 'test-credential', получена '%s'", label)
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

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/data/credential/test-credential", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Создаем данные для сохранения
	credentialData := &domain.CredentialData{
		Login:    "service-login",
		Password: "service-password",
	}

	// Вызываем метод SaveCredential
	dataUseCase.SaveCredential(w, req, "test-credential", credentialData, "")

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
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	if message, ok := response["message"]; !ok || message != "учетные данные успешно сохранены" {
		t.Errorf("Ожидалось сообщение 'учетные данные успешно сохранены', получено '%v'", response)
	}
}

// TestDataUseCase_SaveCredential_TokenError проверяет обработку ошибки при извлечении логина из токена
func TestDataUseCase_SaveCredential_TokenError(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "", errors.New("ошибка извлечения логина из токена")
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/data/credential/test-credential", nil)
	req.Header.Set("Authorization", "Bearer error-token")
	w := httptest.NewRecorder()

	// Создаем данные для сохранения
	credentialData := &domain.CredentialData{
		Login:    "service-login",
		Password: "service-password",
	}

	// Вызываем метод SaveCredential
	dataUseCase.SaveCredential(w, req, "test-credential", credentialData, "")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка получения логина") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка получения логина', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_SaveCredential_Error проверяет обработку ошибки при сохранении учетных данных
func TestDataUseCase_SaveCredential_Error(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		SaveCredentialFunc: func(login string, label string, credentialData *domain.CredentialData, metadata string) error {
			return errors.New("ошибка при сохранении учетных данных")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/data/credential/test-credential", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Создаем данные для сохранения
	credentialData := &domain.CredentialData{
		Login:    "service-login",
		Password: "service-password",
	}

	// Вызываем метод SaveCredential
	dataUseCase.SaveCredential(w, req, "test-credential", credentialData, "")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "ошибка при сохранении учетных данных") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'ошибка при сохранении учетных данных', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_GetCard_Success проверяет успешное получение данных карты
func TestDataUseCase_GetCard_Success(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		GetCardFunc: func(login string, label string) (*domain.CardData, string, error) {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-card" {
				t.Errorf("Ожидалась метка 'test-card', получена '%s'", label)
			}
			return &domain.CardData{
				Number:     "1234567890123456",
				Holder:     "Test User",
				ExpiryDate: "12/25",
				CVV:        "123",
			}, "", nil
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/data/card/test-card", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GetCard
	dataUseCase.GetCard(w, req, "test-card")

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
		CardData *domain.CardData `json:"card_data"`
		Metadata string           `json:"metadata"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	if response.CardData.Number != "1234567890123456" {
		t.Errorf("Ожидался номер карты '1234567890123456', получен '%s'", response.CardData.Number)
	}
	if response.CardData.Holder != "Test User" {
		t.Errorf("Ожидался держатель карты 'Test User', получен '%s'", response.CardData.Holder)
	}
	if response.CardData.ExpiryDate != "12/25" {
		t.Errorf("Ожидался срок действия '12/25', получен '%s'", response.CardData.ExpiryDate)
	}
	if response.CardData.CVV != "123" {
		t.Errorf("Ожидался CVV '123', получен '%s'", response.CardData.CVV)
	}
}

// TestDataUseCase_GetCard_NotFound проверяет обработку ошибки "данные карты не найдены"
func TestDataUseCase_GetCard_NotFound(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		GetCardFunc: func(login string, label string) (*domain.CardData, string, error) {
			return nil, "", errors.New("данные карты не найдены")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/data/card/test-card", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GetCard
	dataUseCase.GetCard(w, req, "test-card")

	// Проверяем статус ответа
	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "данные карты не найдены") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'данные карты не найдены', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_SaveCard_Success проверяет успешное сохранение данных карты
func TestDataUseCase_SaveCard_Success(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		SaveCardFunc: func(login string, label string, cardData *domain.CardData, metadata string) error {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-card" {
				t.Errorf("Ожидалась метка 'test-card', получена '%s'", label)
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

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/data/card/test-card", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Создаем данные для сохранения
	cardData := &domain.CardData{
		Number:     "1234567890123456",
		Holder:     "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}

	// Вызываем метод SaveCard
	dataUseCase.SaveCard(w, req, "test-card", cardData, "")

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
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	if message, ok := response["message"]; !ok || message != "данные карты успешно сохранены" {
		t.Errorf("Ожидалось сообщение 'данные карты успешно сохранены', получено '%v'", response)
	}
}

// TestDataUseCase_GetText_Success проверяет успешное получение текстовых данных
func TestDataUseCase_GetText_Success(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		GetTextFunc: func(login string, label string) (*domain.TextData, string, error) {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-text" {
				t.Errorf("Ожидалась метка 'test-text', получена '%s'", label)
			}
			return &domain.TextData{
				Content: "test text content",
			}, "", nil
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/data/text/test-text", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GetText
	dataUseCase.GetText(w, req, "test-text")

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
		TextData *domain.TextData `json:"text_data"`
		Metadata string           `json:"metadata"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	if response.TextData.Content != "test text content" {
		t.Errorf("Ожидался текст 'test text content', получен '%s'", response.TextData.Content)
	}
}

// TestDataUseCase_GetText_NotFound проверяет обработку ошибки "текстовые данные не найдены"
func TestDataUseCase_GetText_NotFound(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		GetTextFunc: func(login string, label string) (*domain.TextData, string, error) {
			return nil, "", errors.New("текстовые данные не найдены")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("GET", "/api/data/text/test-text", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод GetText
	dataUseCase.GetText(w, req, "test-text")

	// Проверяем статус ответа
	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "текстовые данные не найдены") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'текстовые данные не найдены', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_SaveText_Success проверяет успешное сохранение текстовых данных
func TestDataUseCase_SaveText_Success(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		SaveTextFunc: func(login string, label string, textData *domain.TextData, metadata string) error {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-text" {
				t.Errorf("Ожидалась метка 'test-text', получена '%s'", label)
			}
			if textData.Content != "test text content" {
				t.Errorf("Ожидался текст 'test text content', получен '%s'", textData.Content)
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

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/data/text/test-text", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Создаем данные для сохранения
	textData := &domain.TextData{
		Content: "test text content",
	}

	// Вызываем метод SaveText
	dataUseCase.SaveText(w, req, "test-text", textData, "")

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
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	if message, ok := response["message"]; !ok || message != "текстовые данные успешно сохранены" {
		t.Errorf("Ожидалось сообщение 'текстовые данные успешно сохранены', получено '%v'", response)
	}
}

// TestDataUseCase_SaveText_TokenError проверяет обработку ошибки при извлечении логина из токена
func TestDataUseCase_SaveText_TokenError(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "", errors.New("ошибка извлечения логина из токена")
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/data/text/test-text", nil)
	req.Header.Set("Authorization", "Bearer error-token")
	w := httptest.NewRecorder()

	// Создаем данные для сохранения
	textData := &domain.TextData{
		Content: "test text content",
	}

	// Вызываем метод SaveText
	dataUseCase.SaveText(w, req, "test-text", textData, "")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка получения логина") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка получения логина', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_SaveText_Error проверяет обработку ошибки при сохранении текстовых данных
func TestDataUseCase_SaveText_Error(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		SaveTextFunc: func(login string, label string, textData *domain.TextData, metadata string) error {
			return errors.New("ошибка при сохранении текстовых данных")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("POST", "/api/data/text/test-text", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Создаем данные для сохранения
	textData := &domain.TextData{
		Content: "test text content",
	}

	// Вызываем метод SaveText
	dataUseCase.SaveText(w, req, "test-text", textData, "")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "ошибка при сохранении текстовых данных") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'ошибка при сохранении текстовых данных', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_DeleteCredential_Success проверяет успешное удаление учетных данных
func TestDataUseCase_DeleteCredential_Success(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		DeleteCredentialFunc: func(login string, label string) error {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-credential" {
				t.Errorf("Ожидалась метка 'test-credential', получена '%s'", label)
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

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/credential/test-credential", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteCredential
	dataUseCase.DeleteCredential(w, req, "test-credential")

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
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	if message, ok := response["message"]; !ok || message != "учетные данные успешно удалены" {
		t.Errorf("Ожидалось сообщение 'учетные данные успешно удалены', получено '%v'", response)
	}
}

// TestDataUseCase_DeleteCredential_NotFound проверяет обработку ошибки "учетные данные не найдены"
func TestDataUseCase_DeleteCredential_NotFound(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		DeleteCredentialFunc: func(login string, label string) error {
			return errors.New("учетные данные не найдены")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/credential/test-credential", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteCredential
	dataUseCase.DeleteCredential(w, req, "test-credential")

	// Проверяем статус ответа
	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "учетные данные не найдены") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'учетные данные не найдены', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_DeleteCard_Success проверяет успешное удаление данных карты
func TestDataUseCase_DeleteCard_Success(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		DeleteCardFunc: func(login string, label string) error {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-card" {
				t.Errorf("Ожидалась метка 'test-card', получена '%s'", label)
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

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/card/test-card", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteCard
	dataUseCase.DeleteCard(w, req, "test-card")

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
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	if message, ok := response["message"]; !ok || message != "данные карты успешно удалены" {
		t.Errorf("Ожидалось сообщение 'данные карты успешно удалены', получено '%v'", response)
	}
}

// TestDataUseCase_DeleteCard_TokenError проверяет обработку ошибки при извлечении логина из токена
func TestDataUseCase_DeleteCard_TokenError(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "", errors.New("ошибка извлечения логина из токена")
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/card/test-card", nil)
	req.Header.Set("Authorization", "Bearer error-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteCard
	dataUseCase.DeleteCard(w, req, "test-card")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка получения логина") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка получения логина', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_DeleteCard_NotFound проверяет обработку ошибки "данные карты не найдены"
func TestDataUseCase_DeleteCard_NotFound(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		DeleteCardFunc: func(login string, label string) error {
			return errors.New("данные карты не найдены")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/card/test-card", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteCard
	dataUseCase.DeleteCard(w, req, "test-card")

	// Проверяем статус ответа
	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "данные карты не найдены") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'данные карты не найдены', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_DeleteCard_OtherError проверяет обработку других ошибок
func TestDataUseCase_DeleteCard_OtherError(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		DeleteCardFunc: func(login string, label string) error {
			return errors.New("внутренняя ошибка сервера")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/card/test-card", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteCard
	dataUseCase.DeleteCard(w, req, "test-card")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "внутренняя ошибка сервера") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'внутренняя ошибка сервера', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_DeleteText_Success проверяет успешное удаление текстовых данных
func TestDataUseCase_DeleteText_Success(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		DeleteTextFunc: func(login string, label string) error {
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			if label != "test-text" {
				t.Errorf("Ожидалась метка 'test-text', получена '%s'", label)
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

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/text/test-text", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteText
	dataUseCase.DeleteText(w, req, "test-text")

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
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Ошибка при разборе JSON ответа: %v", err)
	}

	if message, ok := response["message"]; !ok || message != "текстовые данные успешно удалены" {
		t.Errorf("Ожидалось сообщение 'текстовые данные успешно удалены', получено '%v'", response)
	}
}

// TestDataUseCase_DeleteText_TokenError проверяет обработку ошибки при извлечении логина из токена
func TestDataUseCase_DeleteText_TokenError(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "", errors.New("ошибка извлечения логина из токена")
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/text/test-text", nil)
	req.Header.Set("Authorization", "Bearer error-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteText
	dataUseCase.DeleteText(w, req, "test-text")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "Ошибка получения логина") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'Ошибка получения логина', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_DeleteText_NotFound проверяет обработку ошибки "текстовые данные не найдены"
func TestDataUseCase_DeleteText_NotFound(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		DeleteTextFunc: func(login string, label string) error {
			return errors.New("текстовые данные не найдены")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/text/test-text", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteText
	dataUseCase.DeleteText(w, req, "test-text")

	// Проверяем статус ответа
	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "текстовые данные не найдены") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'текстовые данные не найдены', получено '%s'", w.Body.String())
	}
}

// TestDataUseCase_DeleteText_OtherError проверяет обработку других ошибок
func TestDataUseCase_DeleteText_OtherError(t *testing.T) {
	// Создаем мок для DataService
	mockDataService := &MockDataService{
		DeleteTextFunc: func(login string, label string) error {
			return errors.New("внутренняя ошибка сервера")
		},
	}

	// Создаем мок для JwtService
	mockJwtService := &MockJwtService{
		ExtractLoginFromTokenFunc: func(tokenString string) (string, error) {
			return "testuser", nil
		},
	}

	// Создаем экземпляр DataUseCase
	dataUseCase := &DataUseCase{
		dataService: mockDataService,
		jwtService:  mockJwtService,
	}

	// Создаем тестовый HTTP запрос и ответ
	req := httptest.NewRequest("DELETE", "/api/data/text/test-text", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Вызываем метод DeleteText
	dataUseCase.DeleteText(w, req, "test-text")

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(w.Body.String(), "внутренняя ошибка сервера") {
		t.Errorf("Ожидалось сообщение об ошибке с текстом 'внутренняя ошибка сервера', получено '%s'", w.Body.String())
	}
}
