package service

import (
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"testing"
)

// TestNewUserService тестирует функцию NewUserService
func TestNewUserService(t *testing.T) {
	// Создаем моки для репозитория и сервиса аутентификации
	mockUserRepo := &MockUserRepo{}
	mockAuthService := &MockAuthService{}

	// Вызываем функцию NewUserService
	userService := NewUserService(mockUserRepo, mockAuthService)

	// Проверяем, что возвращенный объект не nil
	if userService == nil {
		t.Fatal("Функция NewUserService вернула nil")
	}

	// Проверяем, что возвращенный объект реализует интерфейс UserService
	_, ok := userService.(interfaces.UserService)
	if !ok {
		t.Fatal("Возвращенный объект не реализует интерфейс UserService")
	}
}

// TestUserService_FindUser тестирует метод FindUser
func TestUserService_FindUser(t *testing.T) {
	// Создаем мок для репозитория
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			// Проверяем параметры
			if login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}
			return &domain.User{
				Id: "user123",
				Credentials: domain.Credentials{
					Login:    "testuser",
					PassHash: "hash",
				},
			}, nil
		},
	}

	// Создаем мок для сервиса аутентификации
	mockAuthService := &MockAuthService{}

	// Создаем экземпляр UserService
	userService := &UserService{
		repo:        mockUserRepo,
		authService: mockAuthService,
	}

	// Вызываем метод FindUser
	user, err := userService.FindUser("testuser")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове FindUser: %v", err)
	}
	if user == nil {
		t.Fatal("Результат не должен быть nil")
	}
	if user.Login != "testuser" {
		t.Errorf("Ожидался логин 'testuser', получен '%s'", user.Login)
	}
	if user.Id != "user123" {
		t.Errorf("Ожидался ID 'user123', получен '%s'", user.Id)
	}
}

// TestUserService_FindUser_NotFound тестирует метод FindUser с ошибкой "пользователь не найден"
func TestUserService_FindUser_NotFound(t *testing.T) {
	// Создаем мок для репозитория
	mockUserRepo := &MockUserRepo{
		FindUserFunc: func(login string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	// Создаем мок для сервиса аутентификации
	mockAuthService := &MockAuthService{}

	// Создаем экземпляр UserService
	userService := &UserService{
		repo:        mockUserRepo,
		authService: mockAuthService,
	}

	// Вызываем метод FindUser
	_, err := userService.FindUser("testuser")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
	if err != domain.ErrNotFound {
		t.Errorf("Ожидалась ошибка domain.ErrNotFound, получена '%v'", err)
	}
}

// TestUserService_SaveUser тестирует метод SaveUser
func TestUserService_SaveUser(t *testing.T) {
	// Создаем мок для репозитория
	mockUserRepo := &MockUserRepo{
		SaveUserFunc: func(user *domain.User) error {
			// Проверяем параметры
			if user.Login != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", user.Login)
			}
			if user.PassHash != "hashed_password" {
				t.Errorf("Ожидался хеш пароля 'hashed_password', получен '%s'", user.PassHash)
			}

			// Устанавливаем ID пользователя
			user.Id = "user123"
			return nil
		},
	}

	// Создаем мок для сервиса аутентификации
	mockAuthService := &MockAuthService{
		HashPasswordFunc: func(password string) (string, error) {
			// Проверяем параметры
			if password != "password" {
				t.Errorf("Ожидался пароль 'password', получен '%s'", password)
			}
			return "hashed_password", nil
		},
	}

	// Создаем экземпляр UserService
	userService := &UserService{
		repo:        mockUserRepo,
		authService: mockAuthService,
	}

	// Вызываем метод SaveUser
	user, err := userService.SaveUser("testuser", "password")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveUser: %v", err)
	}
	if user == nil {
		t.Fatal("Результат не должен быть nil")
	}
	if user.Login != "testuser" {
		t.Errorf("Ожидался логин 'testuser', получен '%s'", user.Login)
	}
	if user.PassHash != "hashed_password" {
		t.Errorf("Ожидался хеш пароля 'hashed_password', получен '%s'", user.PassHash)
	}
	if user.Id != "user123" {
		t.Errorf("Ожидался ID 'user123', получен '%s'", user.Id)
	}
}

// TestUserService_SaveUser_HashError тестирует метод SaveUser с ошибкой хеширования пароля
func TestUserService_SaveUser_HashError(t *testing.T) {
	// Создаем мок для репозитория
	mockUserRepo := &MockUserRepo{}

	// Создаем мок для сервиса аутентификации
	mockAuthService := &MockAuthService{
		HashPasswordFunc: func(password string) (string, error) {
			return "", errors.New("ошибка хеширования пароля")
		},
	}

	// Создаем экземпляр UserService
	userService := &UserService{
		repo:        mockUserRepo,
		authService: mockAuthService,
	}

	// Вызываем метод SaveUser
	_, err := userService.SaveUser("testuser", "password")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// TestUserService_SaveUser_SaveError тестирует метод SaveUser с ошибкой сохранения пользователя
func TestUserService_SaveUser_SaveError(t *testing.T) {
	// Создаем мок для репозитория
	mockUserRepo := &MockUserRepo{
		SaveUserFunc: func(user *domain.User) error {
			return errors.New("ошибка сохранения пользователя")
		},
	}

	// Создаем мок для сервиса аутентификации
	mockAuthService := &MockAuthService{
		HashPasswordFunc: func(password string) (string, error) {
			return "hashed_password", nil
		},
	}

	// Создаем экземпляр UserService
	userService := &UserService{
		repo:        mockUserRepo,
		authService: mockAuthService,
	}

	// Вызываем метод SaveUser
	_, err := userService.SaveUser("testuser", "password")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}