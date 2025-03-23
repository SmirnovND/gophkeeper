package domain

import (
	"encoding/json"
	"time"
)

// UserData представляет собой структуру для хранения данных пользователя
type UserData struct {
	ID        string          `json:"id" db:"id"`
	UserID    string          `json:"user_id" db:"user_id"`
	Label     string          `json:"label" db:"label"`
	Type      string          `json:"type" db:"type"`
	Data      json.RawMessage `json:"data" db:"data"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// FileMetadata представляет собой структуру для хранения метаданных файла
type FileMetadata struct {
	FileName  string `json:"file_name"`
	Extension string `json:"extension"`
}