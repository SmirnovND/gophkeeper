package repo

import "sync"

// TokenStorage - структура для хранения JWT-токена в памяти.
type TokenStorage struct {
	mu    sync.RWMutex
	token string
}

// NewTokenStorage создает новый экземпляр TokenStorage.
func NewTokenStorage() *TokenStorage {
	return &TokenStorage{}
}

// SaveToken сохраняет токен в памяти.
func (s *TokenStorage) SaveToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
}

// GetToken возвращает сохраненный токен.
func (s *TokenStorage) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token
}

// ClearToken очищает токен в памяти.
func (s *TokenStorage) ClearToken() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = ""
}
