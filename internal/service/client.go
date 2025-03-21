package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"net/http"
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

func (c *ClientService) Login(login, password string) (string, error) {
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

func (c *ClientService) FindUser(login string) (*domain.User, error) {
	return nil, domain.ErrNotFound
}

func (c *ClientService) SaveUser(login, password string) (*domain.User, error) {
	return nil, fmt.Errorf("метод не реализован на клиенте")
}
