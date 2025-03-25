package service

import (
	"encoding/json"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Тестирование метода Login
func TestClientService_Login(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		if r.URL.Path != "/api/user/login" {
			t.Errorf("Ожидался путь /api/user/login, получен %s", r.URL.Path)
		}

		// Проверяем заголовок Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Ожидался Content-Type application/json, получен %s", r.Header.Get("Content-Type"))
		}

		// Читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Ошибка при чтении тела запроса: %v", err)
		}

		// Декодируем JSON
		var credentials domain.Credentials
		if err := json.Unmarshal(body, &credentials); err != nil {
			t.Fatalf("Ошибка при декодировании JSON: %v", err)
		}

		// Проверяем логин и пароль
		if credentials.Login != "testuser" || credentials.Password != "testpassword" {
			t.Errorf("Ожидался логин testuser и пароль testpassword, получены %s и %s", credentials.Login, credentials.Password)
		}

		// Устанавливаем заголовок Authorization
		w.Header().Set("Authorization", "test-token")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Создаем клиентский сервис с адресом тестового сервера
	serverAddr := server.URL[7:] // Убираем "http://" из URL
	clientService := NewClientService(serverAddr)

	// Вызываем метод Login
	token, err := clientService.Login("testuser", "testpassword")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове Login: %v", err)
	}
	if token != "test-token" {
		t.Errorf("Ожидался токен test-token, получен %s", token)
	}
}

// Тестирование метода Login с ошибкой аутентификации
func TestClientService_Login_AuthError(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Возвращаем ошибку аутентификации
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	// Создаем клиентский сервис с адресом тестового сервера
	serverAddr := server.URL[7:] // Убираем "http://" из URL
	clientService := NewClientService(serverAddr)

	// Вызываем метод Login
	_, err := clientService.Login("testuser", "wrongpassword")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// Тестирование метода Register
func TestClientService_Register(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		if r.URL.Path != "/api/user/register" {
			t.Errorf("Ожидался путь /api/user/register, получен %s", r.URL.Path)
		}

		// Проверяем заголовок Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Ожидался Content-Type application/json, получен %s", r.Header.Get("Content-Type"))
		}

		// Читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Ошибка при чтении тела запроса: %v", err)
		}

		// Декодируем JSON
		var credentials domain.Credentials
		if err := json.Unmarshal(body, &credentials); err != nil {
			t.Fatalf("Ошибка при декодировании JSON: %v", err)
		}

		// Проверяем логин и пароль
		if credentials.Login != "newuser" || credentials.Password != "newpassword" {
			t.Errorf("Ожидался логин newuser и пароль newpassword, получены %s и %s", credentials.Login, credentials.Password)
		}

		// Устанавливаем заголовок Authorization
		w.Header().Set("Authorization", "new-token")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Создаем клиентский сервис с адресом тестового сервера
	serverAddr := server.URL[7:] // Убираем "http://" из URL
	clientService := NewClientService(serverAddr)

	// Вызываем метод Register
	token, err := clientService.Register("newuser", "newpassword")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове Register: %v", err)
	}
	if token != "new-token" {
		t.Errorf("Ожидался токен new-token, получен %s", token)
	}
}

// Тестирование метода Register с ошибкой (пользователь уже существует)
func TestClientService_Register_UserExists(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Возвращаем ошибку "пользователь уже существует"
		w.WriteHeader(http.StatusConflict)
	}))
	defer server.Close()

	// Создаем клиентский сервис с адресом тестового сервера
	serverAddr := server.URL[7:] // Убираем "http://" из URL
	clientService := NewClientService(serverAddr)

	// Вызываем метод Register
	_, err := clientService.Register("existinguser", "password")

	// Проверяем результаты
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}
}

