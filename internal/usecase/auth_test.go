package usecase

import (
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Создаем мок для UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindUser(login string) (*domain.User, error) {
	args := m.Called(login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) SaveUser(login string, password string) (*domain.User, error) {
	args := m.Called(login, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Создаем мок для AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateToken(login string) (string, error) {
	args := m.Called(login)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*domain.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Claims), args.Error(1)
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

func (m *MockAuthService) SetResponseAuthData(w http.ResponseWriter, token string) {
	m.Called(w, token)
}

func TestAuthUseCase_Register_Success(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	
	user := &domain.User{
		Credentials: domain.Credentials{
			Id:       1,
			Login:    "testuser",
			Password: "password123",
			PassHash: "hashedpassword",
		},
	}
	
	// Настраиваем поведение моков
	mockUserService.On("FindUser", "testuser").Return(nil, domain.ErrNotFound)
	mockUserService.On("SaveUser", "testuser", "password123").Return(user, nil)
	mockAuthService.On("GenerateToken", "testuser").Return("jwt-token", nil)
	mockAuthService.On("SetResponseAuthData", mock.Anything, "jwt-token").Return()
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Register(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `{"status": "success"}`)
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestAuthUseCase_Register_UserAlreadyExists(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "existinguser",
		Password: "password123",
	}
	
	existingUser := &domain.User{
		Credentials: domain.Credentials{
			Id:       1,
			Login:    "existinguser",
			PassHash: "hashedpassword",
		},
	}
	
	// Настраиваем поведение моков - пользователь уже существует
	mockUserService.On("FindUser", "existinguser").Return(existingUser, nil)
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Register(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
}

func TestAuthUseCase_Register_FindUserError(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	
	// Настраиваем поведение моков - ошибка при поиске пользователя
	mockUserService.On("FindUser", "testuser").Return(nil, errors.New("database error"))
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Register(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
}

func TestAuthUseCase_Register_SaveUserError(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	
	// Настраиваем поведение моков
	mockUserService.On("FindUser", "testuser").Return(nil, domain.ErrNotFound)
	mockUserService.On("SaveUser", "testuser", "password123").Return(nil, errors.New("save error"))
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Register(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error saving user")
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
}

func TestAuthUseCase_Register_GenerateTokenError(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	
	user := &domain.User{
		Credentials: domain.Credentials{
			Id:       1,
			Login:    "testuser",
			Password: "password123",
			PassHash: "hashedpassword",
		},
	}
	
	// Настраиваем поведение моков
	mockUserService.On("FindUser", "testuser").Return(nil, domain.ErrNotFound)
	mockUserService.On("SaveUser", "testuser", "password123").Return(user, nil)
	mockAuthService.On("GenerateToken", "testuser").Return("", errors.New("token generation error"))
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Register(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error generating token")
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestAuthUseCase_Login_Success(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	
	user := &domain.User{
		Credentials: domain.Credentials{
			Id:       1,
			Login:    "testuser",
			PassHash: "hashedpassword",
		},
	}
	
	// Настраиваем поведение моков
	mockUserService.On("FindUser", "testuser").Return(user, nil)
	mockAuthService.On("CheckPasswordHash", "password123", "hashedpassword").Return(true)
	mockAuthService.On("GenerateToken", "testuser").Return("jwt-token", nil)
	mockAuthService.On("SetResponseAuthData", mock.Anything, "jwt-token").Return()
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Login(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `{"status": "success"}`)
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestAuthUseCase_Login_UserNotFound(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "nonexistentuser",
		Password: "password123",
	}
	
	// Настраиваем поведение моков - пользователь не найден
	mockUserService.On("FindUser", "nonexistentuser").Return(nil, domain.ErrNotFound)
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Login(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
}

func TestAuthUseCase_Login_FindUserError(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	
	// Настраиваем поведение моков - ошибка при поиске пользователя
	mockUserService.On("FindUser", "testuser").Return(nil, errors.New("database error"))
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Login(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
}

func TestAuthUseCase_Login_InvalidPassword(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "wrongpassword",
	}
	
	user := &domain.User{
		Credentials: domain.Credentials{
			Id:       1,
			Login:    "testuser",
			PassHash: "hashedpassword",
		},
	}
	
	// Настраиваем поведение моков - неверный пароль
	mockUserService.On("FindUser", "testuser").Return(user, nil)
	mockAuthService.On("CheckPasswordHash", "wrongpassword", "hashedpassword").Return(false)
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Login(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestAuthUseCase_Login_GenerateTokenError(t *testing.T) {
	// Arrange
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)
	
	authUseCase := NewAuthUseCase(mockUserService, mockAuthService)
	
	credentials := &domain.Credentials{
		Login:    "testuser",
		Password: "password123",
	}
	
	user := &domain.User{
		Credentials: domain.Credentials{
			Id:       1,
			Login:    "testuser",
			PassHash: "hashedpassword",
		},
	}
	
	// Настраиваем поведение моков - ошибка при генерации токена
	mockUserService.On("FindUser", "testuser").Return(user, nil)
	mockAuthService.On("CheckPasswordHash", "password123", "hashedpassword").Return(true)
	mockAuthService.On("GenerateToken", "testuser").Return("", errors.New("token generation error"))
	
	// Создаем ResponseWriter для теста
	w := httptest.NewRecorder()
	
	// Act
	authUseCase.Login(w, credentials)
	
	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	// Проверяем, что все ожидаемые методы были вызваны
	mockUserService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}