package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type ClientService struct {
	client     *http.Client
	serverAddr string
}

func NewClientService(serverAddr string) interfaces.ClientService {
	return &ClientService{
		client:     &http.Client{},
		serverAddr: serverAddr,
	}
}

func (c *ClientService) sendRequest(method, url string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при маршалинге данных: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}

	return resp, nil
}

func (c *ClientService) Login(login string, password string) (string, error) {
	credentials := domain.Credentials{Login: login, Password: password}
	resp, err := c.sendRequest("POST", "http://"+c.serverAddr+"/api/user/login", credentials)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка аутентификации, код ответа: %d", resp.StatusCode)
	}

	token := resp.Header.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("токен не найден в ответе")
	}

	return token, nil
}

func (c *ClientService) Register(login, password string) (string, error) {
	credentials := domain.Credentials{Login: login, Password: password}
	resp, err := c.sendRequest("POST", "http://"+c.serverAddr+"/api/user/register", credentials)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return "", fmt.Errorf("пользователь с таким логином уже существует")
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка регистрации, код ответа: %d", resp.StatusCode)
	}

	token := resp.Header.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("токен не найден в ответе")
	}

	return token, nil
}

func (c *ClientService) GetUploadLink(label string, extension string, token string) (string, error) {
	// Запрос на получение ссылки для загрузки файла
	url := "http://" + c.serverAddr + "/api/file/upload"

	// Создаем структуру для запроса с параметрами из аргументов функции
	requestData := struct {
		Name      string `json:"name"`
		Extension string `json:"extension"`
	}{
		Name:      label,
		Extension: extension,
	}

	// Преобразуем структуру в JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("ошибка при маршалинге данных: %w", err)
	}

	// Создаем запрос вместо использования http.Post для добавления заголовка авторизации
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// Выполняем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при запросе к серверу: %v\n", err))
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка при получении ссылки для загрузки, код ответа: %d", resp.StatusCode)
	}

	// Чтение ответа сервера
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при чтении ответа сервера: %v\n", err))
	}

	// Извлекаем URL из ответа
	var response struct {
		URL         string `json:"url"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при парсинге ответа: %v\n", err))
	}

	return response.URL, nil
}

func (c *ClientService) SendFileToServer(url string, file *os.File) (string, error) {
	// Получаем информацию о файле для определения его размера
	fileInfo, err := file.Stat()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при получении информации о файле: %v\n", err))
	}

	// Сбрасываем указатель чтения файла в начало
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при перемещении указателя файла: %v\n", err))
	}

	// Загрузка файла
	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при подготовке запроса на загрузку: %v\n", err))
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = fileInfo.Size() // Устанавливаем размер файла в заголовке Content-Length

	client := &http.Client{}
	fileUploadResp, err := client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка при загрузке файла: %v\n", err))
	}
	defer fileUploadResp.Body.Close()

	if fileUploadResp.StatusCode == http.StatusOK {
		return "Файл успешно загружен!", nil
	} else {
		return "", errors.New(fmt.Sprintf("Ошибка при загрузке файла: %s\n", fileUploadResp.Status))
	}
}

func (c *ClientService) GetDownloadLink(label string, token string) (string, error) {
	// Формируем URL для запроса на получение ссылки для скачивания
	url := fmt.Sprintf("http://%s/api/file/download?label=%s", c.serverAddr, label)

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Устанавливаем заголовок авторизации
	req.Header.Set("Authorization", token)

	// Выполняем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка при получении ссылки для скачивания, код ответа: %d", resp.StatusCode)
	}

	// Чтение ответа сервера
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа сервера: %w", err)
	}

	// Извлекаем URL из ответа
	var response struct {
		URL         string `json:"url"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("ошибка при парсинге ответа: %w", err)
	}

	return response.URL, nil
}

func (c *ClientService) DownloadFileFromServer(url string, outputPath string) error {
	// Создаем запрос на скачивание файла
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса на скачивание: %w", err)
	}

	// Выполняем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса на скачивание: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка при скачивании файла, код ответа: %d", resp.StatusCode)
	}

	// Создаем файл для сохранения скачанных данных
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла для сохранения: %w", err)
	}
	defer outputFile.Close()

	// Копируем данные из ответа в файл
	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении файла: %w", err)
	}

	return nil
}
