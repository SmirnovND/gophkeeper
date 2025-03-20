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

func (t *TokenService) SaveToken(name string, token string) {
	t.ts.SaveToken(name, token)
}
