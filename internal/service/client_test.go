package service

import (
	"encoding/json"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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
	
	// Тестируем ошибки в методах работы с текстовыми данными
	// Создаем тестовый сервер для проверки ошибок
	textErrorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок Authorization
		if r.Header.Get("Authorization") != "test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		// Обрабатываем разные методы и пути
		if r.Method == "POST" && r.URL.Path == "/api/data/text/error-text" {
			// Возвращаем ошибку сервера при сохранении
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/text/not-found-text" {
			// Возвращаем ошибку "не найдено" при получении
			w.WriteHeader(http.StatusNotFound)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/text/error-text" {
			// Возвращаем ошибку сервера при получении
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/text/invalid-json" {
			// Возвращаем некорректный JSON при получении
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{invalid json}"))
		} else if r.Method == "DELETE" && r.URL.Path == "/api/data/text/not-found-text" {
			// Возвращаем ошибку "не найдено" при удалении
			w.WriteHeader(http.StatusNotFound)
		} else if r.Method == "DELETE" && r.URL.Path == "/api/data/text/error-text" {
			// Возвращаем ошибку сервера при удалении
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer textErrorServer.Close()
	
	// Создаем клиентский сервис с адресом тестового сервера для ошибок
	textErrorServerAddr := textErrorServer.URL[7:] // Убираем "http://" из URL
	textErrorClientService := NewClientService(textErrorServerAddr)
	
	// Тестируем ошибки в SaveText
	// Тест на ошибку авторизации
	err = textErrorClientService.SaveText("test-text", &domain.TextData{Content: "test"}, "", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации при сохранении текста, но ее не было")
	}
	
	// Тест на ошибку сервера
	err = textErrorClientService.SaveText("error-text", &domain.TextData{Content: "test"}, "", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера при сохранении текста, но ее не было")
	}
	
	// Тестируем ошибки в GetText
	// Тест на ошибку авторизации
	_, _, err = textErrorClientService.GetText("test-text", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации при получении текста, но ее не было")
	}
	
	// Тест на ошибку "не найдено"
	_, _, err = textErrorClientService.GetText("not-found-text", "test-token")
	if err == nil || err.Error() != "текстовые данные не найдены" {
		t.Errorf("Ожидалась ошибка 'текстовые данные не найдены', получено: %v", err)
	}
	
	// Тест на ошибку сервера
	_, _, err = textErrorClientService.GetText("error-text", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера при получении текста, но ее не было")
	}
	
	// Тест на некорректный JSON
	_, _, err = textErrorClientService.GetText("invalid-json", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка при парсинге JSON, но ее не было")
	}
	
	// Тестируем ошибки в DeleteText
	// Тест на ошибку авторизации
	err = textErrorClientService.DeleteText("test-text", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации при удалении текста, но ее не было")
	}
	
	// Тест на ошибку "не найдено"
	err = textErrorClientService.DeleteText("not-found-text", "test-token")
	if err == nil || err.Error() != "текстовые данные не найдены" {
		t.Errorf("Ожидалась ошибка 'текстовые данные не найдены', получено: %v", err)
	}
	
	// Тест на ошибку сервера
	err = textErrorClientService.DeleteText("error-text", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера при удалении текста, но ее не было")
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
	
	// Тестируем ошибки в Register
	// Создаем тестовый сервер для проверки ошибок
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/user/register" {
			// Проверяем тело запроса
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}
			
			var credentials domain.Credentials
			if err := json.Unmarshal(body, &credentials); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}
			
			// Возвращаем ошибку конфликта, если логин "existinguser"
			if credentials.Login == "existinguser" {
				w.WriteHeader(http.StatusConflict)
				return
			}
			
			// Возвращаем другую ошибку для других случаев
			if credentials.Login == "erroruser" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			
			// Возвращаем успешный ответ, но без токена
			if credentials.Login == "notokenuser" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}))
	defer errorServer.Close()
	
	// Создаем клиентский сервис с адресом тестового сервера для ошибок
	errorServerAddr := errorServer.URL[7:] // Убираем "http://" из URL
	errorClientService := NewClientService(errorServerAddr)
	
	// Тест на конфликт (пользователь уже существует)
	_, err = errorClientService.Register("existinguser", "password")
	if err == nil || err.Error() != "пользователь с таким логином уже существует" {
		t.Errorf("Ожидалась ошибка 'пользователь с таким логином уже существует', получено: %v", err)
	}
	
	// Тест на другую ошибку сервера
	_, err = errorClientService.Register("erroruser", "password")
	if err == nil {
		t.Error("Ожидалась ошибка при регистрации, но ее не было")
	}
	
	// Тест на отсутствие токена в ответе
	_, err = errorClientService.Register("notokenuser", "password")
	if err == nil || err.Error() != "токен не найден в ответе" {
		t.Errorf("Ожидалась ошибка 'токен не найден в ответе', получено: %v", err)
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
		} else if r.Method == "PUT" && r.URL.Path == "/upload" {
			// Проверяем заголовок Content-Type для загрузки файла
			if r.Header.Get("Content-Type") != "application/octet-stream" {
				t.Errorf("Ожидался Content-Type application/octet-stream, получен %s", r.Header.Get("Content-Type"))
			}
			
			// Читаем содержимое файла
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}
			
			// Проверяем содержимое файла
			if string(body) != "test file content" {
				t.Errorf("Ожидалось содержимое 'test file content', получено '%s'", string(body))
			}
			
			w.WriteHeader(http.StatusOK)
		} else if r.Method == "GET" && r.URL.Path == "/download" {
			// Отправляем тестовое содержимое файла
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("downloaded file content"))
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
	
	// Тестируем ошибки в GetUploadLink
	// Создаем тестовый сервер для проверки ошибок
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/file/upload" {
			// Проверяем заголовок Authorization
			if r.Header.Get("Authorization") != "test-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			
			// Возвращаем ошибку сервера
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer errorServer.Close()
	
	// Создаем клиентский сервис с адресом тестового сервера для ошибок
	errorServerAddr := errorServer.URL[7:] // Убираем "http://" из URL
	errorClientService := NewClientService(errorServerAddr)
	
	// Тест на ошибку сервера
	_, err = errorClientService.GetUploadLink("test-file", "txt", "test metadata", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка при получении ссылки для загрузки, но ее не было")
	}
	
	// Тест на ошибку авторизации
	_, err = errorClientService.GetUploadLink("test-file", "txt", "test metadata", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации, но ее не было")
	}
	
	// Тест на ошибку маршалинга данных (невозможно создать в реальном коде, но можно проверить через мок)
	// Этот тест будет покрыт в других тестах, так как маршалинг используется во многих методах

	// Тестируем SendFileToServer
	// Создаем временный файл для тестирования
	tempFile, err := ioutil.TempFile("", "test_file_*.txt")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Записываем тестовое содержимое в файл
	_, err = tempFile.WriteString("test file content")
	if err != nil {
		t.Fatalf("Ошибка при записи в файл: %v", err)
	}
	tempFile.Seek(0, 0) // Перемещаем указатель в начало файла
	
	// Тестируем отправку файла
	message, err := clientService.SendFileToServer(server.URL+"/upload", tempFile)
	if err != nil {
		t.Fatalf("Ошибка при вызове SendFileToServer: %v", err)
	}
	if message != "Файл успешно загружен!" {
		t.Errorf("Ожидалось сообщение 'Файл успешно загружен!', получено '%s'", message)
	}
	
	// Тестируем DownloadFileFromServer
	// Создаем временный файл для сохранения скачанного содержимого
	downloadTempFile, err := ioutil.TempFile("", "download_test_*.txt")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла для скачивания: %v", err)
	}
	downloadTempFile.Close() // Закрываем файл, так как он будет открыт в DownloadFileFromServer
	defer os.Remove(downloadTempFile.Name())
	
	// Тестируем скачивание файла
	err = clientService.DownloadFileFromServer(server.URL+"/download", downloadTempFile.Name())
	if err != nil {
		t.Fatalf("Ошибка при вызове DownloadFileFromServer: %v", err)
	}
	
	// Проверяем содержимое скачанного файла
	downloadedContent, err := ioutil.ReadFile(downloadTempFile.Name())
	if err != nil {
		t.Fatalf("Ошибка при чтении скачанного файла: %v", err)
	}
	if string(downloadedContent) != "downloaded file content" {
		t.Errorf("Ожидалось содержимое 'downloaded file content', получено '%s'", string(downloadedContent))
	}

	// Тестируем GetDownloadLink
	downloadURL, fileMetadata, metaInfo, err := clientService.GetDownloadLink("test-file", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове GetDownloadLink: %v", err)
	}
	
	// Тестируем ошибки в GetDownloadLink
	// Создаем тестовый сервер для проверки ошибок
	downloadErrorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/file/download" {
			// Проверяем заголовок Authorization
			if r.Header.Get("Authorization") != "test-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			
			// Проверяем параметр label
			if r.URL.Query().Get("label") == "not-found" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			
			// Возвращаем ошибку сервера для других случаев
			if r.URL.Query().Get("label") == "server-error" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			
			// Возвращаем некорректный JSON
			if r.URL.Query().Get("label") == "invalid-json" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{invalid json}"))
				return
			}
		}
	}))
	defer downloadErrorServer.Close()
	
	// Создаем клиентский сервис с адресом тестового сервера для ошибок
	downloadErrorServerAddr := downloadErrorServer.URL[7:] // Убираем "http://" из URL
	downloadErrorClientService := NewClientService(downloadErrorServerAddr)
	
	// Тест на ошибку авторизации
	_, _, _, err = downloadErrorClientService.GetDownloadLink("test-file", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации, но ее не было")
	}
	
	// Тест на ошибку "не найдено"
	_, _, _, err = downloadErrorClientService.GetDownloadLink("not-found", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка 'не найдено', но ее не было")
	}
	
	// Тест на ошибку сервера
	_, _, _, err = downloadErrorClientService.GetDownloadLink("server-error", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера, но ее не было")
	}
	
	// Тест на некорректный JSON
	_, _, _, err = downloadErrorClientService.GetDownloadLink("invalid-json", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка при парсинге JSON, но ее не было")
	}
	
	// Тестируем ошибки в SendFileToServer
	// Тест с некорректным URL
	_, err = clientService.SendFileToServer("http://invalid-url", tempFile)
	if err == nil {
		t.Error("Ожидалась ошибка при отправке файла на некорректный URL, но ее не было")
	}
	
	// Тест с некорректным файлом
	invalidFile, err := os.Open("/non-existent-file.txt")
	if err == nil {
		defer invalidFile.Close()
		_, err = clientService.SendFileToServer(server.URL+"/upload", invalidFile)
		if err == nil {
			t.Error("Ожидалась ошибка при отправке некорректного файла, но ее не было")
		}
	}
	
	// Тестируем ошибки в DownloadFileFromServer
	// Тест с некорректным URL
	err = clientService.DownloadFileFromServer("http://invalid-url", downloadTempFile.Name())
	if err == nil {
		t.Error("Ожидалась ошибка при скачивании файла с некорректного URL, но ее не было")
	}
	
	// Тест с некорректным путем для сохранения
	err = clientService.DownloadFileFromServer(server.URL+"/download", "/invalid/path/for/saving.txt")
	if err == nil {
		t.Error("Ожидалась ошибка при сохранении файла по некорректному пути, но ее не было")
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
