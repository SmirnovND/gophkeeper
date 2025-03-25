// +build test

package pkg

import (
	"errors"
	"strings"
)

// ExtractLoginFromToken - мок-версия функции для тестов
// Эта версия функции используется только в тестах.
// Токен должен быть в формате "Bearer <jwt-token>".
func ExtractLoginFromToken(tokenString string) (string, error) {
	// Проверяем, что токен начинается с "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return "", errors.New("токен должен начинаться с 'Bearer '")
	}

	// Удаляем префикс "Bearer "
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Для тестов просто извлекаем логин из токена, созданного функцией createValidJWTToken
	if strings.Contains(tokenString, "valid-token-for-") {
		login := strings.TrimPrefix(tokenString, "valid-token-for-")
		return login, nil
	}

	return "", errors.New("неверный токен")
}