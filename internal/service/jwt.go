package service

import (
	"errors"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type JwtService struct {
}

// NewJwtService создает новый экземпляр JwtService
func NewJwtService() interfaces.JwtService {
	return &JwtService{}
}

func (j *JwtService) ExtractLoginFromToken(tokenString string) (string, error) {
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
