package service

import (
	"encoding/json"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тестирование методов для работы с текстовыми данными
func TestClientService_TextData(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок Authorization
		if r.Header.Get("Authorization") != "test-token" {
			t.Errorf("Ожидался Authorization test-token, получен %s", r.Header.Get("Authorization"))
		}

		// Проверяем заголовок Content-Type
		if r.Method == "POST" && r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Ожидался Content-Type application/json, получен %s", r.Header.Get("Content-Type"))
		}

		// Обрабатываем разные методы и пути
		if r.Method == "POST" && r.URL.Path == "/api/data/text/test-text" {
			// Сохранение текстовых данных
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			var requestData struct {
				TextData *domain.TextData `json:"text_data"`
				Metadata string           `json:"metadata"`
			}
			if err := json.Unmarshal(body, &requestData); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if requestData.TextData == nil || requestData.TextData.Content != "test text content" {
				content := ""
				if requestData.TextData != nil {
					content = requestData.TextData.Content
				}
				t.Errorf("Ожидался текст 'test text content', получен '%s'", content)
			}

			w.WriteHeader(http.StatusOK)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/text/test-text" {
			// Получение текстовых данных
			response := struct {
				TextData *domain.TextData `json:"text_data"`
				Metadata string           `json:"metadata"`
			}{
				TextData: &domain.TextData{
					Content: "test text content",
				},
				Metadata: "",
			}
			responseJSON, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
		} else if r.Method == "DELETE" && r.URL.Path == "/api/data/text/test-text" {
			// Удаление текстовых данных
			w.WriteHeader(http.StatusOK)
		} else {
			t.Errorf("Неожиданный запрос: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Создаем клиентский сервис с адресом тестового сервера
	serverAddr := server.URL[7:] // Убираем "http://" из URL
	clientService := NewClientService(serverAddr)

	// Тестируем SaveText
	textData := &domain.TextData{
		Content: "test text content",
	}
	err := clientService.SaveText("test-text", textData, "", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveText: %v", err)
	}

	// Тестируем GetText
	retrievedTextData, _, err := clientService.GetText("test-text", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове GetText: %v", err)
	}
	if retrievedTextData.Content != "test text content" {
		t.Errorf("Ожидался текст 'test text content', получен '%s'", retrievedTextData.Content)
	}

	// Тестируем DeleteText
	err = clientService.DeleteText("test-text", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове DeleteText: %v", err)
	}
}