// Тестирование метода GetUploadLink
func TestClientService_GetUploadLink(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		if r.URL.Path != "/api/file/upload" {
			t.Errorf("Ожидался путь /api/file/upload, получен %s", r.URL.Path)
		}

		// Проверяем заголовок Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Ожидался Content-Type application/json, получен %s", r.Header.Get("Content-Type"))
		}

		// Проверяем заголовок Authorization
		if r.Header.Get("Authorization") != "test-token" {
			t.Errorf("Ожидался Authorization test-token, получен %s", r.Header.Get("Authorization"))
		}

		// Читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Ошибка при чтении тела запроса: %v", err)
		}

		// Декодируем JSON
		var requestData struct {
			Name      string `json:"name"`
			Extension string `json:"extension"`
		}
		if err := json.Unmarshal(body, &requestData); err != nil {
			t.Fatalf("Ошибка при декодировании JSON: %v", err)
		}

		// Проверяем имя и расширение
		if requestData.Name != "test-file" || requestData.Extension != "txt" {
			t.Errorf("Ожидалось имя test-file и расширение txt, получены %s и %s", requestData.Name, requestData.Extension)
		}

		// Возвращаем ответ
		response := struct {
			URL         string `json:"url"`
			Description string `json:"description"`
		}{
			URL:         "http://example.com/upload/test-file.txt",
			Description: "Upload URL for test-file.txt",
		}
		responseJSON, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))
	defer server.Close()

	// Создаем клиентский сервис с адресом тестового сервера
	serverAddr := server.URL[7:] // Убираем "http://" из URL
	clientService := NewClientService(serverAddr)

	// Вызываем метод GetUploadLink
	url, err := clientService.GetUploadLink("test-file", "txt", "test-token")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове GetUploadLink: %v", err)
	}
	if url != "http://example.com/upload/test-file.txt" {
		t.Errorf("Ожидался URL http://example.com/upload/test-file.txt, получен %s", url)
	}
}

// Тестирование метода GetDownloadLink
func TestClientService_GetDownloadLink(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		if r.URL.Path != "/api/file/download" {
			t.Errorf("Ожидался путь /api/file/download, получен %s", r.URL.Path)
		}

		// Проверяем параметр label
		if r.URL.Query().Get("label") != "test-file" {
			t.Errorf("Ожидался параметр label=test-file, получен %s", r.URL.Query().Get("label"))
		}

		// Проверяем заголовок Authorization
		if r.Header.Get("Authorization") != "test-token" {
			t.Errorf("Ожидался Authorization test-token, получен %s", r.Header.Get("Authorization"))
		}

		// Возвращаем ответ
		response := struct {
			URL         string              `json:"url"`
			Description string              `json:"description"`
			Metadata    domain.FileMetadata `json:"metadata"`
		}{
			URL:         "http://example.com/download/test-file.txt",
			Description: "Download URL for test-file.txt",
			Metadata: domain.FileMetadata{
				FileName:  "test-file",
				Extension: "txt",
			},
		}
		responseJSON, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))
	defer server.Close()

	// Создаем клиентский сервис с адресом тестового сервера
	serverAddr := server.URL[7:] // Убираем "http://" из URL
	clientService := NewClientService(serverAddr)

	// Вызываем метод GetDownloadLink
	url, metadata, err := clientService.GetDownloadLink("test-file", "test-token")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове GetDownloadLink: %v", err)
	}
	if url != "http://example.com/download/test-file.txt" {
		t.Errorf("Ожидался URL http://example.com/download/test-file.txt, получен %s", url)
	}
	if metadata.FileName != "test-file" || metadata.Extension != "txt" {
		t.Errorf("Ожидались метаданные {FileName: test-file, Extension: txt}, получены {FileName: %s, Extension: %s}", metadata.FileName, metadata.Extension)
	}
}

// Тестирование метода SendFileToServer
func TestClientService_SendFileToServer(t *testing.T) {
	// Создаем временный файл для тестирования
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "test-file.txt")
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	_, err = tempFile.WriteString("test content")
	if err != nil {
		t.Fatalf("Ошибка при записи во временный файл: %v", err)
	}
	tempFile.Close()

	// Открываем файл для чтения
	file, err := os.Open(tempFilePath)
	if err != nil {
		t.Fatalf("Ошибка при открытии временного файла: %v", err)
	}
	defer file.Close()

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PUT" {
			t.Errorf("Ожидался метод PUT, получен %s", r.Method)
		}

		// Проверяем заголовок Content-Type
		if r.Header.Get("Content-Type") != "application/octet-stream" {
			t.Errorf("Ожидался Content-Type application/octet-stream, получен %s", r.Header.Get("Content-Type"))
		}

		// Читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Ошибка при чтении тела запроса: %v", err)
		}

		// Проверяем содержимое файла
		if string(body) != "test content" {
			t.Errorf("Ожидалось содержимое 'test content', получено '%s'", string(body))
		}

		// Возвращаем успешный ответ
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Создаем клиентский сервис
	clientService := NewClientService("")

	// Вызываем метод SendFileToServer
	message, err := clientService.SendFileToServer(server.URL, file)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове SendFileToServer: %v", err)
	}
	if message != "Файл успешно загружен!" {
		t.Errorf("Ожидалось сообщение 'Файл успешно загружен!', получено '%s'", message)
	}
}

