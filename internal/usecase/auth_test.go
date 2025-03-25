package usecase

import (
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockUserService - мок для интерфейса UserService
type MockUserService struct {
	FindUserFunc func(login string) (*domain.User, error)
	SaveUserFunc func(login string, password string) (*domain.User, error)
}

func (m *MockUserService) FindUser(login string) (*domain.User, error) {
	return m.FindUserFunc(login)
}

func (m *MockUserService) SaveUser(login string, password string) (*domain.User, error) {
	return m.SaveUserFunc(login, password)
}

// MockAuthService - мок для интерфейса AuthService
type MockAuthService struct {
	GenerateTokenFunc       func(login string) (string, error)
	ValidateTokenFunc       func(tokenString string) (*domain.Claims, error)
	HashPasswordFunc        func(password string) (string, error)
	CheckPasswordHashFunc   func(password, hash string) bool
	SetResponseAuthDataFunc func(w http.ResponseWriter, token string)
}

func (m *MockAuthService) GenerateToken(login string) (string, error) {
	return m.GenerateTokenFunc(login)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*domain.Claims, error) {
	return m.ValidateTokenFunc(tokenString)
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	return m.HashPasswordFunc(password)
}

func (m *MockAuthService) CheckPasswordHash(password, hash string) bool {
	return m.CheckPasswordHashFunc(password, hash)
}

func (m *MockAuthService) SetResponseAuthData(w http.ResponseWriter, token string) {
	m.SetResponseAuthDataFunc(w, token)
}

// TestAuthUseCase_Register_Success тестирует успешную регистрацию пользователя
func TestAuthUseCase_Register_Success(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Пользователь не найден - это хорошо для регистрации
			return nil, domain.ErrNotFound
		},
		SaveUserFunc: func(login string, password string) (*domain.User, error) {
			// Успешное сохранение пользователя
			return &domain.User{
				Credentials: domain.Credentials{
					Login:    login,
					PassHash: "hashed_password",
				},
			}, nil
		},
	}

	mockAuthService := &MockAuthService{
		GenerateTokenFunc: func(login string) (string, error) {
			// Успешная генерация токена
			return "test_token", nil
		},
		SetResponseAuthDataFunc: func(w http.ResponseWriter, token string) {
			// Устанавливаем заголовок с токеном
			w.Header().Set("Authorization", "Bearer "+token)
		},
	}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Вызываем метод Register
	token, err := authUseCase.Register(w, credentials)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при регистрации: %v", err)
	}

	if token != "test_token" {
		t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}

	// Проверяем заголовок Authorization
	if w.Header().Get("Authorization") != "Bearer test_token" {
		t.Errorf("Ожидался заголовок Authorization 'Bearer test_token', получен '%s'", w.Header().Get("Authorization"))
	}

	// Проверяем тело ответа
	expectedBody := `{"status": "success"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Ожидалось тело ответа '%s', получено '%s'", expectedBody, w.Body.String())
	}
}

// TestAuthUseCase_Register_UserExists тестирует регистрацию с уже существующим пользователем
func TestAuthUseCase_Register_UserExists(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Пользователь уже существует
			return &domain.User{
				Credentials: domain.Credentials{
					Login:    login,
					PassHash: "existing_hash",
				},
			}, nil
		},
	}

	mockAuthService := &MockAuthService{}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "existinguser",
		Password: "testpassword",
	}

	// Вызываем метод Register
	token, err := authUseCase.Register(w, credentials)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	if token != "" {
		t.Errorf("Ожидался пустой токен, получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusConflict {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusConflict, w.Code)
	}
}

// TestAuthUseCase_Register_FindUserError тестирует ошибку при поиске пользователя
func TestAuthUseCase_Register_FindUserError(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Ошибка при поиске пользователя (не ErrNotFound)
			return nil, errors.New("database error")
		},
	}

	mockAuthService := &MockAuthService{}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Вызываем метод Register
	token, err := authUseCase.Register(w, credentials)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	if token != "" {
		t.Errorf("Ожидался пустой токен, получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}
}

// TestAuthUseCase_Register_SaveUserError тестирует ошибку при сохранении пользователя
func TestAuthUseCase_Register_SaveUserError(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Пользователь не найден - это хорошо для регистрации
			return nil, domain.ErrNotFound
		},
		SaveUserFunc: func(login string, password string) (*domain.User, error) {
			// Ошибка при сохранении пользователя
			return nil, errors.New("database error")
		},
	}

	mockAuthService := &MockAuthService{}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Вызываем метод Register
	token, err := authUseCase.Register(w, credentials)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	if token != "" {
		t.Errorf("Ожидался пустой токен, получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}
}

// TestAuthUseCase_Register_GenerateTokenError тестирует ошибку при генерации токена
func TestAuthUseCase_Register_GenerateTokenError(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Пользователь не найден - это хорошо для регистрации
			return nil, domain.ErrNotFound
		},
		SaveUserFunc: func(login string, password string) (*domain.User, error) {
			// Успешное сохранение пользователя
			return &domain.User{
				Credentials: domain.Credentials{
					Login:    login,
					PassHash: "hashed_password",
				},
			}, nil
		},
	}

	mockAuthService := &MockAuthService{
		GenerateTokenFunc: func(login string) (string, error) {
			// Ошибка при генерации токена
			return "", errors.New("token generation error")
		},
	}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Вызываем метод Register
	token, err := authUseCase.Register(w, credentials)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	if token != "" {
		t.Errorf("Ожидался пустой токен, получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}
}

// TestAuthUseCase_Login_Success тестирует успешный вход пользователя
func TestAuthUseCase_Login_Success(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Пользователь найден
			return &domain.User{
				Credentials: domain.Credentials{
					Login:    login,
					PassHash: "hashed_password",
				},
			}, nil
		},
	}

	mockAuthService := &MockAuthService{
		CheckPasswordHashFunc: func(password, hash string) bool {
			// Пароль верный
			return true
		},
		GenerateTokenFunc: func(login string) (string, error) {
			// Успешная генерация токена
			return "test_token", nil
		},
		SetResponseAuthDataFunc: func(w http.ResponseWriter, token string) {
			// Устанавливаем заголовок с токеном
			w.Header().Set("Authorization", "Bearer "+token)
		},
	}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Вызываем метод Login
	token, err := authUseCase.Login(w, credentials)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при входе: %v", err)
	}

	if token != "test_token" {
		t.Errorf("Ожидался токен 'test_token', получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusOK, w.Code)
	}

	// Проверяем заголовок Authorization
	if w.Header().Get("Authorization") != "Bearer test_token" {
		t.Errorf("Ожидался заголовок Authorization 'Bearer test_token', получен '%s'", w.Header().Get("Authorization"))
	}

	// Проверяем тело ответа
	expectedBody := `{"status": "success"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Ожидалось тело ответа '%s', получено '%s'", expectedBody, w.Body.String())
	}
}

