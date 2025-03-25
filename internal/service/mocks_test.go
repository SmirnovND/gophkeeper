package service

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"net/http"
)

// MockUserRepo - мок для интерфейса UserRepo
type MockUserRepo struct {
	FindUserFunc func(login string) (*domain.User, error)
	SaveUserFunc func(user *domain.User) error
}

// FindUser - реализация метода FindUser для мока
func (m *MockUserRepo) FindUser(login string) (*domain.User, error) {
	return m.FindUserFunc(login)
}

// SaveUser - реализация метода SaveUser для мока
func (m *MockUserRepo) SaveUser(user *domain.User) error {
	return m.SaveUserFunc(user)
}

// MockAuthService - мок для интерфейса AuthService
type MockAuthService struct {
	GenerateTokenFunc     func(login string) (string, error)
	ValidateTokenFunc     func(tokenString string) (*domain.Claims, error)
	HashPasswordFunc      func(password string) (string, error)
	CheckPasswordHashFunc func(password, hash string) bool
	SetResponseAuthDataFunc func(w http.ResponseWriter, token string)
}

// GenerateToken - реализация метода GenerateToken для мока
func (m *MockAuthService) GenerateToken(login string) (string, error) {
	return m.GenerateTokenFunc(login)
}

// ValidateToken - реализация метода ValidateToken для мока
func (m *MockAuthService) ValidateToken(tokenString string) (*domain.Claims, error) {
	return m.ValidateTokenFunc(tokenString)
}

// HashPassword - реализация метода HashPassword для мока
func (m *MockAuthService) HashPassword(password string) (string, error) {
	return m.HashPasswordFunc(password)
}

// CheckPasswordHash - реализация метода CheckPasswordHash для мока
func (m *MockAuthService) CheckPasswordHash(password, hash string) bool {
	return m.CheckPasswordHashFunc(password, hash)
}

// SetResponseAuthData - реализация метода SetResponseAuthData для мока
func (m *MockAuthService) SetResponseAuthData(w http.ResponseWriter, token string) {
	m.SetResponseAuthDataFunc(w, token)
}

// MockUserDataRepo - мок для интерфейса UserDataRepo
type MockUserDataRepo struct {
	SaveUserDataFunc              func(userData *domain.UserData) error
	FindUserDataByLabelFunc       func(userID, label string) (*domain.UserData, error)
	GetUserDataByLabelAndTypeFunc func(userID, label string, dataType string) (*domain.UserData, error)
	DeleteUserDataFunc            func(id string) error
}

// SaveUserData - реализация метода SaveUserData для мока
func (m *MockUserDataRepo) SaveUserData(userData *domain.UserData) error {
	return m.SaveUserDataFunc(userData)
}

// FindUserDataByLabel - реализация метода FindUserDataByLabel для мока
func (m *MockUserDataRepo) FindUserDataByLabel(userID, label string) (*domain.UserData, error) {
	return m.FindUserDataByLabelFunc(userID, label)
}

// GetUserDataByLabelAndType - реализация метода GetUserDataByLabelAndType для мока
func (m *MockUserDataRepo) GetUserDataByLabelAndType(userID, label string, dataType string) (*domain.UserData, error) {
	return m.GetUserDataByLabelAndTypeFunc(userID, label, dataType)
}

// DeleteUserData - реализация метода DeleteUserData для мока
func (m *MockUserDataRepo) DeleteUserData(id string) error {
	return m.DeleteUserDataFunc(id)
}