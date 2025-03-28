package config

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

// Создаем временный файл конфигурации для тестов
func createTempConfigFile(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Не удалось записать во временный файл: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Не удалось закрыть временный файл: %v", err)
	}

	return tmpfile.Name()
}

func TestConfig_LoadConfig(t *testing.T) {
	// Подготовка тестовых данных
	yamlContent := `
db:
  dsn: "postgres://user:password@localhost:5432/testdb"
app:
  jwt_secret: "test-secret"
  run_addr: ":8080"
minio:
  bucket_name: "test-bucket"
  access_key: "test-access-key"
  secret_key: "test-secret-key"
  host: "localhost:9000"
`
	configPath := createTempConfigFile(t, yamlContent)
	defer os.Remove(configPath) // Удаляем временный файл после завершения теста

	// Выполнение тестируемого кода
	config := &Config{}
	config.LoadConfig(configPath)

	// Проверка результатов
	if config.GetDBDsn() != "postgres://user:password@localhost:5432/testdb" {
		t.Errorf("Ожидалось GetDBDsn()='postgres://user:password@localhost:5432/testdb', получено '%s'", config.GetDBDsn())
	}

	if config.GetJwtSecret() != "test-secret" {
		t.Errorf("Ожидалось GetJwtSecret()='test-secret', получено '%s'", config.GetJwtSecret())
	}

	if config.GetRunAddr() != ":8080" {
		t.Errorf("Ожидалось GetRunAddr()=':8080', получено '%s'", config.GetRunAddr())
	}

	if config.GetMinioBucketName() != "test-bucket" {
		t.Errorf("Ожидалось GetMinioBucketName()='test-bucket', получено '%s'", config.GetMinioBucketName())
	}

	if config.GetMinioAccessKey() != "test-access-key" {
		t.Errorf("Ожидалось GetMinioAccessKey()='test-access-key', получено '%s'", config.GetMinioAccessKey())
	}

	if config.GetMinioSecretKey() != "test-secret-key" {
		t.Errorf("Ожидалось GetMinioSecretKey()='test-secret-key', получено '%s'", config.GetMinioSecretKey())
	}

	if config.GetMinioHost() != "localhost:9000" {
		t.Errorf("Ожидалось GetMinioHost()='localhost:9000', получено '%s'", config.GetMinioHost())
	}
}

func TestConfig_LoadConfig_InvalidFile(t *testing.T) {
	// Тест на обработку несуществующего файла
	// Поскольку LoadConfig вызывает log.Fatal при ошибке, мы перехватываем вывод лога

	// Сохраняем оригинальный вывод лога
	originalOutput := log.Writer()
	originalFlags := log.Flags()

	// Создаем буфер для перехвата вывода
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0) // Отключаем префиксы времени и т.д. для упрощения проверки

	// Восстанавливаем оригинальный вывод после теста
	defer func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	}()

	// Создаем функцию для восстановления после log.Fatal
	// log.Fatal вызывает os.Exit(1), поэтому мы должны перехватить это
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()

	var exitCode int
	osExit = func(code int) {
		exitCode = code
		// Не выходим из программы, а просто записываем код выхода
	}

	// Выполняем тестируемый код
	config := &Config{}
	config.LoadConfig("несуществующий_файл.yaml")

	// Проверяем, что был вызван log.Fatal
	if exitCode != 1 {
		t.Errorf("Ожидался вызов os.Exit(1), получен код %d", exitCode)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(buf.String(), "ReadConfigFile:") {
		t.Errorf("Ожидалось сообщение об ошибке с 'ReadConfigFile:', получено: %s", buf.String())
	}
}

func TestConfig_Getters(t *testing.T) {
	// Тестирование геттеров напрямую
	config := &Config{
		Db: Db{
			Dsn: "test-dsn",
		},
		App: App{
			JwtSecret: "test-jwt",
			RunAddr:   "test-addr",
		},
		Minio: Minio{
			BucketName: "test-bucket",
			AccessKey:  "test-access",
			SecretKey:  "test-secret",
			Host:       "test-host",
		},
	}

	// Проверка всех геттеров
	if config.GetDBDsn() != "test-dsn" {
		t.Errorf("Ожидалось GetDBDsn()='test-dsn', получено '%s'", config.GetDBDsn())
	}

	if config.GetJwtSecret() != "test-jwt" {
		t.Errorf("Ожидалось GetJwtSecret()='test-jwt', получено '%s'", config.GetJwtSecret())
	}

	if config.GetRunAddr() != "test-addr" {
		t.Errorf("Ожидалось GetRunAddr()='test-addr', получено '%s'", config.GetRunAddr())
	}

	if config.GetMinioBucketName() != "test-bucket" {
		t.Errorf("Ожидалось GetMinioBucketName()='test-bucket', получено '%s'", config.GetMinioBucketName())
	}

	if config.GetMinioAccessKey() != "test-access" {
		t.Errorf("Ожидалось GetMinioAccessKey()='test-access', получено '%s'", config.GetMinioAccessKey())
	}

	if config.GetMinioSecretKey() != "test-secret" {
		t.Errorf("Ожидалось GetMinioSecretKey()='test-secret', получено '%s'", config.GetMinioSecretKey())
	}

	if config.GetMinioHost() != "test-host" {
		t.Errorf("Ожидалось GetMinioHost()='test-host', получено '%s'", config.GetMinioHost())
	}
}

