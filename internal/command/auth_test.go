package command

import (
	"bytes"
	"errors"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"io"
	"os"
	"strings"
	"testing"
)

// MockClientUseCase - мок для интерфейса ClientUseCase
type MockClientUseCase struct {
	LoginFunc    func(username string, password string) error
	RegisterFunc func(username string, password string, passwordCheck string) error
}

func (m *MockClientUseCase) Login(username string, password string) error {
	if m.LoginFunc != nil {
		return m.LoginFunc(username, password)
	}
	return nil
}

func (m *MockClientUseCase) Register(username string, password string, passwordCheck string) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(username, password, passwordCheck)
	}
	return nil
}

// Реализация остальных методов интерфейса ClientUseCase
func (m *MockClientUseCase) Upload(filePath string, label string) (string, error) {
	return "", nil
}

func (m *MockClientUseCase) Download(label string) error {
	return nil
}

func (m *MockClientUseCase) SaveText(label string, textData *domain.TextData) error {
	return nil
}

func (m *MockClientUseCase) GetText(label string) (*domain.TextData, error) {
	return nil, nil
}

func (m *MockClientUseCase) DeleteText(label string) error {
	return nil
}

func (m *MockClientUseCase) SaveCard(label string, cardData *domain.CardData) error {
	return nil
}

func (m *MockClientUseCase) GetCard(label string) (*domain.CardData, error) {
	return nil, nil
}

func (m *MockClientUseCase) DeleteCard(label string) error {
	return nil
}

func (m *MockClientUseCase) SaveCredential(label string, credentialData *domain.CredentialData) error {
	return nil
}

func (m *MockClientUseCase) GetCredential(label string) (*domain.CredentialData, error) {
	return nil, nil
}

func (m *MockClientUseCase) DeleteCredential(label string) error {
	return nil
}

// TestCommand_Login_Success тестирует успешную авторизацию
func TestCommand_Login_Success(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "testuser\ntestpassword\n"
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	// Создаем буфер для перехвата вывода
	var buf bytes.Buffer
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	// Создаем мок для ClientUseCase
	mockClientUseCase := &MockClientUseCase{
		LoginFunc: func(username string, password string) error {
			// Проверяем параметры
			if username != "testuser" {
				t.Errorf("Ожидался логин 'testuser', получен '%s'", username)
			}
			if password != "testpassword" {
				t.Errorf("Ожидался пароль 'testpassword', получен '%s'", password)
			}
			return nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду login
	loginCmd := cmd.Login()

	// Выполняем команду
	loginCmd.Run(loginCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Успешная авторизация!") {
		t.Errorf("Ожидалось сообщение об успешной авторизации, получено: %s", output)
	}
}

// TestCommand_Login_Error тестирует ошибку при авторизации
func TestCommand_Login_Error(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "testuser\nwrongpassword\n"
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	// Создаем буфер для перехвата вывода
	var buf bytes.Buffer
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	// Создаем мок для ClientUseCase
	mockClientUseCase := &MockClientUseCase{
		LoginFunc: func(username string, password string) error {
			return errors.New("неверный логин или пароль")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду login
	loginCmd := cmd.Login()

	// Выполняем команду
	loginCmd.Run(loginCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка авторизации:") {
		t.Errorf("Ожидалось сообщение об ошибке авторизации, получено: %s", output)
	}
}

// TestCommand_RegisterCmd_Success тестирует успешную регистрацию
func TestCommand_RegisterCmd_Success(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "newuser\nnewpassword\nnewpassword\n"
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	// Создаем буфер для перехвата вывода
	var buf bytes.Buffer
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	// Создаем мок для ClientUseCase
	mockClientUseCase := &MockClientUseCase{
		RegisterFunc: func(username string, password string, passwordCheck string) error {
			// Проверяем параметры
			if username != "newuser" {
				t.Errorf("Ожидался логин 'newuser', получен '%s'", username)
			}
			if password != "newpassword" {
				t.Errorf("Ожидался пароль 'newpassword', получен '%s'", password)
			}
			if passwordCheck != "newpassword" {
				t.Errorf("Ожидалось подтверждение пароля 'newpassword', получено '%s'", passwordCheck)
			}
			return nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду register
	registerCmd := cmd.RegisterCmd()

	// Выполняем команду
	registerCmd.Run(registerCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Успешная регистрация!") {
		t.Errorf("Ожидалось сообщение об успешной регистрации, получено: %s", output)
	}
}

// TestCommand_RegisterCmd_Error тестирует ошибку при регистрации
func TestCommand_RegisterCmd_Error(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "existinguser\npassword\npassword\n"
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	// Создаем буфер для перехвата вывода
	var buf bytes.Buffer
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	// Создаем мок для ClientUseCase
	mockClientUseCase := &MockClientUseCase{
		RegisterFunc: func(username string, password string, passwordCheck string) error {
			return errors.New("пользователь уже существует")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду register
	registerCmd := cmd.RegisterCmd()

	// Выполняем команду
	registerCmd.Run(registerCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка регистрации:") {
		t.Errorf("Ожидалось сообщение об ошибке регистрации, получено: %s", output)
	}
}
