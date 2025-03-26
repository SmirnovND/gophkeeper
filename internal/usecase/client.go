package usecase

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
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
	c.TokenService.SaveToken(token)
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
	c.TokenService.SaveToken(token)
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

	token, err := c.TokenService.LoadToken()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при загрузке токена: %v\n", err))
	}
	// Запрашиваем метаинформацию у пользователя
	fmt.Println("Введите метаинформацию для файла (необязательно):")
	fmt.Print("> ")
	var metadata string
	reader := bufio.NewReader(os.Stdin)
	metadata, _ = reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	// Получение ссылки на загрузку файла
	url, err := c.ClientService.GetUploadLink(label, pkg.GetExtensionByPath(filePath), metadata, token)
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

// Download - функция для скачивания файла с сервера.
func (c *ClientUseCase) Download(label string) error {
	// Проверяем, что метка файла указана
	if label == "" {
		return errors.New("Не указана метка файла")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Получаем ссылку на скачивание файла, метаданные и метаинформацию
	downloadURL, fileMetadata, metaInfo, err := c.ClientService.GetDownloadLink(label, token)
	if err != nil {
		return fmt.Errorf("ошибка при получении ссылки на скачивание: %w", err)
	}
	
	// Выводим метаинформацию, если она есть
	if metaInfo != "" {
		fmt.Println("Метаинформация файла:")
		fmt.Println("------------------")
		fmt.Println(metaInfo)
		fmt.Println("------------------")
	}

	// Получаем директорию загрузок
	downloadsDir := pkg.GetDownloadsDir()
	// Формируем имя файла из метки и расширения
	fileName := fmt.Sprintf("%s.%s", label, fileMetadata.Extension)
	outputPath := filepath.Join(downloadsDir, fileName)

	fmt.Printf("Скачивание файла с меткой '%s'\n", label)

	// Скачиваем файл
	err = c.ClientService.DownloadFileFromServer(downloadURL, outputPath)
	if err != nil {
		return fmt.Errorf("ошибка при скачивании файла: %w", err)
	}

	fmt.Printf("Файл успешно скачан и сохранен в '%s'\n", outputPath)
	return nil
}

// SaveText сохраняет текстовые данные
func (c *ClientUseCase) SaveText(label string, textData *domain.TextData, metadata string) error {
	// Проверяем, что метка указана
	if label == "" {
		return errors.New("не указана метка для текстовых данных")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Сохраняем текстовые данные
	err = c.ClientService.SaveText(label, textData, metadata, token)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении текстовых данных: %w", err)
	}

	return nil
}

// GetText получает текстовые данные
func (c *ClientUseCase) GetText(label string) (*domain.TextData, string, error) {
	// Проверяем, что метка указана
	if label == "" {
		return nil, "", errors.New("не указана метка для текстовых данных")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Получаем текстовые данные
	textData, metadata, err := c.ClientService.GetText(label, token)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при получении текстовых данных: %w", err)
	}

	return textData, metadata, nil
}

// DeleteText удаляет текстовые данные
func (c *ClientUseCase) DeleteText(label string) error {
	// Проверяем, что метка указана
	if label == "" {
		return errors.New("не указана метка для текстовых данных")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Удаляем текстовые данные
	err = c.ClientService.DeleteText(label, token)
	if err != nil {
		return fmt.Errorf("ошибка при удалении текстовых данных: %w", err)
	}

	return nil
}

// SaveCard сохраняет данные кредитной карты
func (c *ClientUseCase) SaveCard(label string, cardData *domain.CardData, metadata string) error {
	// Проверяем, что метка указана
	if label == "" {
		return errors.New("не указана метка для данных карты")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Сохраняем данные карты
	err = c.ClientService.SaveCard(label, cardData, metadata, token)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении данных карты: %w", err)
	}

	return nil
}

// GetCard получает данные кредитной карты
func (c *ClientUseCase) GetCard(label string) (*domain.CardData, string, error) {
	// Проверяем, что метка указана
	if label == "" {
		return nil, "", errors.New("не указана метка для данных карты")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Получаем данные карты
	cardData, metadata, err := c.ClientService.GetCard(label, token)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при получении данных карты: %w", err)
	}

	return cardData, metadata, nil
}

// DeleteCard удаляет данные кредитной карты
func (c *ClientUseCase) DeleteCard(label string) error {
	// Проверяем, что метка указана
	if label == "" {
		return errors.New("не указана метка для данных карты")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Удаляем данные карты
	err = c.ClientService.DeleteCard(label, token)
	if err != nil {
		return fmt.Errorf("ошибка при удалении данных карты: %w", err)
	}

	return nil
}

// SaveCredential сохраняет учетные данные
func (c *ClientUseCase) SaveCredential(label string, credentialData *domain.CredentialData, metadata string) error {
	// Проверяем, что метка указана
	if label == "" {
		return errors.New("не указана метка для учетных данных")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Сохраняем учетные данные
	err = c.ClientService.SaveCredential(label, credentialData, metadata, token)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении учетных данных: %w", err)
	}

	return nil
}

// GetCredential получает учетные данные
func (c *ClientUseCase) GetCredential(label string) (*domain.CredentialData, string, error) {
	// Проверяем, что метка указана
	if label == "" {
		return nil, "", errors.New("не указана метка для учетных данных")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Получаем учетные данные
	credentialData, metadata, err := c.ClientService.GetCredential(label, token)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при получении учетных данных: %w", err)
	}

	return credentialData, metadata, nil
}

// DeleteCredential удаляет учетные данные
func (c *ClientUseCase) DeleteCredential(label string) error {
	// Проверяем, что метка указана
	if label == "" {
		return errors.New("не указана метка для учетных данных")
	}

	// Загружаем токен
	token, err := c.TokenService.LoadToken()
	if err != nil {
		return fmt.Errorf("ошибка при загрузке токена: %w", err)
	}

	// Удаляем учетные данные
	err = c.ClientService.DeleteCredential(label, token)
	if err != nil {
		return fmt.Errorf("ошибка при удалении учетных данных: %w", err)
	}

	return nil
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
