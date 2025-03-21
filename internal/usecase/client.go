package usecase

import (
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
)

type ClientUseCase struct {
	TokenService  interfaces.TokenService
	ClientService interfaces.ClientService
}

func NewClientUseCase(
	TokenService interfaces.TokenService,
	ClientService interfaces.ClientService,
) interfaces.ClientUseCase {
	return &ClientUseCase{
		TokenService:  TokenService,
		ClientService: ClientService,
	}
}

func (c *ClientUseCase) Login(username string, password string) error {
	// Получаем токен через ClientService
	token, err := c.ClientService.Login(username, password)
	if err != nil {
		return fmt.Errorf("ошибка при входе: %w", err)
	}

	// Сохраняем полученный токен
	c.TokenService.SaveToken(username, token)
	return nil
}

func (c *ClientUseCase) Register(username string, password string, passwordCheck string) error {
	if password != passwordCheck {
		return fmt.Errorf("пароли не совпадают")
	}
	// Получаем токен через ClientService
	token, err := c.ClientService.Register(username, password)
	if err != nil {
		return fmt.Errorf("ошибка при регистрации: %w", err)
	}

	// Сохраняем полученный токен
	c.TokenService.SaveToken(username, token)
	return nil
}
