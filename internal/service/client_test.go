package service

import (
	"encoding/json"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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

	// Тестируем ошибки в методах работы с данными карт
	// Создаем тестовый сервер для проверки ошибок
	cardErrorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок Authorization
		if r.Header.Get("Authorization") != "test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Обрабатываем разные методы и пути
		if r.Method == "POST" && r.URL.Path == "/api/data/card/error-card" {
			// Возвращаем ошибку сервера при сохранении
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/card/not-found-card" {
			// Возвращаем ошибку "не найдено" при получении
			w.WriteHeader(http.StatusNotFound)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/card/error-card" {
			// Возвращаем ошибку сервера при получении
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/card/invalid-json" {
			// Возвращаем некорректный JSON при получении
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{invalid json}"))
		} else if r.Method == "DELETE" && r.URL.Path == "/api/data/card/not-found-card" {
			// Возвращаем ошибку "не найдено" при удалении
			w.WriteHeader(http.StatusNotFound)
		} else if r.Method == "DELETE" && r.URL.Path == "/api/data/card/error-card" {
			// Возвращаем ошибку сервера при удалении
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer cardErrorServer.Close()

	// Создаем клиентский сервис с адресом тестового сервера для ошибок
	cardErrorServerAddr := cardErrorServer.URL[7:] // Убираем "http://" из URL
	cardErrorClientService := NewClientService(cardErrorServerAddr)

	// Тестируем ошибки в SaveCard
	// Тест на ошибку авторизации
	err = cardErrorClientService.SaveCard("test-card", &domain.CardData{Number: "1234"}, "", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации при сохранении данных карты, но ее не было")
	}

	// Тест на ошибку сервера
	err = cardErrorClientService.SaveCard("error-card", &domain.CardData{Number: "1234"}, "", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера при сохранении данных карты, но ее не было")
	}

	// Тестируем ошибки в GetCard
	// Тест на ошибку авторизации
	_, _, err = cardErrorClientService.GetCard("test-card", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации при получении данных карты, но ее не было")
	}

	// Тест на ошибку "не найдено"
	_, _, err = cardErrorClientService.GetCard("not-found-card", "test-token")
	if err == nil || err.Error() != "данные карты не найдены" {
		t.Errorf("Ожидалась ошибка 'данные карты не найдены', получено: %v", err)
	}

	// Тест на ошибку сервера
	_, _, err = cardErrorClientService.GetCard("error-card", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера при получении данных карты, но ее не было")
	}

	// Тест на некорректный JSON
	_, _, err = cardErrorClientService.GetCard("invalid-json", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка при парсинге JSON, но ее не было")
	}

	// Тестируем ошибки в DeleteCard
	// Тест на ошибку авторизации
	err = cardErrorClientService.DeleteCard("test-card", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации при удалении данных карты, но ее не было")
	}

	// Тест на ошибку "не найдено"
	err = cardErrorClientService.DeleteCard("not-found-card", "test-token")
	if err == nil || err.Error() != "данные карты не найдены" {
		t.Errorf("Ожидалась ошибка 'данные карты не найдены', получено: %v", err)
	}

	// Тест на ошибку сервера
	err = cardErrorClientService.DeleteCard("error-card", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера при удалении данных карты, но ее не было")
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

	// Тестируем ошибки в методах работы с учетными данными
	// Создаем тестовый сервер для проверки ошибок
	credentialErrorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок Authorization
		if r.Header.Get("Authorization") != "test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Обрабатываем разные методы и пути
		if r.Method == "POST" && r.URL.Path == "/api/data/credential/error-credential" {
			// Возвращаем ошибку сервера при сохранении
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/credential/not-found-credential" {
			// Возвращаем ошибку "не найдено" при получении
			w.WriteHeader(http.StatusNotFound)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/credential/error-credential" {
			// Возвращаем ошибку сервера при получении
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.Method == "GET" && r.URL.Path == "/api/data/credential/invalid-json" {
			// Возвращаем некорректный JSON при получении
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{invalid json}"))
		} else if r.Method == "DELETE" && r.URL.Path == "/api/data/credential/not-found-credential" {
			// Возвращаем ошибку "не найдено" при удалении
			w.WriteHeader(http.StatusNotFound)
		} else if r.Method == "DELETE" && r.URL.Path == "/api/data/credential/error-credential" {
			// Возвращаем ошибку сервера при удалении
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer credentialErrorServer.Close()

	// Создаем клиентский сервис с адресом тестового сервера для ошибок
	credentialErrorServerAddr := credentialErrorServer.URL[7:] // Убираем "http://" из URL
	credentialErrorClientService := NewClientService(credentialErrorServerAddr)

	// Тестируем ошибки в SaveCredential
	// Тест на ошибку авторизации
	err = credentialErrorClientService.SaveCredential("test-credential", &domain.CredentialData{Login: "test"}, "", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации при сохранении учетных данных, но ее не было")
	}

	// Тест на ошибку сервера
	err = credentialErrorClientService.SaveCredential("error-credential", &domain.CredentialData{Login: "test"}, "", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера при сохранении учетных данных, но ее не было")
	}

	// Тестируем ошибки в GetCredential
	// Тест на ошибку авторизации
	_, _, err = credentialErrorClientService.GetCredential("test-credential", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации при получении учетных данных, но ее не было")
	}

	// Тест на ошибку "не найдено"
	_, _, err = credentialErrorClientService.GetCredential("not-found-credential", "test-token")
	if err == nil || err.Error() != "учетные данные не найдены" {
		t.Errorf("Ожидалась ошибка 'учетные данные не найдены', получено: %v", err)
	}

	// Тест на ошибку сервера
	_, _, err = credentialErrorClientService.GetCredential("error-credential", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера при получении учетных данных, но ее не было")
	}

	// Тест на некорректный JSON
	_, _, err = credentialErrorClientService.GetCredential("invalid-json", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка при парсинге JSON, но ее не было")
	}

	// Тестируем ошибки в DeleteCredential
	// Тест на ошибку авторизации
	err = credentialErrorClientService.DeleteCredential("test-credential", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации при удалении учетных данных, но ее не было")
	}

	// Тест на ошибку "не найдено"
	err = credentialErrorClientService.DeleteCredential("not-found-credential", "test-token")
	if err == nil || err.Error() != "учетные данные не найдены" {
		t.Errorf("Ожидалась ошибка 'учетные данные не найдены', получено: %v", err)
	}

	// Тест на ошибку сервера
	err = credentialErrorClientService.DeleteCredential("error-credential", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка сервера при удалении учетных данных, но ее не было")
	}
}

// Тестирование внутреннего метода sendRequest
func TestClientService_SendRequest(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод и путь
		if r.Method == "POST" && r.URL.Path == "/test-path" {
			// Проверяем заголовок Content-Type
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Ожидался Content-Type application/json, получен %s", r.Header.Get("Content-Type"))
			}

			// Проверяем тело запроса
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			var data struct {
				Key string `json:"key"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			if data.Key != "value" {
				t.Errorf("Ожидалось значение 'value', получено '%s'", data.Key)
			}

			// Отправляем успешный ответ
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"result": "success"}`))
		} else if r.Method == "GET" && r.URL.Path == "/test-path" {
			// Отправляем успешный ответ для GET-запроса
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"result": "success"}`))
		} else if r.URL.Path == "/error-path" {
			// Отправляем ошибку сервера
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.URL.Path == "/timeout-path" {
			// Имитируем таймаут, не отвечая на запрос
			time.Sleep(100 * time.Millisecond)
		} else {
			t.Errorf("Неожиданный запрос: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Создаем клиентский сервис с адресом тестового сервера
	clientService := &ClientService{
		client:     &http.Client{Timeout: 50 * time.Millisecond}, // Устанавливаем таймаут для тестирования
		serverAddr: server.URL[7:],                               // Убираем "http://" из URL
	}

	// Тестируем успешный POST-запрос
	data := struct {
		Key string `json:"key"`
	}{
		Key: "value",
	}
	resp, err := clientService.sendRequest("POST", server.URL+"/test-path", data)
	if err != nil {
		t.Fatalf("Ошибка при отправке POST-запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	// Тестируем успешный GET-запрос (с nil данными)
	resp, err = clientService.sendRequest("GET", server.URL+"/test-path", nil)
	if err != nil {
		t.Fatalf("Ошибка при отправке GET-запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	// Тестируем ошибку маршалинга данных
	// Создаем структуру, которая вызовет ошибку при маршалинге
	invalidData := make(chan int) // каналы нельзя маршалить в JSON
	_, err = clientService.sendRequest("POST", server.URL+"/test-path", invalidData)
	if err == nil {
		t.Error("Ожидалась ошибка при маршалинге данных, но ее не было")
	}

	// Тестируем ошибку при создании запроса
	_, err = clientService.sendRequest("POST", "://invalid-url", data)
	if err == nil {
		t.Error("Ожидалась ошибка при создании запроса, но ее не было")
	}

	// Тестируем ошибку при выполнении запроса (таймаут)
	_, err = clientService.sendRequest("GET", server.URL+"/timeout-path", nil)
	if err == nil {
		t.Error("Ожидалась ошибка таймаута, но ее не было")
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

	// Тестируем ошибки в Login и Register
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
		} else if r.Method == "POST" && r.URL.Path == "/api/user/login" {
			// Проверяем тело запроса
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			var credentials domain.Credentials
			if err := json.Unmarshal(body, &credentials); err != nil {
				t.Fatalf("Ошибка при декодировании JSON: %v", err)
			}

			// Возвращаем ошибку авторизации
			if credentials.Login == "wronguser" {
				w.WriteHeader(http.StatusUnauthorized)
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

			// Возвращаем успешный ответ с токеном для других случаев
			w.Header().Set("Authorization", "test-token")
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer errorServer.Close()

	// Создаем клиентский сервис с адресом тестового сервера для ошибок
	errorServerAddr := errorServer.URL[7:] // Убираем "http://" из URL
	errorClientService := NewClientService(errorServerAddr)

	// Тесты для Register
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

	// Тесты для Login
	// Тест на ошибку авторизации
	_, err = errorClientService.Login("wronguser", "password")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации, но ее не было")
	}

	// Тест на другую ошибку сервера
	_, err = errorClientService.Login("erroruser", "password")
	if err == nil {
		t.Error("Ожидалась ошибка при входе, но ее не было")
	}

	// Тест на отсутствие токена в ответе
	_, err = errorClientService.Login("notokenuser", "password")
	if err == nil || err.Error() != "токен не найден в ответе" {
		t.Errorf("Ожидалась ошибка 'токен не найден в ответе', получено: %v", err)
	}

	// Тест на успешный вход
	token, err = errorClientService.Login("validuser", "password")
	if err != nil {
		t.Fatalf("Ошибка при вызове Login: %v", err)
	}
	if token != "test-token" {
		t.Errorf("Ожидался токен 'test-token', получен '%s'", token)
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
	uploadErrorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/file/upload" {
			// Проверяем заголовок Authorization
			if r.Header.Get("Authorization") != "test-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Проверяем параметры запроса
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

			// Возвращаем ошибку сервера для определенного имени файла
			if requestData.Name == "error-file" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Возвращаем некорректный JSON для определенного имени файла
			if requestData.Name == "invalid-json" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{invalid json}"))
				return
			}

			// Возвращаем успешный ответ, но без URL для определенного имени файла
			if requestData.Name == "no-url" {
				response := struct {
					Description string `json:"description"`
				}{
					Description: "No URL provided",
				}
				responseJSON, _ := json.Marshal(response)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(responseJSON)
				return
			}

			// Возвращаем успешный ответ с URL для других случаев
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
		}
	}))
	defer uploadErrorServer.Close()

	// Создаем клиентский сервис с адресом тестового сервера для ошибок
	uploadErrorServerAddr := uploadErrorServer.URL[7:] // Убираем "http://" из URL
	uploadErrorClientService := NewClientService(uploadErrorServerAddr)

	// Тест на ошибку сервера
	_, err = uploadErrorClientService.GetUploadLink("error-file", "txt", "test metadata", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка при получении ссылки для загрузки, но ее не было")
	}

	// Тест на ошибку авторизации
	_, err = uploadErrorClientService.GetUploadLink("test-file", "txt", "test metadata", "invalid-token")
	if err == nil {
		t.Error("Ожидалась ошибка авторизации, но ее не было")
	}

	// Тест на некорректный JSON
	_, err = uploadErrorClientService.GetUploadLink("invalid-json", "txt", "test metadata", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка при парсинге JSON, но ее не было")
	}

	// Тест на отсутствие URL в ответе
	_, err = uploadErrorClientService.GetUploadLink("no-url", "txt", "test metadata", "test-token")
	if err == nil {
		t.Error("Ожидалась ошибка при получении URL из ответа, но ее не было")
	}

	// Тест на успешный ответ
	url, err = uploadErrorClientService.GetUploadLink("success-file", "txt", "test metadata", "test-token")
	if err != nil {
		t.Fatalf("Ошибка при вызове GetUploadLink: %v", err)
	}
	if url != "http://example.com/upload" {
		t.Errorf("Ожидался URL 'http://example.com/upload', получен '%s'", url)
	}

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

	// Тестируем ошибки в SendFileToServer
	// Создаем тестовый сервер для проверки ошибок
	fileErrorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" && r.URL.Path == "/upload-error" {
			// Возвращаем ошибку сервера
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer fileErrorServer.Close()

	// Тест на ошибку при отправке файла (ошибка сервера)
	tempFile.Seek(0, 0) // Перемещаем указатель в начало файла
	_, err = clientService.SendFileToServer(fileErrorServer.URL+"/upload-error", tempFile)
	if err == nil {
		t.Error("Ожидалась ошибка при отправке файла, но ее не было")
	}

	// Тест на ошибку при получении информации о файле
	// Создаем временный файл и закрываем его, чтобы вызвать ошибку при Stat
	invalidFile, err := ioutil.TempFile("", "invalid_file_*.txt")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	invalidFile.Close()           // Закрываем файл
	os.Remove(invalidFile.Name()) // Удаляем файл, чтобы вызвать ошибку при Stat

	_, err = clientService.SendFileToServer(server.URL+"/upload", invalidFile)
	if err == nil {
		t.Error("Ожидалась ошибка при получении информации о файле, но ее не было")
	}

	// Тест на ошибку при перемещении указателя файла
	// Создаем временный файл и закрываем его, чтобы вызвать ошибку при Seek
	closedFile, err := ioutil.TempFile("", "closed_file_*.txt")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	closedFile.Close() // Закрываем файл
	defer os.Remove(closedFile.Name())

	_, err = clientService.SendFileToServer(server.URL+"/upload", closedFile)
	if err == nil {
		t.Error("Ожидалась ошибка при перемещении указателя файла, но ее не было")
	}

	// Тест на ошибку при создании запроса
	_, err = clientService.SendFileToServer("://invalid-url", tempFile)
	if err == nil {
		t.Error("Ожидалась ошибка при создании запроса, но ее не было")
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
	invalidFileTest, err := os.Open("/non-existent-file.txt")
	if err == nil {
		defer invalidFileTest.Close()
		_, err = clientService.SendFileToServer(server.URL+"/upload", invalidFileTest)
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
