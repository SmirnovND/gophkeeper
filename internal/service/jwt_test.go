package service

import (
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

// TestNewJwtService проверяет создание нового экземпляра JwtService
func TestNewJwtService(t *testing.T) {
	jwtService := NewJwtService()
	if jwtService == nil {
		t.Error("Ожидалось создание экземпляра JwtService, получен nil")
	}
}

// TestExtractLoginFromToken проверяет извлечение логина из токена
func TestExtractLoginFromToken(t *testing.T) {
	jwtService := NewJwtService()

	// Тест успешного извлечения логина
	t.Run("Success", func(t *testing.T) {
		// Создаем тестовый токен
		claims := jwt.MapClaims{
			"login": "testuser",
			"exp":   time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("test_secret"))
		if err != nil {
			t.Fatalf("Ошибка при создании токена: %v", err)
		}

		// Добавляем префикс Bearer
		tokenWithBearer := "Bearer " + tokenString

		// Извлекаем логин
		login, err := jwtService.ExtractLoginFromToken(tokenWithBearer)
		if err != nil {
			t.Errorf("Не ожидалась ошибка, получена: %v", err)
		}
		if login != "testuser" {
			t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
		}
	})

	// Тест ошибки при отсутствии префикса Bearer
	t.Run("NoBearerPrefix", func(t *testing.T) {
		// Создаем тестовый токен без префикса Bearer
		claims := jwt.MapClaims{
			"login": "testuser",
			"exp":   time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("test_secret"))
		if err != nil {
			t.Fatalf("Ошибка при создании токена: %v", err)
		}

		// Извлекаем логин
		_, err = jwtService.ExtractLoginFromToken(tokenString)
		if err == nil {
			t.Error("Ожидалась ошибка отсутствия префикса Bearer, но ее не было")
		}
	})

	// Тест ошибки при неверном формате токена
	t.Run("InvalidTokenFormat", func(t *testing.T) {
		// Используем неверный формат токена
		tokenWithBearer := "Bearer invalid.token.format"

		// Извлекаем логин
		_, err := jwtService.ExtractLoginFromToken(tokenWithBearer)
		if err == nil {
			t.Error("Ожидалась ошибка неверного формата токена, но ее не было")
		}
	})

	// Тест ошибки при отсутствии поля login в токене
	t.Run("NoLoginField", func(t *testing.T) {
		// Создаем тестовый токен без поля login
		claims := jwt.MapClaims{
			"user_id": "123",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("test_secret"))
		if err != nil {
			t.Fatalf("Ошибка при создании токена: %v", err)
		}

		// Добавляем префикс Bearer
		tokenWithBearer := "Bearer " + tokenString

		// Извлекаем логин
		_, err = jwtService.ExtractLoginFromToken(tokenWithBearer)
		if err == nil {
			t.Error("Ожидалась ошибка отсутствия поля login, но ее не было")
		}
	})
}