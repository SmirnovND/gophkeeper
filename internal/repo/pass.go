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

// DeleteUserData удаляет данные пользователя по ID
func (r *UserDataRepo) DeleteUserData(id string) error {
	query := `DELETE FROM "user_data" WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user data: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetUserDataByLabelAndType ищет данные пользователя по метке и типу
func (r *UserDataRepo) GetUserDataByLabelAndType(userID, label string, dataType string) (*domain.UserData, error) {
	query := `SELECT id, user_id, label, type, data, created_at, updated_at
              FROM "user_data"
              WHERE user_id = $1 AND label = $2 AND type = $3
              LIMIT 1`
	row := r.db.QueryRow(query, userID, label, dataType)

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
		return nil, fmt.Errorf("error querying user data by label and type: %w", err)
	}

	return userData, nil
}
