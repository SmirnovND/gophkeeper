package service

import (
	"errors"
	"testing"
)

// MockTokenStorage - мок для интерфейса TokenStorage
type MockTokenStorage struct {
	SaveTokenFunc func(token string) error
	LoadTokenFunc func() (string, error)
}

func (m *MockTokenStorage) SaveToken(token string) error {
	if m.SaveTokenFunc != nil {
		return m.SaveTokenFunc(token)
	}
	return nil
}

func (m *MockTokenStorage) LoadToken() (string, error) {
	if m.LoadTokenFunc != nil {
		return m.LoadTokenFunc()
	}
	return "", nil
}

// TestNewTokenService проверяет создание нового экземпляра TokenService
func TestNewTokenService(t *testing.T) {
	mockStorage := &MockTokenStorage{}
	tokenService := NewTokenService(mockStorage)
	
	if tokenService == nil {
		t.Error("Ожидалось создание экземпляра TokenService, получен nil")
	}
}

// TestTokenService_SaveToken проверяет метод SaveToken
func TestTokenService_SaveToken(t *testing.T) {
	// Тест успешного сохранения токена
	t.Run("Success", func(t *testing.T) {
		tokenSaved := false
		expectedToken := "test-token"
		
		mockStorage := &MockTokenStorage{
			SaveTokenFunc: func(token string) error {
				tokenSaved = true
				if token != expectedToken {
					t.Errorf("Ожидался токен '%s', получен '%s'", expectedToken, token)
				}
				return nil
			},
		}
		
		tokenService := NewTokenService(mockStorage)
		tokenService.SaveToken(expectedToken)
		
		if !tokenSaved {
			t.Error("Метод SaveToken не был вызван")
		}
	})
}

// TestTokenService_LoadToken проверяет метод LoadToken
func TestTokenService_LoadToken(t *testing.T) {
	// Тест успешной загрузки токена
	t.Run("Success", func(t *testing.T) {
		expectedToken := "test-token"
		
		mockStorage := &MockTokenStorage{
			LoadTokenFunc: func() (string, error) {
				return expectedToken, nil
			},
		}
		
		tokenService := NewTokenService(mockStorage)
		token, err := tokenService.LoadToken()
		
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
		if token != expectedToken {
			t.Errorf("Ожидался токен '%s', получен '%s'", expectedToken, token)
		}
	})
	
	// Тест ошибки при загрузке токена
	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("ошибка загрузки токена")
		
		mockStorage := &MockTokenStorage{
			LoadTokenFunc: func() (string, error) {
				return "", expectedError
			},
		}
		
		tokenService := NewTokenService(mockStorage)
		token, err := tokenService.LoadToken()
		
		if err == nil {
			t.Error("Ожидалась ошибка, но ее не было")
		}
		if err != expectedError {
			t.Errorf("Ожидалась ошибка '%v', получена '%v'", expectedError, err)
		}
		if token != "" {
			t.Errorf("Ожидался пустой токен, получен '%s'", token)
		}
	})
}