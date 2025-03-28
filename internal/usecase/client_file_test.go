package usecase

import (
	"testing"
)

// TestIsTextFile тестирует функцию isTextFile
func TestIsTextFile(t *testing.T) {
	// Тест для текстового файла
	t.Run("TextFile", func(t *testing.T) {
		result := isTextFile("test.txt")
		if !result {
			t.Error("Ожидалось true для файла с расширением .txt")
		}
	})

	// Тест для файла с другим расширением
	t.Run("NonTextFile", func(t *testing.T) {
		result := isTextFile("test.pdf")
		if result {
			t.Error("Ожидалось false для файла с расширением .pdf")
		}
	})

	// Тест для файла без расширения
	t.Run("NoExtension", func(t *testing.T) {
		result := isTextFile("test")
		if result {
			t.Error("Ожидалось false для файла без расширения")
		}
	})
}

// TestIsBinaryFile тестирует функцию isBinaryFile
func TestIsBinaryFile(t *testing.T) {
	// Тест для бинарного файла
	t.Run("BinaryFile", func(t *testing.T) {
		result := isBinaryFile("test.bin")
		if !result {
			t.Error("Ожидалось true для файла с расширением .bin")
		}
	})

	// Тест для файла с другим расширением
	t.Run("NonBinaryFile", func(t *testing.T) {
		result := isBinaryFile("test.exe")
		if result {
			t.Error("Ожидалось false для файла с расширением .exe")
		}
	})

	// Тест для файла без расширения
	t.Run("NoExtension", func(t *testing.T) {
		result := isBinaryFile("test")
		if result {
			t.Error("Ожидалось false для файла без расширения")
		}
	})
}