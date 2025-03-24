package interfaces

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// В interfaces/db.go
type DB interface {
	QueryRow(query string, args ...any) *sqlx.Row
	Ping() error
	Exec(query string, args ...any) (sql.Result, error)
}
