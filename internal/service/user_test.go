package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для UserRepo
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindUser(login string) (*domain.User, error) {
	args := m.Called(login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) SaveUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Мок для AuthService
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

func TestUserService_FindUser(t *testing.T) {
	// Тестовые случаи
	tests := []struct {
		name        string
		login       string
		mockUser    *domain.User
		mockError   error
		expectedErr error
	}{
		{
			name:  "Успешный поиск пользователя",
			login: "testuser",
			mockUser: &domain.User{
				Credentials: domain.Credentials{
					Id:       1,
					Login:    "testuser",
					PassHash: "hashedpassword",
				},
			},
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:        "Пользователь не найден",
			login:       "nonexistentuser",
			mockUser:    nil,
			mockError:   domain.ErrNotFound,
			expectedErr: domain.ErrNotFound,
		},
		{
			name:        "Ошибка базы данных",
			login:       "testuser",
			mockUser:    nil,
			mockError:   errors.New("database error"),
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			mockRepo := new(MockUserRepo)
			mockAuthService := new(MockAuthService)

			// Настраиваем ожидаемое поведение мока
			mockRepo.On("FindUser", tt.login).Return(tt.mockUser, tt.mockError)

			// Создаем сервис с моками
			userService := NewUserService(mockRepo, mockAuthService)

			// Вызываем тестируемый метод
			user, err := userService.FindUser(tt.login)

			// Проверяем результаты
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUser, user)
			}

			// Проверяем, что все ожидаемые вызовы были выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_SaveUser(t *testing.T) {
	// Тестовые случаи
	tests := []struct {
		name           string
		login          string
		password       string
		mockHashedPass string
		mockHashErr    error
		mockSaveErr    error
		expectedUser   *domain.User
		expectedErr    error
	}{
		{
			name:           "Успешное сохранение пользователя",
			login:          "newuser",
			password:       "password123",
			mockHashedPass: "hashedpassword123",
			mockHashErr:    nil,
			mockSaveErr:    nil,
			expectedUser: &domain.User{
				Credentials: domain.Credentials{
					Login:    "newuser",
					PassHash: "hashedpassword123",
				},
			},
			expectedErr: nil,
		},
		{
			name:           "Ошибка при хешировании пароля",
			login:          "newuser",
			password:       "password123",
			mockHashedPass: "",
			mockHashErr:    errors.New("hash error"),
			mockSaveErr:    nil,
			expectedUser:   nil,
			expectedErr:    errors.New("hash error"),
		},
		{
			name:           "Ошибка при сохранении пользователя",
			login:          "newuser",
			password:       "password123",
			mockHashedPass: "hashedpassword123",
			mockHashErr:    nil,
			mockSaveErr:    errors.New("save error"),
			expectedUser:   nil,
			expectedErr:    errors.New("save error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			mockRepo := new(MockUserRepo)
			mockAuthService := new(MockAuthService)

			// Настраиваем ожидаемое поведение моков
			mockAuthService.On("HashPassword", tt.password).Return(tt.mockHashedPass, tt.mockHashErr)

			// Если нет ошибки хеширования, настраиваем ожидание для SaveUser
			if tt.mockHashErr == nil {
				// Используем mock.MatchedBy для проверки аргументов
				mockRepo.On("SaveUser", mock.MatchedBy(func(user *domain.User) bool {
					return user.Login == tt.login && user.PassHash == tt.mockHashedPass
				})).Return(tt.mockSaveErr)
			}

			// Создаем сервис с моками
			userService := NewUserService(mockRepo, mockAuthService)

			// Вызываем тестируемый метод
			user, err := userService.SaveUser(tt.login, tt.password)

			// Проверяем результаты
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.Login, user.Login)
				assert.Equal(t, tt.expectedUser.PassHash, user.PassHash)
			}

			// Проверяем, что все ожидаемые вызовы были выполнены
			mockRepo.AssertExpectations(t)
			mockAuthService.AssertExpectations(t)
		})
	}
}
