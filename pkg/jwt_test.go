package pkg

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestExtractLoginFromToken(t *testing.T) {
	// Создаем валидный JWT токен с полем login
	validToken := createValidJWTToken("testuser")
	validTokenWithBearer := "Bearer " + validToken

	// Создаем невалидный JWT токен
	invalidToken := "invalid-token"
	invalidTokenWithBearer := "Bearer " + invalidToken

	// Создаем валидный JWT токен без поля login
	tokenWithoutLogin := createJWTTokenWithoutLogin()
	tokenWithoutLoginWithBearer := "Bearer " + tokenWithoutLogin

	tests := []struct {
		name          string
		tokenString   string
		expectedLogin string
		expectError   bool
		errorMessage  string
	}{
		{
			name:          "Валидный токен с префиксом Bearer",
			tokenString:   validTokenWithBearer,
			expectedLogin: "testuser",
			expectError:   false,
			errorMessage:  "",
		},
		{
			name:          "Токен без префикса Bearer",
			tokenString:   validToken,
			expectedLogin: "",
			expectError:   true,
			errorMessage:  "токен должен начинаться с 'Bearer '",
		},
		{
			name:          "Невалидный токен с префиксом Bearer",
			tokenString:   invalidTokenWithBearer,
			expectedLogin: "",
			expectError:   true,
			errorMessage:  "ошибка при парсинге токена",
		},
		{
			name:          "Пустой токен",
			tokenString:   "",
			expectedLogin: "",
			expectError:   true,
			errorMessage:  "токен должен начинаться с 'Bearer '",
		},
		{
			name:          "Токен без поля login",
			tokenString:   tokenWithoutLoginWithBearer,
			expectedLogin: "",
			expectError:   true,
			errorMessage:  "поле 'login' отсутствует в токене или имеет неверный формат",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			login, err := ExtractLoginFromToken(tt.tokenString)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				assert.Empty(t, login)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLogin, login)
			}
		})
	}
}

// Вспомогательная функция для создания валидного JWT токена
func createValidJWTToken(login string) string {
	claims := jwt.MapClaims{
		"login": login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret-key"))
	return tokenString
}

// Вспомогательная функция для создания JWT токена без поля login
func createJWTTokenWithoutLogin() string {
	claims := jwt.MapClaims{
		"other_field": "some_value",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret-key"))
	return tokenString
}