// TestAuthUseCase_Login_UserNotFound тестирует вход с несуществующим пользователем
func TestAuthUseCase_Login_UserNotFound(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Пользователь не найден
			return nil, domain.ErrNotFound
		},
	}

	mockAuthService := &MockAuthService{}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "nonexistentuser",
		Password: "testpassword",
	}

	// Вызываем метод Login
	token, err := authUseCase.Login(w, credentials)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	if token != "" {
		t.Errorf("Ожидался пустой токен, получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}
}

// TestAuthUseCase_Login_FindUserError тестирует ошибку при поиске пользователя
func TestAuthUseCase_Login_FindUserError(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Ошибка при поиске пользователя (не ErrNotFound)
			return nil, errors.New("database error")
		},
	}

	mockAuthService := &MockAuthService{}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Вызываем метод Login
	token, err := authUseCase.Login(w, credentials)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	if token != "" {
		t.Errorf("Ожидался пустой токен, получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}
}

// TestAuthUseCase_Login_InvalidPassword тестирует вход с неверным паролем
func TestAuthUseCase_Login_InvalidPassword(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Пользователь найден
			return &domain.User{
				Credentials: domain.Credentials{
					Login:    login,
					PassHash: "hashed_password",
				},
			}, nil
		},
	}

	mockAuthService := &MockAuthService{
		CheckPasswordHashFunc: func(password, hash string) bool {
			// Пароль неверный
			return false
		},
	}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "wrongpassword",
	}

	// Вызываем метод Login
	token, err := authUseCase.Login(w, credentials)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	if token != "" {
		t.Errorf("Ожидался пустой токен, получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusUnauthorized, w.Code)
	}
}

// TestAuthUseCase_Login_GenerateTokenError тестирует ошибку при генерации токена
func TestAuthUseCase_Login_GenerateTokenError(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Пользователь найден
			return &domain.User{
				Credentials: domain.Credentials{
					Login:    login,
					PassHash: "hashed_password",
				},
			}, nil
		},
	}

	mockAuthService := &MockAuthService{
		CheckPasswordHashFunc: func(password, hash string) bool {
			// Пароль верный
			return true
		},
		GenerateTokenFunc: func(login string) (string, error) {
			// Ошибка при генерации токена
			return "", errors.New("token generation error")
		},
	}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Создаем тестовый ResponseWriter
	w := httptest.NewRecorder()

	// Создаем тестовые учетные данные
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Вызываем метод Login
	token, err := authUseCase.Login(w, credentials)

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	if token != "" {
		t.Errorf("Ожидался пустой токен, получен '%s'", token)
	}

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusInternalServerError, w.Code)
	}
}

// TestAuthUseCase_ValidateToken_Success тестирует успешную валидацию токена
func TestAuthUseCase_ValidateToken_Success(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{}

	expectedClaims := &domain.Claims{
		Login: "testuser",
	}

	mockAuthService := &MockAuthService{
		ValidateTokenFunc: func(tokenString string) (*domain.Claims, error) {
			// Успешная валидация токена
			return expectedClaims, nil
		},
	}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Вызываем метод ValidateToken
	claims, err := authUseCase.ValidateToken("test_token")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при валидации токена: %v", err)
	}

	if claims != expectedClaims {
		t.Errorf("Ожидались claims %v, получены %v", expectedClaims, claims)
	}
}

// TestAuthUseCase_ValidateToken_Error тестирует ошибку при валидации токена
func TestAuthUseCase_ValidateToken_Error(t *testing.T) {
	// Создаем моки
	mockUserService := &MockUserService{}

	mockAuthService := &MockAuthService{
		ValidateTokenFunc: func(tokenString string) (*domain.Claims, error) {
			// Ошибка при валидации токена
			return nil, errors.New("invalid token")
		},
	}

	// Создаем экземпляр AuthUseCase
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)

	// Вызываем метод ValidateToken
	claims, err := authUseCase.ValidateToken("invalid_token")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	if claims != nil {
		t.Errorf("Ожидались nil claims, получены %v", claims)
	}
}
