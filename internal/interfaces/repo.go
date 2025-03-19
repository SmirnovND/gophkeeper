package interfaces

import "github.com/SmirnovND/gophkeeper/internal/domain"

// UserRepo описывает интерфейс для работы с пользователями.
type UserRepo interface {
	// FindUser ищет пользователя по логину и возвращает его, если найден.
	// Возвращает ошибку, если пользователь не найден или произошла другая ошибка.
	FindUser(login string) (*domain.User, error)

	// SaveUser сохраняет нового пользователя в базе данных.
	// Возвращает ошибку, если произошла ошибка при сохранении.
	SaveUser(user *domain.User) error
}
