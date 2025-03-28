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

// MockFileClientUseCase - мок для интерфейса ClientUseCase с методами для работы с файлами
type MockFileClientUseCase struct {
	UploadFunc   func(filePath string, label string) (string, error)
	DownloadFunc func(label string) error
}

// Реализация методов интерфейса ClientUseCase для работы с файлами
func (m *MockFileClientUseCase) Upload(filePath string, label string) (string, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(filePath, label)
	}
	return "", nil
}

func (m *MockFileClientUseCase) Download(label string) error {
	if m.DownloadFunc != nil {
		return m.DownloadFunc(label)
	}
	return nil
}

// Реализация остальных методов интерфейса ClientUseCase, которые не используются в тестах
func (m *MockFileClientUseCase) Login(username string, password string) error {
	return nil
}

func (m *MockFileClientUseCase) Register(username string, password string, passwordCheck string) error {
	return nil
}

func (m *MockFileClientUseCase) SaveText(label string, textData *domain.TextData, metadata string) error {
	return nil
}

func (m *MockFileClientUseCase) GetText(label string) (*domain.TextData, string, error) {
	return nil, "", nil
}

func (m *MockFileClientUseCase) DeleteText(label string) error {
	return nil
}

func (m *MockFileClientUseCase) SaveCard(label string, cardData *domain.CardData, metadata string) error {
	return nil
}

func (m *MockFileClientUseCase) GetCard(label string) (*domain.CardData, string, error) {
	return nil, "", nil
}

func (m *MockFileClientUseCase) DeleteCard(label string) error {
	return nil
}

func (m *MockFileClientUseCase) SaveCredential(label string, credentialData *domain.CredentialData, metadata string) error {
	return nil
}

func (m *MockFileClientUseCase) GetCredential(label string) (*domain.CredentialData, string, error) {
	return nil, "", nil
}

func (m *MockFileClientUseCase) DeleteCredential(label string) error {
	return nil
}

// TestCommand_UploadCmd_Success тестирует успешную загрузку файла
func TestCommand_UploadCmd_Success(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "/path/to/file.txt\ntest_label\n"
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
	mockClientUseCase := &MockFileClientUseCase{
		UploadFunc: func(filePath string, label string) (string, error) {
			// Проверяем параметры
			if filePath != "/path/to/file.txt" {
				t.Errorf("Ожидался путь к файлу '/path/to/file.txt', получен '%s'", filePath)
			}
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			return "https://example.com/file.txt", nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду upload
	uploadCmd := cmd.UploadCmd()

	// Выполняем команду
	uploadCmd.Run(uploadCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Файл успешно загружен:") {
		t.Errorf("Ожидалось сообщение об успешной загрузке файла, получено: %s", output)
	}
	if !strings.Contains(output, "https://example.com/file.txt") {
		t.Errorf("Ожидался URL файла в выводе, получено: %s", output)
	}
}

// TestCommand_UploadCmd_Error тестирует ошибку при загрузке файла
func TestCommand_UploadCmd_Error(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "/path/to/file.txt\ntest_label\n"
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
	mockClientUseCase := &MockFileClientUseCase{
		UploadFunc: func(filePath string, label string) (string, error) {
			return "", errors.New("ошибка при загрузке файла")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду upload
	uploadCmd := cmd.UploadCmd()

	// Выполняем команду
	uploadCmd.Run(uploadCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка:") {
		t.Errorf("Ожидалось сообщение об ошибке при загрузке файла, получено: %s", output)
	}
}

// TestCommand_DownloadCmd_Success тестирует успешное скачивание файла
func TestCommand_DownloadCmd_Success(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "test_label\n"
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
	mockClientUseCase := &MockFileClientUseCase{
		DownloadFunc: func(label string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			return nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду download
	downloadCmd := cmd.DownloadCmd()

	// Выполняем команду
	downloadCmd.Run(downloadCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод - в данном случае, успешное скачивание не выводит сообщение
	// Мы просто проверяем, что нет сообщения об ошибке
	output := buf.String()
	if strings.Contains(output, "Ошибка при скачивании файла:") {
		t.Errorf("Не ожидалось сообщение об ошибке при скачивании файла, получено: %s", output)
	}
}

// TestCommand_DownloadCmd_Error тестирует ошибку при скачивании файла
func TestCommand_DownloadCmd_Error(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "test_label\n"
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
	mockClientUseCase := &MockFileClientUseCase{
		DownloadFunc: func(label string) error {
			return errors.New("ошибка при скачивании файла")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду download
	downloadCmd := cmd.DownloadCmd()

	// Выполняем команду
	downloadCmd.Run(downloadCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при скачивании файла:") {
		t.Errorf("Ожидалось сообщение об ошибке при скачивании файла, получено: %s", output)
	}
}