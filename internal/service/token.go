package service

import "github.com/SmirnovND/gophkeeper/internal/interfaces"

type TokenService struct {
	ts interfaces.TokenStorage
}

func NewTokenService(ts interfaces.TokenStorage) interfaces.TokenService {
	return &TokenService{
		ts: ts,
	}
}

func (t *TokenService) SaveToken(token string) {
	t.ts.SaveToken(token)
}

func (t *TokenService) LoadToken() (string, error) {
	return t.ts.LoadToken()
}
