package service

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"net/http/httptest"
	"testing"
	"time"
)

// Мок для ConfigServer
type MockConfigServer struct {
	jwtSecret string
}

func NewMockConfigServer() *MockConfigServer {
	return &MockConfigServer{
		jwtSecret: "test-secret-key",
	}
}

func (m *MockConfigServer) GetJwtSecret() string {
	return m.jwtSecret
}

func (m *MockConfigServer) GetDBDsn() string {
	return ""
}

func (m *MockConfigServer) GetRunAddr() string {
	return ""
}

func (m *MockConfigServer) GetMinioBucketName() string {
	return ""
}

func (m *MockConfigServer) GetMinioAccessKey() string {
	return ""
}

func (m *MockConfigServer) GetMinioSecretKey() string {
	return ""
}

func (m *MockConfigServer) GetMinioHost() string {
	return ""
}

func TestGenerateToken(t *testing.T) {
	// Arrange
	mockConfig := NewMockConfigServer()
	authService := NewAuthService(mockConfig)
	login := "testuser"

	// Act
	token, err := authService.GenerateToken(login)

	// Assert
	if err != nil {
		t.Fatalf("Ошибка при генерации токена: %v", err)
	}
	if token == "" {
		t.Fatal("Токен не должен быть пустым")
	}

	// Проверяем, что токен можно валидировать
	claims, err := authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Ошибка при валидации токена: %v", err)
	}
	if claims.Login != login {
		t.Fatalf("Логин в токене не совпадает: ожидается %s, получено %s", login, claims.Login)
	}
}

func TestValidateToken_Valid(t *testing.T) {
	// Arrange
	mockConfig := NewMockConfigServer()
	authService := NewAuthService(mockConfig)
	login := "testuser"

	// Создаем токен
	token, err := authService.GenerateToken(login)
	if err != nil {
		t.Fatalf("Ошибка при генерации токена: %v", err)
	}

	// Act
	claims, err := authService.ValidateToken(token)

	// Assert
	if err != nil {
		t.Fatalf("Ошибка при валидации токена: %v", err)
	}
	if claims == nil {
		t.Fatal("Claims не должны быть nil")
	}
	if claims.Login != login {
		t.Fatalf("Логин в токене не совпадает: ожидается %s, получено %s", login, claims.Login)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	// Arrange
	mockConfig := NewMockConfigServer()
	authService := NewAuthService(mockConfig)

	// Создаем невалидный токен
	invalidToken := "invalid.token.string"

	// Act
	claims, err := authService.ValidateToken(invalidToken)

	// Assert
	if err == nil {
		t.Fatal("Ожидалась ошибка при валидации невалидного токена")
	}
	if claims != nil {
		t.Fatal("Claims должны быть nil для невалидного токена")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	// Arrange
	mockConfig := NewMockConfigServer()
	authService := &AuthService{cf: mockConfig}
	login := "testuser"

	// Создаем токен с истекшим сроком действия
	expirationTime := time.Now().Add(-1 * time.Hour) // Истек 1 час назад
	claims := &domain.Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString([]byte(mockConfig.GetJwtSecret()))
	if err != nil {
		t.Fatalf("Ошибка при создании истекшего токена: %v", err)
	}

	// Act
	resultClaims, err := authService.ValidateToken(expiredToken)

	// Assert
	if err == nil {
		t.Fatal("Ожидалась ошибка при валидации истекшего токена")
	}
	if resultClaims != nil {
		t.Fatal("Claims должны быть nil для истекшего токена")
	}
}

func TestHashPassword(t *testing.T) {
	// Arrange
	mockConfig := NewMockConfigServer()
	authService := NewAuthService(mockConfig)
	password := "securepassword123"

	// Act
	hash, err := authService.HashPassword(password)

	// Assert
	if err != nil {
		t.Fatalf("Ошибка при хешировании пароля: %v", err)
	}
	if hash == "" {
		t.Fatal("Хеш не должен быть пустым")
	}
	if hash == password {
		t.Fatal("Хеш не должен совпадать с исходным паролем")
	}
}

func TestCheckPasswordHash_Valid(t *testing.T) {
	// Arrange
	mockConfig := NewMockConfigServer()
	authService := NewAuthService(mockConfig)
	password := "securepassword123"

	// Хешируем пароль
	hash, err := authService.HashPassword(password)
	if err != nil {
		t.Fatalf("Ошибка при хешировании пароля: %v", err)
	}

	// Act
	result := authService.CheckPasswordHash(password, hash)

	// Assert
	if !result {
		t.Fatal("Проверка хеша пароля должна быть успешной")
	}
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	// Arrange
	mockConfig := NewMockConfigServer()
	authService := NewAuthService(mockConfig)
	password := "securepassword123"
	wrongPassword := "wrongpassword456"

	// Хешируем пароль
	hash, err := authService.HashPassword(password)
	if err != nil {
		t.Fatalf("Ошибка при хешировании пароля: %v", err)
	}

	// Act
	result := authService.CheckPasswordHash(wrongPassword, hash)

	// Assert
	if result {
		t.Fatal("Проверка хеша пароля должна быть неуспешной для неверного пароля")
	}
}

func TestSetResponseAuthData(t *testing.T) {
	// Arrange
	mockConfig := NewMockConfigServer()
	authService := NewAuthService(mockConfig)
	token := "test-token"

	// Создаем тестовый HTTP-ответ
	w := httptest.NewRecorder()

	// Act
	authService.SetResponseAuthData(w, token)

	// Assert
	authHeader := w.Header().Get("Authorization")
	expectedHeader := "Bearer " + token
	if authHeader != expectedHeader {
		t.Fatalf("Заголовок Authorization не совпадает: ожидается %s, получено %s", expectedHeader, authHeader)
	}
}
