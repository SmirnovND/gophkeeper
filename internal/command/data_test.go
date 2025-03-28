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

// MockDataClientUseCase - расширенный мок для интерфейса ClientUseCase с методами для работы с данными
type MockDataClientUseCase struct {
	SaveTextFunc         func(label string, textData *domain.TextData, metadata string) error
	GetTextFunc          func(label string) (*domain.TextData, string, error)
	DeleteTextFunc       func(label string) error
	SaveCardFunc         func(label string, cardData *domain.CardData, metadata string) error
	GetCardFunc          func(label string) (*domain.CardData, string, error)
	DeleteCardFunc       func(label string) error
	SaveCredentialFunc   func(label string, credentialData *domain.CredentialData, metadata string) error
	GetCredentialFunc    func(label string) (*domain.CredentialData, string, error)
	DeleteCredentialFunc func(label string) error
}

// Реализация методов интерфейса ClientUseCase для работы с текстовыми данными
func (m *MockDataClientUseCase) SaveText(label string, textData *domain.TextData, metadata string) error {
	if m.SaveTextFunc != nil {
		return m.SaveTextFunc(label, textData, metadata)
	}
	return nil
}

func (m *MockDataClientUseCase) GetText(label string) (*domain.TextData, string, error) {
	if m.GetTextFunc != nil {
		return m.GetTextFunc(label)
	}
	return nil, "", nil
}

func (m *MockDataClientUseCase) DeleteText(label string) error {
	if m.DeleteTextFunc != nil {
		return m.DeleteTextFunc(label)
	}
	return nil
}

// Реализация методов интерфейса ClientUseCase для работы с данными карт
func (m *MockDataClientUseCase) SaveCard(label string, cardData *domain.CardData, metadata string) error {
	if m.SaveCardFunc != nil {
		return m.SaveCardFunc(label, cardData, metadata)
	}
	return nil
}

func (m *MockDataClientUseCase) GetCard(label string) (*domain.CardData, string, error) {
	if m.GetCardFunc != nil {
		return m.GetCardFunc(label)
	}
	return nil, "", nil
}

func (m *MockDataClientUseCase) DeleteCard(label string) error {
	if m.DeleteCardFunc != nil {
		return m.DeleteCardFunc(label)
	}
	return nil
}

// Реализация методов интерфейса ClientUseCase для работы с учетными данными
func (m *MockDataClientUseCase) SaveCredential(label string, credentialData *domain.CredentialData, metadata string) error {
	if m.SaveCredentialFunc != nil {
		return m.SaveCredentialFunc(label, credentialData, metadata)
	}
	return nil
}

func (m *MockDataClientUseCase) GetCredential(label string) (*domain.CredentialData, string, error) {
	if m.GetCredentialFunc != nil {
		return m.GetCredentialFunc(label)
	}
	return nil, "", nil
}

func (m *MockDataClientUseCase) DeleteCredential(label string) error {
	if m.DeleteCredentialFunc != nil {
		return m.DeleteCredentialFunc(label)
	}
	return nil
}

// Реализация остальных методов интерфейса ClientUseCase, которые не используются в тестах
func (m *MockDataClientUseCase) Login(username string, password string) error {
	return nil
}

func (m *MockDataClientUseCase) Register(username string, password string, passwordCheck string) error {
	return nil
}

func (m *MockDataClientUseCase) Upload(filePath string, label string) (string, error) {
	return "", nil
}

func (m *MockDataClientUseCase) Download(label string) error {
	return nil
}

// Тесты для команд работы с текстовыми данными