func TestConfig_LoadConfig_InvalidYAML(t *testing.T) {
	// Создаем временный файл с некорректным YAML
	invalidYAML := `
db:
  dsn: "postgres://user:password@localhost:5432/testdb"
app:
  jwt_secret: "test-secret"
  run_addr: ":8080"
minio:
  bucket_name: "test-bucket"
  access_key: "test-access-key"
  secret_key: "test-secret-key"
  host: "localhost:9000"
  invalid_yaml: [
`
	configPath := createTempConfigFile(t, invalidYAML)
	defer os.Remove(configPath)

	// Сохраняем оригинальный вывод лога
	originalOutput := log.Writer()
	originalFlags := log.Flags()

	// Создаем буфер для перехвата вывода
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0) // Отключаем префиксы времени и т.д. для упрощения проверки

	// Восстанавливаем оригинальный вывод после теста
	defer func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	}()

	// Создаем функцию для восстановления после log.Fatal
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()

	var exitCode int
	osExit = func(code int) {
		exitCode = code
		// Не выходим из программы, а просто записываем код выхода
	}

	// Выполняем тестируемый код
	config := &Config{}
	config.LoadConfig(configPath)

	// Проверяем, что был вызван log.Fatal
	if exitCode != 1 {
		t.Errorf("Ожидался вызов os.Exit(1), получен код %d", exitCode)
	}

	// Проверяем сообщение об ошибке
	if !strings.Contains(buf.String(), "DecodeConfigFile:") {
		t.Errorf("Ожидалось сообщение об ошибке с 'DecodeConfigFile:', получено: %s", buf.String())
	}
}

func TestNewConfig(t *testing.T) {
	// Сохраняем оригинальные аргументы
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Создаем временный конфиг
	yamlContent := `
db:
  dsn: "postgres://user:password@localhost:5432/testdb"
app:
  jwt_secret: "test-secret"
  run_addr: ":8080"
minio:
  bucket_name: "test-bucket"
  access_key: "test-access-key"
  secret_key: "test-secret-key"
  host: "localhost:9000"
`
	configPath := createTempConfigFile(t, yamlContent)
	defer os.Remove(configPath)

	// Подменяем аргументы командной строки
	os.Args = []string{"cmd", configPath}

	// Вызываем тестируемую функцию
	config := NewConfig()

	// Проверяем результаты
	if config.GetDBDsn() != "postgres://user:password@localhost:5432/testdb" {
		t.Errorf("Ожидалось GetDBDsn()='postgres://user:password@localhost:5432/testdb', получено '%s'", config.GetDBDsn())
	}

	if config.GetJwtSecret() != "test-secret" {
		t.Errorf("Ожидалось GetJwtSecret()='test-secret', получено '%s'", config.GetJwtSecret())
	}

	if config.GetRunAddr() != ":8080" {
		t.Errorf("Ожидалось GetRunAddr()=':8080', получено '%s'", config.GetRunAddr())
	}

	if config.GetMinioBucketName() != "test-bucket" {
		t.Errorf("Ожидалось GetMinioBucketName()='test-bucket', получено '%s'", config.GetMinioBucketName())
	}

	if config.GetMinioAccessKey() != "test-access-key" {
		t.Errorf("Ожидалось GetMinioAccessKey()='test-access-key', получено '%s'", config.GetMinioAccessKey())
	}

	if config.GetMinioSecretKey() != "test-secret-key" {
		t.Errorf("Ожидалось GetMinioSecretKey()='test-secret-key', получено '%s'", config.GetMinioSecretKey())
	}

	if config.GetMinioHost() != "localhost:9000" {
		t.Errorf("Ожидалось GetMinioHost()='localhost:9000', получено '%s'", config.GetMinioHost())
	}
}

func TestNewConfig_PanicRecovery(t *testing.T) {
	// Сохраняем оригинальные аргументы
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Устанавливаем пустые аргументы, чтобы вызвать панику
	os.Args = []string{"cmd"}

	// Перехватываем вывод лога
	originalOutput := log.Writer()
	originalFlags := log.Flags()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	}()

	// Вызываем тестируемую функцию - она должна восстановиться после паники
	// и вернуть nil, но не должна вызвать падение теста
	config := NewConfig()

	// Проверяем, что функция вернула nil после восстановления от паники
	if config != nil {
		t.Error("Ожидалось, что функция вернет nil после восстановления от паники")
	}

	// Проверяем, что в логе есть сообщение о панике
	if !strings.Contains(buf.String(), "runtime error") && !strings.Contains(buf.String(), "index out of range") {
		t.Errorf("Ожидалось сообщение о панике в логе, получено: %s", buf.String())
	}
}
