package config

import (
	"os"
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
	// Поскольку LoadConfig вызывает log.Fatal при ошибке, мы не можем напрямую проверить это поведение
	// в стандартном тесте. В реальном коде лучше изменить LoadConfig, чтобы он возвращал ошибку
	// вместо вызова log.Fatal.
	
	// Этот тест демонстрирует проблему с текущей реализацией
	t.Skip("Этот тест пропущен, так как LoadConfig вызывает log.Fatal при ошибке")
	
	config := &Config{}
	config.LoadConfig("несуществующий_файл.yaml")
	// Этот код никогда не будет выполнен из-за log.Fatal в LoadConfig
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

// Тест для NewConfig сложнее написать, так как он зависит от аргументов командной строки
// и вызывает log.Fatal при ошибке. Для полноценного тестирования нужно рефакторить код.
// Вот пример, как можно было бы тестировать NewConfig, если бы он был более тестируемым:

/*
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
	// и т.д. для других полей
}
*/