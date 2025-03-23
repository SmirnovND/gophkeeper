package pkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// ExtractLoginFromToken извлекает значение поля login из JWT токена.
// Токен должен быть в формате "Bearer <jwt-token>".
func ExtractLoginFromToken(tokenString string) (string, error) {
	// Проверяем, что токен начинается с "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return "", errors.New("токен должен начинаться с 'Bearer '")
	}

	// Удаляем префикс "Bearer "
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Парсим JWT токен
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("ошибка при парсинге токена: %w", err)
	}

	// Извлекаем claims из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("не удалось извлечь claims из токена")
	}

	// Извлекаем значение поля login
	login, ok := claims["login"].(string)
	if !ok {
		return "", errors.New("поле 'login' отсутствует в токене или имеет неверный формат")
	}

	return login, nil
}
