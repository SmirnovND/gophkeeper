package usecase

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/service"
	"github.com/SmirnovND/toolbox/pkg/db"
	"net/http"
)

type AuthUseCase struct {
	userService        *service.UserService
	authService        *service.AuthService
	transactionManager *db.TransactionManager
}

func NewAuthUseCase(
	UserService *service.UserService,
	AuthService *service.AuthService,
	TransactionManager *db.TransactionManager,
) *AuthUseCase {
	return &AuthUseCase{
		userService:        UserService,
		authService:        AuthService,
		transactionManager: TransactionManager,
	}
}

func (a *AuthUseCase) Register(w http.ResponseWriter, credentials *domain.Credentials) {
	w.Header().Set("Content-Type", "application/json")

	_, err := a.userService.FindUser(credentials.Login)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err != domain.ErrNotFound {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	var user *domain.User

	user, err = a.userService.SaveUser(credentials.Login, credentials.Password)
	if err != nil {
		// Обработка ошибки сохранения пользователя
		http.Error(w, fmt.Sprintf("Error saving user: %v", err), http.StatusInternalServerError)
		return
	}

	token, err := a.authService.GenerateToken(user.Login)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating token: %v", err), http.StatusInternalServerError)
		return
	}

	a.authService.SetResponseAuthData(w, token)

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}

func (a *AuthUseCase) Login(w http.ResponseWriter, credentials *domain.Credentials) {
	w.Header().Set("Content-Type", "application/json")

	user, err := a.userService.FindUser(credentials.Login)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "Error", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
	}

	// Генерируем токен
	passValid := a.authService.CheckPasswordHash(credentials.Password, user.PassHash)
	if !passValid {
		http.Error(w, "Error", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := a.authService.GenerateToken(user.Login)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	a.authService.SetResponseAuthData(w, token)

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}
