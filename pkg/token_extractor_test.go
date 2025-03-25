package pkg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenExtractor(t *testing.T) {
	// Сохраняем оригинальную функцию TokenExtractor
	originalTokenExtractor := TokenExtractor
	defer func() {
		// Восстанавливаем оригинальную функцию после теста
		TokenExtractor = originalTokenExtractor
	}()

	t.Run("Стандартная функция извлечения", func(t *testing.T) {
		// Создаем валидный JWT токен для тестирования
		validToken := "Bearer " + createValidJWTToken("testuser")
		
		// Вызываем TokenExtractor
		login, err := TokenExtractor(validToken)
		
		// Проверяем результат
		assert.NoError(t, err)
		assert.Equal(t, "testuser", login)
	})

	t.Run("Переопределение функции извлечения для тестов", func(t *testing.T) {
		// Переопределяем TokenExtractor для теста
		TokenExtractor = func(tokenString string) (string, error) {
			if tokenString == "Bearer test-token" {
				return "mock-user", nil
			}
			return "", errors.New("неверный токен")
		}

		// Тестируем успешный случай
		login, err := TokenExtractor("Bearer test-token")
		assert.NoError(t, err)
		assert.Equal(t, "mock-user", login)

		// Тестируем случай с ошибкой
		login, err = TokenExtractor("Bearer invalid-token")
		assert.Error(t, err)
		assert.Equal(t, "", login)
	})
}