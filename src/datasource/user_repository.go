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
	GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
}

type UserRepositoryStruct struct {
	strorage *Storage
}

func NewUserRepositoryStruct(s *Storage) *UserRepositoryStruct {
	return &UserRepositoryStruct{
		strorage: s,
	}
}

func (r *UserRepositoryStruct) Save(ctx context.Context, u *domain.User) error {
	if u == nil {
		return errors.New("user cannot be nil")
	}
	if r.strorage.db == nil {
		return errors.New("database connection is not initialized")
	}
	// Проверяем, существует ли пользователь
	_, err := r.strorage.db.Exec(ctx, `INSERT INTO player (id, login, password, created_at, updated_at) 
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
		ON CONFLICT (id) DO NOTHING `,
		u.ID, u.Login, u.Password)
	if err != nil {
		return errors.New("failed to save user")
	}
	return nil

}

func (r *UserRepositoryStruct) GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	var user domain.User

	if userID == uuid.Nil {
		return nil, errors.New("user ID is empty")
	}
	if r.strorage.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	err := r.strorage.db.QueryRow(ctx, "SELECT id, login, password FROM player WHERE id = $1", userID).Scan(&user.ID, &user.Login, &user.Password)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, errors.New("failed to get user")
	}

	return &user, nil
}

func (r *UserRepositoryStruct) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	if r.strorage.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	err := r.strorage.db.QueryRow(ctx, "SELECT id, login, password FROM player WHERE login = $1", login).Scan(&user.ID, &user.Login, &user.Password)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, errors.New("failed to get user")
	}

	return &user, nil
}
