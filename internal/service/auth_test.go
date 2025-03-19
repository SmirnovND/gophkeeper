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

func (m *MockConfigServer) GetJwtSecret() string {
	return m.jwtSecret
}

func (m *MockConfigServer) GetDBDsn() string {
	return ""
}

func (m *MockConfigServer) GetRabbitMQURI() string {
	return ""
}

func (m *MockConfigServer) GetRunAddr() string {
	return ""
}

func TestAuthService_GenerateToken(t *testing.T) {
	// Arrange
	mockConfig := &MockConfigServer{jwtSecret: "test_secret"}
	authService := NewAuthService(mockConfig)
	login := "testuser"

	// Act
	token, err := authService.GenerateToken(login)

	// Assert
	if err != nil {
		t.Errorf("Ошибка при генерации токена: %v", err)
	}
	if token == "" {
		t.Error("Сгенерированный токен пустой")
	}

	// Проверяем, что токен можно валидировать
	claims, err := authService.ValidateToken(token)
	if err != nil {
		t.Errorf("Ошибка при валидации токена: %v", err)
	}
	if claims.Login != login {
		t.Errorf("Логин в токене не совпадает: ожидается %s, получено %s", login, claims.Login)
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	// Arrange
	mockConfig := &MockConfigServer{jwtSecret: "test_secret"}
	authService := NewAuthService(mockConfig)
	login := "testuser"

	// Создаем валидный токен
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &domain.Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, _ := token.SignedString([]byte(mockConfig.GetJwtSecret()))

	// Создаем невалидный токен с другим секретом
	invalidToken, _ := token.SignedString([]byte("wrong_secret"))

	// Создаем просроченный токен
	expiredClaims := &domain.Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte(mockConfig.GetJwtSecret()))

	tests := []struct {
		name        string
		tokenString string
		wantErr     bool
		wantLogin   string
	}{
		{
			name:        "Валидный токен",
			tokenString: validToken,
			wantErr:     false,
			wantLogin:   login,
		},
		{
			name:        "Невалидный токен (неверный секрет)",
			tokenString: invalidToken,
			wantErr:     true,
			wantLogin:   "",
		},
		{
			name:        "Просроченный токен",
			tokenString: expiredTokenString,
			wantErr:     true,
			wantLogin:   "",
		},
		{
			name:        "Пустой токен",
			tokenString: "",
			wantErr:     true,
			wantLogin:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			claims, err := authService.ValidateToken(tt.tokenString)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && claims.Login != tt.wantLogin {
				t.Errorf("ValidateToken() login = %v, want %v", claims.Login, tt.wantLogin)
			}
		})
	}
}

func TestAuthService_HashPassword(t *testing.T) {
	// Arrange
	mockConfig := &MockConfigServer{jwtSecret: "test_secret"}
	authService := NewAuthService(mockConfig)
	password := "password123"

	// Act
	hash, err := authService.HashPassword(password)

	// Assert
	if err != nil {
		t.Errorf("Ошибка при хешировании пароля: %v", err)
	}
	if hash == "" {
		t.Error("Хеш пароля пустой")
	}
	if hash == password {
		t.Error("Хеш пароля совпадает с исходным паролем")
	}

	// Проверяем, что пароль можно проверить с помощью хеша
	if !authService.CheckPasswordHash(password, hash) {
		t.Error("Проверка пароля с хешем не прошла")
	}
}

func TestAuthService_CheckPasswordHash(t *testing.T) {
	// Arrange
	mockConfig := &MockConfigServer{jwtSecret: "test_secret"}
	authService := NewAuthService(mockConfig)
	password := "password123"
	wrongPassword := "wrongpassword"

	// Создаем хеш пароля
	hash, _ := authService.HashPassword(password)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Правильный пароль",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "Неправильный пароль",
			password: wrongPassword,
			hash:     hash,
			want:     false,
		},
		{
			name:     "Пустой пароль",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Пустой хеш",
			password: password,
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := authService.CheckPasswordHash(tt.password, tt.hash)

			// Assert
			if result != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestAuthService_SetResponseAuthData(t *testing.T) {
	// Arrange
	mockConfig := &MockConfigServer{jwtSecret: "test_secret"}
	authService := NewAuthService(mockConfig)
	token := "test_token"
	w := httptest.NewRecorder()

	// Act
	authService.SetResponseAuthData(w, token)

	// Assert
	// Проверяем заголовок Authorization
	authHeader := w.Header().Get("Authorization")
	expectedAuthHeader := "Bearer " + token
	if authHeader != expectedAuthHeader {
		t.Errorf("Заголовок Authorization = %v, want %v", authHeader, expectedAuthHeader)
	}

	// Проверяем cookie
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Errorf("Ожидается 1 cookie, получено %d", len(cookies))
		return
	}

	cookie := cookies[0]
	if cookie.Name != "auth_token" {
		t.Errorf("Имя cookie = %v, want %v", cookie.Name, "auth_token")
	}
	if cookie.Value != token {
		t.Errorf("Значение cookie = %v, want %v", cookie.Value, token)
	}
	if cookie.Path != "/" {
		t.Errorf("Путь cookie = %v, want %v", cookie.Path, "/")
	}
	if !cookie.HttpOnly {
		t.Error("Cookie должен быть HttpOnly")
	}
}