// Тестирование метода DownloadFileFromServer
func TestClientService_DownloadFileFromServer(t *testing.T) {
	// Создаем временную директорию для сохранения скачанного файла
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "downloaded-file.txt")

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Возвращаем содержимое файла
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("downloaded content"))
	}))
	defer server.Close()

	// Создаем клиентский сервис
	clientService := NewClientService("")

	// Вызываем метод DownloadFileFromServer
	err := clientService.DownloadFileFromServer(server.URL, outputPath)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове DownloadFileFromServer: %v", err)
	}

	// Проверяем, что файл был создан
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Файл не был создан по пути %s", outputPath)
	}

	// Проверяем содержимое скачанного файла
	content, err := ioutil.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Ошибка при чтении скачанного файла: %v", err)
	}
	if string(content) != "downloaded content" {
		t.Errorf("Ожидалось содержимое 'downloaded content', получено '%s'", string(content))
	}
}

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

			var textData domain.TextData
			if err := json.Unmarshal(body, &textData); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if textData.Content != "test text content" {
				t.Errorf("Ожидался текст 'test text content', получен '%s'", textData.Content)
			}

			w.WriteHeader(http.StatusOK)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/text/test-text" {
			// Получение текстовых данных
			textData := domain.TextData{
				Content: "test text content",
			}
			responseJSON, _ := json.Marshal(textData)
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
	err := clientService.SaveText("test-text", textData, "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveText: %v", err)
	}

	// Тестируем GetText
	retrievedTextData, err := clientService.GetText("test-text", "test-token")
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

			var cardData domain.CardData
			if err := json.Unmarshal(body, &cardData); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if cardData.Number != "1234567890123456" || cardData.Holder != "Test User" {
				t.Errorf("Ожидались данные карты {Number: 1234567890123456, Holder: Test User}, получены {Number: %s, Holder: %s}", cardData.Number, cardData.Holder)
			}

			w.WriteHeader(http.StatusOK)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/card/test-card" {
			// Получение данных карты
			cardData := domain.CardData{
				Number:     "1234567890123456",
				Holder:     "Test User",
				CVV:        "123",
				ExpiryDate: "12/25",
			}
			responseJSON, _ := json.Marshal(cardData)
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
		CVV:        "123",
		ExpiryDate: "12/25",
	}
	err := clientService.SaveCard("test-card", cardData, "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveCard: %v", err)
	}

	// Тестируем GetCard
	retrievedCardData, err := clientService.GetCard("test-card", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове GetCard: %v", err)
	}
	if retrievedCardData.Number != "1234567890123456" || retrievedCardData.Holder != "Test User" {
		t.Errorf("Ожидались данные карты {Number: 1234567890123456, Holder: Test User}, получены {Number: %s, Holder: %s}", retrievedCardData.Number, retrievedCardData.Holder)
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

			var credentialData domain.CredentialData
			if err := json.Unmarshal(body, &credentialData); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if credentialData.Login != "testuser" || credentialData.Password != "testpassword" {
				t.Errorf("Ожидались учетные данные {Login: testuser, Password: testpassword}, получены {Login: %s, Password: %s}", credentialData.Login, credentialData.Password)
			}

			w.WriteHeader(http.StatusOK)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/credential/test-credential" {
			// Получение учетных данных
			credentialData := domain.CredentialData{
				Login:    "testuser",
				Password: "testpassword",
			}
			responseJSON, _ := json.Marshal(credentialData)
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
		Password: "testpassword",
	}
	err := clientService.SaveCredential("test-credential", credentialData, "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове SaveCredential: %v", err)
	}

	// Тестируем GetCredential
	retrievedCredentialData, err := clientService.GetCredential("test-credential", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове GetCredential: %v", err)
	}
	if retrievedCredentialData.Login != "testuser" || retrievedCredentialData.Password != "testpassword" {
		t.Errorf("Ожидались учетные данные {Login: testuser, Password: testpassword}, получены {Login: %s, Password: %s}", retrievedCredentialData.Login, retrievedCredentialData.Password)
	}

	// Тестируем DeleteCredential
	err = clientService.DeleteCredential("test-credential", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове DeleteCredential: %v", err)
	}
}
