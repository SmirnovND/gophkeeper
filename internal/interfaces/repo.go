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

// UserDataRepo описывает интерфейс для работы с данными пользователя.
type UserDataRepo interface {
	// SaveUserData сохраняет данные пользователя в базе данных.
	// Возвращает ошибку, если произошла ошибка при сохранении.
	SaveUserData(userData *domain.UserData) error

	// FindUserDataByLabel ищет данные пользователя по метке.
	// Возвращает данные и nil, если данные найдены.
	// Возвращает nil и ошибку, если данные не найдены или произошла другая ошибка.
	FindUserDataByLabel(userID, label string) (*domain.UserData, error)

	// GetUserDataByLabelAndType ищет данные пользователя по метке и типу.
	// Возвращает данные и nil, если данные найдены.
	// Возвращает nil и ошибку, если данные не найдены или произошла другая ошибка.
	GetUserDataByLabelAndType(userID, label string, dataType string) (*domain.UserData, error)

	// FindAllUserData возвращает все данные пользователя.
	// Возвращает список данных и nil, если данные найдены.
	// Возвращает nil и ошибку, если данные не найдены или произошла другая ошибка.
	FindAllUserData(userID string) ([]*domain.UserData, error)

	// DeleteUserData удаляет данные пользователя по ID.
	// Возвращает ошибку, если произошла ошибка при удалении.
	DeleteUserData(id string) error
}

// TokenStorage описывает интерфейс для хранения и управления токеном авторизации.
type TokenStorage interface {
	// SaveToken сохраняет токен.
	SaveToken(token string) error

	// LoadToken загружает токен из файла.
	LoadToken() (string, error)
}
