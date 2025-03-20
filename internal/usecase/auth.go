package usecase

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"net/http"
)

type AuthUseCase struct {
	userService interfaces.UserService
	authService interfaces.AuthService
}

func NewAuthUseCase(
	UserService interfaces.UserService,
	AuthService interfaces.AuthService,
) interfaces.AuthUseCase {
	return &AuthUseCase{
		userService: UserService,
		authService: AuthService,
	}
}

func (a *AuthUseCase) Register(w http.ResponseWriter, credentials *domain.Credentials) (string, error) {
	w.Header().Set("Content-Type", "application/json")

	_, err := a.userService.FindUser(credentials.Login)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		return "", fmt.Errorf("user already exists")
	}

	if err != domain.ErrNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", fmt.Errorf(err.Error()+" error finding user: %w", err)
	}

	var user *domain.User

	user, err = a.userService.SaveUser(credentials.Login, credentials.Password)
	if err != nil {
		// Обработка ошибки сохранения пользователя
		http.Error(w, fmt.Sprintf("Error saving user: %v", err), http.StatusInternalServerError)
		return "", fmt.Errorf("error saving user: %w", err)
	}

	token, err := a.authService.GenerateToken(user.Login)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating token: %v", err), http.StatusInternalServerError)
		return "", fmt.Errorf("error generating token: %w", err)
	}

	a.authService.SetResponseAuthData(w, token)

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))

	return token, nil
}

func (a *AuthUseCase) Login(w http.ResponseWriter, credentials *domain.Credentials) (string, error) {
	w.Header().Set("Content-Type", "application/json")

	user, err := a.userService.FindUser(credentials.Login)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "Error: user not found", http.StatusUnauthorized)
			return "", fmt.Errorf("user not found")
		} else {
			http.Error(w, "Error: error finding user", http.StatusInternalServerError)
			return "", fmt.Errorf("error finding user: %w", err)
		}
	}

	// Проверяем пароль
	passValid := a.authService.CheckPasswordHash(credentials.Password, user.PassHash)
	if !passValid {
		http.Error(w, "Error: invalid password", http.StatusUnauthorized)
		return "", fmt.Errorf("invalid password")
	}

	// Генерируем токен
	token, err := a.authService.GenerateToken(user.Login)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return "", fmt.Errorf("error generating token: %w", err)
	}

	a.authService.SetResponseAuthData(w, token)

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))

	return token, nil
}

// ValidateToken проверяет валидность JWT токена и возвращает claims
func (a *AuthUseCase) ValidateToken(token string) (*domain.Claims, error) {
	return a.authService.ValidateToken(token)
}
