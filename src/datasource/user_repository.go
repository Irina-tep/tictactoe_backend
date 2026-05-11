package datasource

import (
	"context"
	"domain"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// UserRepository - интерфейс репозитория для работы с пользователями
type UserRepository interface {
	Save(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, userID uuid.UUID) (*UserData, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
}

type UserStorage struct {
	strorage *GameStorage
}

// type UserData struct {
// 	UserD     *domain.User
// 	// CreatedAt time.Time
// 	// UpdatedAt time.Time
// }

func NewUserStorage(s *GameStorage) *UserStorage {
	return &UserStorage{
		strorage: s,
	}
}

func (s *UserStorage) Save(ctx context.Context, u *domain.User) error {
	if u == nil {
		return errors.New("user cannot be nil")
	}
	if s.strorage.db == nil {
		return errors.New("database connection is not initialized")
	}
	// Проверяем, существует ли пользователь
	_, err := s.strorage.db.Exec(ctx, `INSERT INTO player (id, login, password, created_at, updated_at) 
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
		ON CONFLICT (id) DO NOTHING `,
		u.ID, u.Login, u.Password)
	if err != nil {
		return errors.New("failed to save user")
	}
	return nil

}

func (s *UserStorage) GetByID(ctx context.Context, userID uuid.UUID) (*UserData, error) {
	var user UserData
	// user.UserD = &domain.User{}
	if userID == uuid.Nil {
		return nil, errors.New("user ID is empty")
	}
	if s.strorage.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	// var createdAt, updatedAt time.Time
	// err := s.strorage.db.Query(ctx, "SELECT id, login, password, created_at, updated_at FROM player WHERE id = $1", userID)
	// if err == pgx.ErrNoRows {
	// 	return nil, errors.New("user not found")
	// }

	// if err != nil {
	// 	return nil, errors.New("failed to get user")
	// }

	row, err := s.strorage.db.Query(ctx, "SELECT id, login, password, created_at, updated_at FROM player WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}

	defer row.Close()
	user, err = pgx.CollectOneRow(row, pgx.RowToStructByName[UserData])
	if err != nil {
		return nil, ErrGameIDNotFound //поменять ошибку
	}

	return &user, nil
}

func (s *UserStorage) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	if s.strorage.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	err := s.strorage.db.QueryRow(ctx, "SELECT id, login, password, created_at, updated_at FROM player WHERE login = $1", login).Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, errors.New("failed to get user")
	}

	return &user, nil
}
