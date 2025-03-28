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

// Тестирование методов для работы с данными кредитных карт
func TestClientService_CardData(t *testing.T) {
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
		if r.Method == "POST" && r.URL.Path == "/api/data/card/test-card" {
			// Сохранение данных карты
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			var requestData struct {
				CardData *domain.CardData `json:"card_data"`
				Metadata string           `json:"metadata"`
			}
			if err := json.Unmarshal(body, &requestData); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if requestData.CardData == nil || requestData.CardData.Number != "1234567890123456" {
				number := ""
				if requestData.CardData != nil {
					number = requestData.CardData.Number
				}
				t.Errorf("Ожидался номер карты '1234567890123456', получен '%s'", number)
			}

			w.WriteHeader(http.StatusOK)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/card/test-card" {
			// Получение данных карты
			response := struct {
				CardData *domain.CardData `json:"card_data"`
				Metadata string           `json:"metadata"`
			}{
				CardData: &domain.CardData{
					Number:     "1234567890123456",
					Holder:     "Test User",
					ExpiryDate: "12/25",
					CVV:        "123",
				},
				Metadata: "test metadata",
			}
			responseJSON, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
		} else if r.Method == "DELETE" && r.URL.Path == "/api/data/card/test-card" {
			// Удаление данных карты
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

	// Тестируем SaveCard
	cardData := &domain.CardData{
		Number:     "1234567890123456",
		Holder:     "Test User",
		ExpiryDate: "12/25",
		CVV:        "123",
	}
	err := clientService.SaveCard("test-card", cardData, "test metadata", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveCard: %v", err)
	}

	// Тестируем GetCard
	retrievedCardData, metadata, err := clientService.GetCard("test-card", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове GetCard: %v", err)
	}
	if retrievedCardData.Number != "1234567890123456" {
		t.Errorf("Ожидался номер карты '1234567890123456', получен '%s'", retrievedCardData.Number)
	}
	if metadata != "test metadata" {
		t.Errorf("Ожидались метаданные 'test metadata', получены '%s'", metadata)
	}

	// Тестируем DeleteCard
	err = clientService.DeleteCard("test-card", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове DeleteCard: %v", err)
	}
}

// Тестирование методов для работы с учетными данными
func TestClientService_CredentialData(t *testing.T) {
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
		if r.Method == "POST" && r.URL.Path == "/api/data/credential/test-credential" {
			// Сохранение учетных данных
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			var requestData struct {
				CredentialData *domain.CredentialData `json:"credential_data"`
				Metadata       string                 `json:"metadata"`
			}
			if err := json.Unmarshal(body, &requestData); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if requestData.CredentialData == nil || requestData.CredentialData.Login != "testuser" {
				login := ""
				if requestData.CredentialData != nil {
					login = requestData.CredentialData.Login
				}
				t.Errorf("Ожидался логин 'testuser', получен '%s'", login)
			}

			w.WriteHeader(http.StatusOK)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/credential/test-credential" {
			// Получение учетных данных
			response := struct {
				CredentialData *domain.CredentialData `json:"credential_data"`
				Metadata       string                 `json:"metadata"`
			}{
				CredentialData: &domain.CredentialData{
					Login:    "testuser",
					Password: "testpass",
				},
				Metadata: "test metadata",
			}
			responseJSON, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
		} else if r.Method == "DELETE" && r.URL.Path == "/api/data/credential/test-credential" {
			// Удаление учетных данных
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

	// Тестируем SaveCredential
	credentialData := &domain.CredentialData{
		Login:    "testuser",
		Password: "testpass",
	}
	err := clientService.SaveCredential("test-credential", credentialData, "test metadata", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveCredential: %v", err)
	}

	// Тестируем GetCredential
	retrievedCredentialData, metadata, err := clientService.GetCredential("test-credential", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове GetCredential: %v", err)
	}
	if retrievedCredentialData.Login != "testuser" {
		t.Errorf("Ожидался логин 'testuser', получен '%s'", retrievedCredentialData.Login)
	}
	if retrievedCredentialData.Password != "testpass" {
		t.Errorf("Ожидался пароль 'testpass', получен '%s'", retrievedCredentialData.Password)
	}
	if metadata != "test metadata" {
		t.Errorf("Ожидались метаданные 'test metadata', получены '%s'", metadata)
	}

	// Тестируем DeleteCredential
	err = clientService.DeleteCredential("test-credential", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове DeleteCredential: %v", err)
	}
}

// Тестирование методов для работы с аутентификацией
func TestClientService_Auth(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Обрабатываем разные методы и пути
		if r.Method == "POST" && r.URL.Path == "/api/user/login" {
			// Проверяем тело запроса
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			var credentials domain.Credentials
			if err := json.Unmarshal(body, &credentials); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if credentials.Login != "testuser" || credentials.Password != "testpass" {
				t.Errorf("Ожидались логин 'testuser' и пароль 'testpass', получены '%s' и '%s'", credentials.Login, credentials.Password)
			}

			// Устанавливаем заголовок Authorization
			w.Header().Set("Authorization", "test-token")
			w.WriteHeader(http.StatusOK)
		} else if r.Method == "POST" && r.URL.Path == "/api/user/register" {
			// Проверяем тело запроса
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			var credentials domain.Credentials
			if err := json.Unmarshal(body, &credentials); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if credentials.Login != "newuser" || credentials.Password != "newpass" {
				t.Errorf("Ожидались логин 'newuser' и пароль 'newpass', получены '%s' и '%s'", credentials.Login, credentials.Password)
			}

			// Устанавливаем заголовок Authorization
			w.Header().Set("Authorization", "new-token")
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

	// Тестируем Login
	token, err := clientService.Login("testuser", "testpass")
	if err != nil {
		t.Fatalf("Ошибка при вызове Login: %v", err)
	}
	if token != "test-token" {
		t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
	}

	// Тестируем Register
	token, err = clientService.Register("newuser", "newpass")
	if err != nil {
		t.Fatalf("Ошибка при вызове Register: %v", err)
	}
	if token != "new-token" {
		t.Errorf("Ожидался токен 'new-token', получен '%s'", token)
	}
}

// Тестирование методов для работы с файлами
func TestClientService_FileOperations(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок Authorization для запросов, требующих авторизации
		if r.URL.Path == "/api/file/upload" || r.URL.Path == "/api/file/download" {
			if r.Header.Get("Authorization") != "test-token" {
				t.Errorf("Ожидался Authorization test-token, получен %s", r.Header.Get("Authorization"))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		// Обрабатываем разные методы и пути
		if r.Method == "POST" && r.URL.Path == "/api/file/upload" {
			// Проверяем тело запроса
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			var requestData struct {
				Name      string `json:"name"`
				Extension string `json:"extension"`
				Metadata  string `json:"metadata"`
			}
			if err := json.Unmarshal(body, &requestData); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if requestData.Name != "test-file" || requestData.Extension != "txt" {
				t.Errorf("Ожидались имя 'test-file' и расширение 'txt', получены '%s' и '%s'", requestData.Name, requestData.Extension)
			}

			// Отправляем ответ с URL для загрузки
			response := struct {
				URL         string `json:"url"`
				Description string `json:"description"`
			}{
				URL:         "http://example.com/upload",
				Description: "Upload URL",
			}
			responseJSON, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
		} else if r.Method == "GET" && r.URL.Path == "/api/file/download" {
			// Проверяем параметры запроса
			if r.URL.Query().Get("label") != "test-file" {
				t.Errorf("Ожидалась метка 'test-file', получена '%s'", r.URL.Query().Get("label"))
			}

			// Отправляем ответ с URL для скачивания
			response := struct {
				URL         string              `json:"url"`
				Description string              `json:"description"`
				Metadata    domain.FileMetadata `json:"metadata"`
				MetaInfo    string              `json:"meta_info"`
			}{
				URL:         "http://example.com/download",
				Description: "Download URL",
				Metadata: domain.FileMetadata{
					FileName:  "test-file",
					Extension: "txt",
				},
				MetaInfo: "test metadata",
			}
			responseJSON, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseJSON)
		} else {
			t.Errorf("Неожиданный запрос: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Создаем клиентский сервис с адресом тестового сервера
	serverAddr := server.URL[7:] // Убираем "http://" из URL
	clientService := NewClientService(serverAddr)

	// Тестируем GetUploadLink
	url, err := clientService.GetUploadLink("test-file", "txt", "test metadata", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове GetUploadLink: %v", err)
	}
	if url != "http://example.com/upload" {
		t.Errorf("Ожидался URL 'http://example.com/upload', получен '%s'", url)
	}

	// Тестируем GetDownloadLink
	downloadURL, fileMetadata, metaInfo, err := clientService.GetDownloadLink("test-file", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове GetDownloadLink: %v", err)
	}
	if downloadURL != "http://example.com/download" {
		t.Errorf("Ожидался URL 'http://example.com/download', получен '%s'", downloadURL)
	}
	if fileMetadata.FileName != "test-file" || fileMetadata.Extension != "txt" {
		t.Errorf("Ожидались имя файла 'test-file' и расширение 'txt', получены '%s' и '%s'", fileMetadata.FileName, fileMetadata.Extension)
	}
	if metaInfo != "test metadata" {
		t.Errorf("Ожидались метаданные 'test metadata', получены '%s'", metaInfo)
	}
}
