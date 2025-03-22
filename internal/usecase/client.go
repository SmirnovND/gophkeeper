package usecase

import (
	"errors"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/gophkeeper/pkg"
	"os"
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

// Upload - функция для загрузки файла на сервер.
func (c *ClientUseCase) Upload(filePath string, label string) (string, error) {
	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при открытии файла: %v\n", err))
	}
	defer file.Close()

	// Получение ссылки на загрузку файла
	url, err := c.ClientService.GetUploadLink(label, pkg.GetExtensionByPath(filePath))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при получении ссылки на загрузку: %v\n", err))
	}

	return c.ClientService.SendFileToServer(url, file)
}
