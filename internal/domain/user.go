package domain

// User представляет собой структуру для хранения информации о пользователе
type User struct {
	Id string `json:"id" db:"id"`
	Credentials
}
