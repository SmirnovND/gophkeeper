package domain

import "errors"

// Константы для типов данных пользователя
const (
	UserDataTypeCredential = "credential" // Тип данных для учетных данных (логин/пароль)
	UserDataTypeCard       = "card"       // Тип данных для кредитных карт
	UserDataTypeText       = "text"       // Тип данных для текстовых данных
	UserDataTypeFile       = "file"       // Тип данных для файлов
)

// Ошибки
var (
	ErrNotFound = errors.New("not found")
)