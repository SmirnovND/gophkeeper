package service

import (
	"context"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/minio/minio-go/v7"
	"net/url"
	"testing"
	"time"
)

// MinioClientInterface определяет интерфейс для методов minio.Client, которые мы используем
type MinioClientInterface interface {
	PresignedPutObject(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error)
	PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error)
}

// MockMinioClient - мок для MinioClientInterface
type MockMinioClient struct {
	PresignedPutObjectFunc func(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error)
	PresignedGetObjectFunc func(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error)
}

// PresignedPutObject - мок для метода PresignedPutObject
func (m *MockMinioClient) PresignedPutObject(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error) {
	return m.PresignedPutObjectFunc(ctx, bucketName, objectName, expires)
}

// PresignedGetObject - мок для метода PresignedGetObject
func (m *MockMinioClient) PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	return m.PresignedGetObjectFunc(ctx, bucketName, objectName, expires, reqParams)
}

// TestCloud_GenerateUploadLink тестирует метод GenerateUploadLink
func TestCloud_GenerateUploadLink(t *testing.T) {
	// Создаем URL для тестирования
	testURL, _ := url.Parse("https://example.com/upload/test-file.txt")

	// Создаем мок для minio.Client
	mockMinioClient := &MockMinioClient{
		PresignedPutObjectFunc: func(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error) {
			// Проверяем параметры
			if bucketName != "test-bucket" {
				t.Errorf("Ожидалось имя бакета 'test-bucket', получено '%s'", bucketName)
			}
			if objectName != "test-file.txt" {
				t.Errorf("Ожидалось имя объекта 'test-file.txt', получено '%s'", objectName)
			}
			if expires != 15*time.Minute {
				t.Errorf("Ожидалось время жизни ссылки 15 минут, получено %v", expires)
			}
			return testURL, nil
		},
	}

	// Создаем экземпляр Cloud
	cloud := &Cloud{
		minio:      mockMinioClient,
		bucketName: "test-bucket",
	}

	// Вызываем метод GenerateUploadLink
	url, err := cloud.GenerateUploadLink("test-file.txt")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове GenerateUploadLink: %v", err)
	}
	if url != "https://example.com/upload/test-file.txt" {
		t.Errorf("Ожидался URL 'https://example.com/upload/test-file.txt', получен '%s'", url)
	}
}

// TestCloud_GenerateUploadLink_Error тестирует обработку ошибок в методе GenerateUploadLink
func TestCloud_GenerateUploadLink_Error(t *testing.T) {
	// Создаем мок для minio.Client, который возвращает ошибку
	expectedError := minio.ErrorResponse{
		Code:       "AccessDenied",
		Message:    "Access Denied",
		StatusCode: 403,
	}

	mockMinioClient := &MockMinioClient{
		PresignedPutObjectFunc: func(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error) {
			return nil, expectedError
		},
	}

	// Создаем экземпляр Cloud
	cloud := &Cloud{
		minio:      mockMinioClient,
		bucketName: "test-bucket",
	}

	// Вызываем метод GenerateUploadLink
	_, err := cloud.GenerateUploadLink("test-file.txt")

	// Проверяем, что возникла ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	// Проверяем, что ошибка имеет правильный тип и содержимое
	minioErr, ok := err.(minio.ErrorResponse)
	if !ok {
		t.Fatalf("Ожидалась ошибка типа minio.ErrorResponse, получена %T", err)
	}

	if minioErr.Code != expectedError.Code ||
		minioErr.Message != expectedError.Message ||
		minioErr.StatusCode != expectedError.StatusCode {
		t.Errorf("Ожидалась ошибка %v, получена %v", expectedError, minioErr)
	}
}

// TestCloud_GenerateDownloadLink тестирует метод GenerateDownloadLink
func TestCloud_GenerateDownloadLink(t *testing.T) {
	// Создаем URL для тестирования
	testURL, _ := url.Parse("https://example.com/download/test-file.txt")

	// Создаем мок для minio.Client
	mockMinioClient := &MockMinioClient{
		PresignedGetObjectFunc: func(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
			// Проверяем параметры
			if bucketName != "test-bucket" {
				t.Errorf("Ожидалось имя бакета 'test-bucket', получено '%s'", bucketName)
			}
			if objectName != "test-file.txt" {
				t.Errorf("Ожидалось имя объекта 'test-file.txt', получено '%s'", objectName)
			}
			if expires != 15*time.Minute {
				t.Errorf("Ожидалось время жизни ссылки 15 минут, получено %v", expires)
			}
			// Проверяем, что reqParams - это пустой url.Values
			if len(reqParams) != 0 {
				t.Errorf("Ожидался пустой reqParams, получено %v", reqParams)
			}
			return testURL, nil
		},
	}

	// Создаем экземпляр Cloud
	cloud := &Cloud{
		minio:      mockMinioClient,
		bucketName: "test-bucket",
	}

	// Вызываем метод GenerateDownloadLink
	url, err := cloud.GenerateDownloadLink("test-file.txt")

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при вызове GenerateDownloadLink: %v", err)
	}
	if url != "https://example.com/download/test-file.txt" {
		t.Errorf("Ожидался URL 'https://example.com/download/test-file.txt', получен '%s'", url)
	}
}

// TestCloud_GenerateDownloadLink_Error тестирует обработку ошибок в методе GenerateDownloadLink
func TestCloud_GenerateDownloadLink_Error(t *testing.T) {
	// Создаем мок для minio.Client, который возвращает ошибку
	expectedError := minio.ErrorResponse{
		Code:       "NoSuchKey",
		Message:    "The specified key does not exist",
		StatusCode: 404,
	}

	mockMinioClient := &MockMinioClient{
		PresignedGetObjectFunc: func(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
			return nil, expectedError
		},
	}

	// Создаем экземпляр Cloud
	cloud := &Cloud{
		minio:      mockMinioClient,
		bucketName: "test-bucket",
	}

	// Вызываем метод GenerateDownloadLink
	_, err := cloud.GenerateDownloadLink("test-file.txt")

	// Проверяем, что возникла ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка, но ее не было")
	}

	// Проверяем, что ошибка имеет правильный тип и содержимое
	minioErr, ok := err.(minio.ErrorResponse)
	if !ok {
		t.Fatalf("Ожидалась ошибка типа minio.ErrorResponse, получена %T", err)
	}

	if minioErr.Code != expectedError.Code ||
		minioErr.Message != expectedError.Message ||
		minioErr.StatusCode != expectedError.StatusCode {
		t.Errorf("Ожидалась ошибка %v, получена %v", expectedError, minioErr)
	}
}

// TestNewCloud тестирует функцию NewCloud
func TestNewCloud(t *testing.T) {
	// Создаем мок для MinioClientInterface
	mockMinioClient := &MockMinioClient{
		PresignedPutObjectFunc: func(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error) {
			return nil, nil
		},
		PresignedGetObjectFunc: func(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
			return nil, nil
		},
	}

	// Вызываем функцию NewCloud
	cloud := NewCloud(mockMinioClient, "test-bucket")

	// Проверяем, что возвращенный объект не nil
	if cloud == nil {
		t.Fatal("Функция NewCloud вернула nil")
	}

	// Проверяем, что возвращенный объект реализует интерфейс CloudService
	_, ok := cloud.(interfaces.CloudService)
	if !ok {
		t.Fatal("Возвращенный объект не реализует интерфейс CloudService")
	}
}
