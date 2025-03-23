package repo

import (
	"database/sql"
	"fmt"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
)

// UserDataRepo реализует интерфейс interfaces.UserDataRepo
type UserDataRepo struct {
	db interfaces.DB
}

// NewUserDataRepo создает новый экземпляр UserDataRepo
func NewUserDataRepo(db interfaces.DB) interfaces.UserDataRepo {
	return &UserDataRepo{
		db: db,
	}
}

// SaveUserData сохраняет данные пользователя в базе данных
func (r *UserDataRepo) SaveUserData(userData *domain.UserData) error {
	exec := r.db.QueryRow

	// Запрос с RETURNING id, чтобы получить вставленный id
	query := `INSERT INTO "user_data" (user_id, label, type, data)
              VALUES ($1, $2, $3, $4)
              RETURNING id, created_at, updated_at`

	// Выполняем запрос и получаем id, created_at, updated_at
	err := exec(query, userData.UserID, userData.Label, userData.Type, userData.Data).
		Scan(&userData.ID, &userData.CreatedAt, &userData.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error saving user data: %w", err)
	}

	return nil
}

// FindUserDataByLabel ищет данные пользователя по метке
func (r *UserDataRepo) FindUserDataByLabel(userID, label string) (*domain.UserData, error) {
	query := `SELECT id, user_id, label, type, data, created_at, updated_at
              FROM "user_data"
              WHERE user_id = $1 AND label = $2
              LIMIT 1`
	row := r.db.QueryRow(query, userID, label)

	userData := &domain.UserData{}
	err := row.Scan(
		&userData.ID,
		&userData.UserID,
		&userData.Label,
		&userData.Type,
		&userData.Data,
		&userData.CreatedAt,
		&userData.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error querying user data: %w", err)
	}

	return userData, nil
}

// FindAllUserData возвращает все данные пользователя
func (r *UserDataRepo) FindAllUserData(userID string) ([]*domain.UserData, error) {
	// Здесь должен быть код для получения всех данных пользователя
	// Но так как мы не используем этот метод в текущей задаче, оставим его реализацию на будущее
	return nil, fmt.Errorf("method not implemented")
}

// DeleteUserData удаляет данные пользователя по ID
func (r *UserDataRepo) DeleteUserData(id string) error {
	// Здесь должен быть код для удаления данных пользователя
	// Но так как мы не используем этот метод в текущей задаче, оставим его реализацию на будущее
	return fmt.Errorf("method not implemented")
}