// TestCommand_SaveTextCmd_Success тестирует успешное сохранение текстовых данных
func TestCommand_SaveTextCmd_Success(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "test_label\ntest_content\n"
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
	mockClientUseCase := &MockDataClientUseCase{
		SaveTextFunc: func(label string, textData *domain.TextData, metadata string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if textData.Content != "test_content" {
				t.Errorf("Ожидалось содержимое 'test_content', получено '%s'", textData.Content)
			}
			return nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду save-text
	saveTextCmd := cmd.SaveTextCmd()

	// Выполняем команду
	saveTextCmd.Run(saveTextCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Текст успешно сохранен!") {
		t.Errorf("Ожидалось сообщение об успешном сохранении текста, получено: %s", output)
	}
}

// TestCommand_SaveTextCmd_Error тестирует ошибку при сохранении текстовых данных
func TestCommand_SaveTextCmd_Error(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "test_label\ntest_content\n"
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
	mockClientUseCase := &MockDataClientUseCase{
		SaveTextFunc: func(label string, textData *domain.TextData, metadata string) error {
			return errors.New("ошибка при сохранении текста")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду save-text
	saveTextCmd := cmd.SaveTextCmd()

	// Выполняем команду
	saveTextCmd.Run(saveTextCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при сохранении текста:") {
		t.Errorf("Ожидалось сообщение об ошибке при сохранении текста, получено: %s", output)
	}
}

// TestCommand_GetTextCmd_Success тестирует успешное получение текстовых данных
func TestCommand_GetTextCmd_Success(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		GetTextFunc: func(label string) (*domain.TextData, string, error) {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			return &domain.TextData{
				Content: "test_content",
			}, "", nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду get-text
	getTextCmd := cmd.GetTextCmd()

	// Выполняем команду
	getTextCmd.Run(getTextCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "test_content") {
		t.Errorf("Ожидалось содержимое 'test_content' в выводе, получено: %s", output)
	}
}

// TestCommand_GetTextCmd_Error тестирует ошибку при получении текстовых данных
func TestCommand_GetTextCmd_Error(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		GetTextFunc: func(label string) (*domain.TextData, string, error) {
			return nil, "", errors.New("ошибка при получении текста")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду get-text
	getTextCmd := cmd.GetTextCmd()

	// Выполняем команду
	getTextCmd.Run(getTextCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при получении текста:") {
		t.Errorf("Ожидалось сообщение об ошибке при получении текста, получено: %s", output)
	}
}

// TestCommand_DeleteTextCmd_Success тестирует успешное удаление текстовых данных
func TestCommand_DeleteTextCmd_Success(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		DeleteTextFunc: func(label string) error {
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

	// Получаем команду delete-text
	deleteTextCmd := cmd.DeleteTextCmd()

	// Выполняем команду
	deleteTextCmd.Run(deleteTextCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Текст успешно удален!") {
		t.Errorf("Ожидалось сообщение об успешном удалении текста, получено: %s", output)
	}
}

// TestCommand_DeleteTextCmd_Error тестирует ошибку при удалении текстовых данных
func TestCommand_DeleteTextCmd_Error(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		DeleteTextFunc: func(label string) error {
			return errors.New("ошибка при удалении текста")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду delete-text
	deleteTextCmd := cmd.DeleteTextCmd()

	// Выполняем команду
	deleteTextCmd.Run(deleteTextCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при удалении текста:") {
		t.Errorf("Ожидалось сообщение об ошибке при удалении текста, получено: %s", output)
	}
}

// Тесты для команд работы с данными кредитных карт

// TestCommand_SaveCardCmd_Success тестирует успешное сохранение данных кредитной карты
func TestCommand_SaveCardCmd_Success(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	// Учитываем, что fmt.Fscanln считывает только до первого пробела
	input := "test_label\n1234567890123456\nJohn\n12/25\n123\n"
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
	mockClientUseCase := &MockDataClientUseCase{
		SaveCardFunc: func(label string, cardData *domain.CardData, metadata string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if cardData.Number != "1234567890123456" {
				t.Errorf("Ожидался номер карты '1234567890123456', получен '%s'", cardData.Number)
			}
			if cardData.Holder != "John" {
				t.Errorf("Ожидался держатель карты 'John', получен '%s'", cardData.Holder)
			}
			if cardData.ExpiryDate != "12/25" {
				t.Errorf("Ожидался срок действия '12/25', получен '%s'", cardData.ExpiryDate)
			}
			if cardData.CVV != "123" {
				t.Errorf("Ожидался CVV '123', получен '%s'", cardData.CVV)
			}
			return nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду save-card
	saveCardCmd := cmd.SaveCardCmd()

	// Выполняем команду
	saveCardCmd.Run(saveCardCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Данные карты успешно сохранены!") {
		t.Errorf("Ожидалось сообщение об успешном сохранении данных карты, получено: %s", output)
	}
}

// TestCommand_SaveCardCmd_Error тестирует ошибку при сохранении данных кредитной карты
func TestCommand_SaveCardCmd_Error(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	// Учитываем, что fmt.Fscanln считывает только до первого пробела
	input := "test_label\n1234567890123456\nJohn\n12/25\n123\n"
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
	mockClientUseCase := &MockDataClientUseCase{
		SaveCardFunc: func(label string, cardData *domain.CardData, metadata string) error {
			return errors.New("ошибка при сохранении данных карты")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду save-card
	saveCardCmd := cmd.SaveCardCmd()

	// Выполняем команду
	saveCardCmd.Run(saveCardCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при сохранении данных карты:") {
		t.Errorf("Ожидалось сообщение об ошибке при сохранении данных карты, получено: %s", output)
	}
}

// TestCommand_GetCardCmd_Success тестирует успешное получение данных кредитной карты
func TestCommand_GetCardCmd_Success(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		GetCardFunc: func(label string) (*domain.CardData, string, error) {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			return &domain.CardData{
				Number:     "1234567890123456",
				Holder:     "John",
				ExpiryDate: "12/25",
				CVV:        "123",
			}, "", nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду get-card
	getCardCmd := cmd.GetCardCmd()

	// Выполняем команду
	getCardCmd.Run(getCardCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "1234567890123456") ||
		!strings.Contains(output, "John") ||
		!strings.Contains(output, "12/25") ||
		!strings.Contains(output, "123") {
		t.Errorf("Ожидались данные карты в выводе, получено: %s", output)
	}
}

// TestCommand_GetCardCmd_Error тестирует ошибку при получении данных кредитной карты
func TestCommand_GetCardCmd_Error(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		GetCardFunc: func(label string) (*domain.CardData, string, error) {
			return nil, "", errors.New("ошибка при получении данных карты")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду get-card
	getCardCmd := cmd.GetCardCmd()

	// Выполняем команду
	getCardCmd.Run(getCardCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при получении данных карты:") {
		t.Errorf("Ожидалось сообщение об ошибке при получении данных карты, получено: %s", output)
	}
}

// TestCommand_DeleteCardCmd_Success тестирует успешное удаление данных кредитной карты
func TestCommand_DeleteCardCmd_Success(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		DeleteCardFunc: func(label string) error {
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

	// Получаем команду delete-card
	deleteCardCmd := cmd.DeleteCardCmd()

	// Выполняем команду
	deleteCardCmd.Run(deleteCardCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Данные карты успешно удалены!") {
		t.Errorf("Ожидалось сообщение об успешном удалении данных карты, получено: %s", output)
	}
}

// TestCommand_DeleteCardCmd_Error тестирует ошибку при удалении данных кредитной карты
func TestCommand_DeleteCardCmd_Error(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		DeleteCardFunc: func(label string) error {
			return errors.New("ошибка при удалении данных карты")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду delete-card
	deleteCardCmd := cmd.DeleteCardCmd()

	// Выполняем команду
	deleteCardCmd.Run(deleteCardCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при удалении данных карты:") {
		t.Errorf("Ожидалось сообщение об ошибке при удалении данных карты, получено: %s", output)
	}
}

// Тесты для команд работы с учетными данными

// TestCommand_SaveCredentialCmd_Success тестирует успешное сохранение учетных данных
func TestCommand_SaveCredentialCmd_Success(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "test_label\ntest_login\ntest_password\n"
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
	mockClientUseCase := &MockDataClientUseCase{
		SaveCredentialFunc: func(label string, credentialData *domain.CredentialData, metadata string) error {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			if credentialData.Login != "test_login" {
				t.Errorf("Ожидался логин 'test_login', получен '%s'", credentialData.Login)
			}
			if credentialData.Password != "test_password" {
				t.Errorf("Ожидался пароль 'test_password', получен '%s'", credentialData.Password)
			}
			return nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду save-credential
	saveCredentialCmd := cmd.SaveCredentialCmd()

	// Выполняем команду
	saveCredentialCmd.Run(saveCredentialCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Учетные данные успешно сохранены!") {
		t.Errorf("Ожидалось сообщение об успешном сохранении учетных данных, получено: %s", output)
	}
}

// TestCommand_SaveCredentialCmd_Error тестирует ошибку при сохранении учетных данных
func TestCommand_SaveCredentialCmd_Error(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем пайп для эмуляции ввода пользователя
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Пишем тестовые данные в пайп
	input := "test_label\ntest_login\ntest_password\n"
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
	mockClientUseCase := &MockDataClientUseCase{
		SaveCredentialFunc: func(label string, credentialData *domain.CredentialData, metadata string) error {
			return errors.New("ошибка при сохранении учетных данных")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду save-credential
	saveCredentialCmd := cmd.SaveCredentialCmd()

	// Выполняем команду
	saveCredentialCmd.Run(saveCredentialCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при сохранении учетных данных:") {
		t.Errorf("Ожидалось сообщение об ошибке при сохранении учетных данных, получено: %s", output)
	}
}

// TestCommand_GetCredentialCmd_Success тестирует успешное получение учетных данных
func TestCommand_GetCredentialCmd_Success(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		GetCredentialFunc: func(label string) (*domain.CredentialData, string, error) {
			// Проверяем параметры
			if label != "test_label" {
				t.Errorf("Ожидалась метка 'test_label', получена '%s'", label)
			}
			return &domain.CredentialData{
				Login:    "test_login",
				Password: "test_password",
			}, "", nil
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду get-credential
	getCredentialCmd := cmd.GetCredentialCmd()

	// Выполняем команду
	getCredentialCmd.Run(getCredentialCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "test_login") || !strings.Contains(output, "test_password") {
		t.Errorf("Ожидались учетные данные в выводе, получено: %s", output)
	}
}

// TestCommand_GetCredentialCmd_Error тестирует ошибку при получении учетных данных
func TestCommand_GetCredentialCmd_Error(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		GetCredentialFunc: func(label string) (*domain.CredentialData, string, error) {
			return nil, "", errors.New("ошибка при получении учетных данных")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду get-credential
	getCredentialCmd := cmd.GetCredentialCmd()

	// Выполняем команду
	getCredentialCmd.Run(getCredentialCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при получении учетных данных:") {
		t.Errorf("Ожидалось сообщение об ошибке при получении учетных данных, получено: %s", output)
	}
}

// TestCommand_DeleteCredentialCmd_Success тестирует успешное удаление учетных данных
func TestCommand_DeleteCredentialCmd_Success(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		DeleteCredentialFunc: func(label string) error {
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

	// Получаем команду delete-credential
	deleteCredentialCmd := cmd.DeleteCredentialCmd()

	// Выполняем команду
	deleteCredentialCmd.Run(deleteCredentialCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Учетные данные успешно удалены!") {
		t.Errorf("Ожидалось сообщение об успешном удалении учетных данных, получено: %s", output)
	}
}

// TestCommand_DeleteCredentialCmd_Error тестирует ошибку при удалении учетных данных
func TestCommand_DeleteCredentialCmd_Error(t *testing.T) {
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
	mockClientUseCase := &MockDataClientUseCase{
		DeleteCredentialFunc: func(label string) error {
			return errors.New("ошибка при удалении учетных данных")
		},
	}

	// Создаем экземпляр Command
	cmd := &Command{
		clientUseCase: mockClientUseCase,
	}

	// Получаем команду delete-credential
	deleteCredentialCmd := cmd.DeleteCredentialCmd()

	// Выполняем команду
	deleteCredentialCmd.Run(deleteCredentialCmd, []string{})

	// Закрываем pipe для записи и копируем вывод в буфер
	w2.Close()
	io.Copy(&buf, r2)

	// Проверяем вывод
	output := buf.String()
	if !strings.Contains(output, "Ошибка при удалении учетных данных:") {
		t.Errorf("Ожидалось сообщение об ошибке при удалении учетных данных, получено: %s", output)
	}
}
