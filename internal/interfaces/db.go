package interfaces

import "github.com/jmoiron/sqlx"

// В interfaces/db.go
type DB interface {
	QueryRow(query string, args ...any) *sqlx.Row
	Ping() error
}
