package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExtensionByPath(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "Файл с расширением .txt",
			filePath: "file.txt",
			expected: "txt",
		},
		{
			name:     "Файл с расширением .jpg",
			filePath: "image.jpg",
			expected: "jpg",
		},
		{
			name:     "Файл с путем и расширением",
			filePath: "/path/to/document.pdf",
			expected: "pdf",
		},
		{
			name:     "Файл с путем Windows и расширением",
			filePath: "C:\\Users\\user\\file.docx",
			expected: "docx",
		},
		{
			name:     "Файл без расширения",
			filePath: "file_without_extension",
			expected: "",
		},
		{
			name:     "Файл с точкой в имени, но без расширения",
			filePath: "file.",
			expected: "",
		},
		{
			name:     "Файл с несколькими точками",
			filePath: "archive.tar.gz",
			expected: "gz",
		},
		{
			name:     "Пустой путь",
			filePath: "",
			expected: "",
		},
		{
			name:     "Только точка",
			filePath: ".",
			expected: "",
		},
		{
			name:     "Скрытый файл в Unix",
			filePath: ".hidden",
			expected: "hidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetExtensionByPath(tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}