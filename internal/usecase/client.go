package usecase

import (
	"errors"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/gophkeeper/pkg"
	"os"
	"path/filepath"
	"strings"
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
	// Проверяем, существует ли файл
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при проверке файла: %v\n", err))
	}

	// Проверяем, что это файл, а не директория
	if fileInfo.IsDir() {
		return "", errors.New("Указанный путь является директорией, а не файлом")
	}

	// Проверяем тип файла (текстовый или бинарный)
	isText := isTextFile(filePath)
	isBinary := isBinaryFile(filePath)

	// Проверяем, что файл является текстовым или бинарным
	if !isText && !isBinary {
		return "", errors.New("Файл не является ни текстовым, ни бинарным")
	}

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

	// Выводим информацию о типе файла
	fileType := "бинарный"
	if isText {
		fileType = "текстовый"
	}
	fmt.Printf("Загрузка %s файла: %s\n", fileType, filePath)

	return c.ClientService.SendFileToServer(url, file)
}

// isTextFile проверяет, является ли файл текстовым
func isTextFile(filePath string) bool {
	// Проверка по расширению файла (быстрый метод)
	if ".txt" == strings.ToLower(filepath.Ext(filePath)) {
		return true
	}

	return false
}

// isBinaryFile проверяет, является ли файл бинарным
func isBinaryFile(filePath string) bool {
	// Проверка по расширению файла (быстрый метод)
	if ".bin" == strings.ToLower(filepath.Ext(filePath)) {
		return true
	}

	return false
}
