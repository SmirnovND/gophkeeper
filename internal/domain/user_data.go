package domain

import "time"

// UserData представляет собой структуру для хранения пользовательских данных различных типов
type UserData struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Label     string    `json:"label" db:"label"`
	Type      string    `json:"type" db:"type"`
	Data      []byte    `json:"data" db:"data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Константы для типов данных пользователя
const (
	UserDataTypeCredential = "credential" // Учетные данные (логин/пароль)
	UserDataTypeCard       = "card"       // Данные банковской карты
	UserDataTypeText       = "text"       // Произвольный текст
	UserDataTypeFile       = "file"       // Файл
)

type FileMetadata struct {
	FileName  string `json:"file_name"`
	Extension string `json:"extension"`
	URL       string `json:"url"`
}
