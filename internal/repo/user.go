package repo

import (
	"database/sql"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
)

type UserRepo struct {
	db interfaces.DB
}

func NewUserRepo(db interfaces.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) FindUser(login string) (*domain.User, error) {
	query := `SELECT id, login, pass_hash FROM "user" WHERE login = $1 LIMIT 1`
	row := r.db.QueryRow(query, login)

	user := &domain.User{}
	err := row.Scan(&user.Id, &user.Login, &user.PassHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) SaveUser(user *domain.User) error {
	exec := r.db.QueryRow

	// Запрос с RETURNING id, чтобы получить вставленный id
	query := `INSERT INTO "user" (login, pass_hash) VALUES ($1, $2) RETURNING id`

	// Выполняем запрос и получаем id
	err := exec(query, user.Login, user.PassHash).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}
