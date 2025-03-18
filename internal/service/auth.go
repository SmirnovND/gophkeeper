package service

import (
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AuthService struct {
	cf interfaces.ConfigServer
}

func NewAuthService(cf interfaces.ConfigServer) *AuthService {
	return &AuthService{
		cf: cf,
	}
}

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func (a *AuthService) GenerateToken(login string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Генерируем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	return token.SignedString([]byte(a.cf.GetJwtSecret()))
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return a.cf.GetJwtSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем, является ли токен действительным
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// Хеширование пароля
func (a *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверка пароля
func (a *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *AuthService) SetResponseAuthData(w http.ResponseWriter, token string) {
	// Отправляем токен в заголовке Authorization
	w.Header().Set("Authorization", "Bearer "+token)

	// Установка cookie с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})
}
