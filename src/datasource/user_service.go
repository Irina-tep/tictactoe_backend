package datasource

import (
	"context"
	"domain"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// Создай сервис авторизации, который использует UserService
// UserServiceImpl - реализация UserService
type AuthorizationService struct {
	userRepo UserRepository
}

func NewAutorizationService(uRepo UserRepository) *AuthorizationService {
	return &AuthorizationService{
		userRepo: uRepo,
	}
}

func (authS *AuthorizationService) Registration(login, password string) (bool, error) {
	ctx := context.Background()
	// Создаем доменного пользователя
	user, err := domain.NewUser(login, password)
	if err != nil {
		return false, err
	}
	// Сохраняем в БД
	if err := authS.userRepo.Save(ctx, &user); err != nil {
		return false, err
	}
	return true, nil
}

// Метод авторизации, который принимает в заголовке логин и пароль в виде base64(login:password) и возвращает UUID пользователя.
func (authS *AuthorizationService) Authenticate(authHeader string) (uuid.UUID, error) {
	if authHeader == "" {
		return uuid.Nil, errors.New("authorization header is empty")
	}
	// Проверяем формат "Basic base64(login:password)"
	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoder, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return uuid.Nil, errors.New("failed to decode base64")
	}
	// Разбираем login:password
	parts := strings.SplitN(string(decoder), ":", 2)
	if len(parts) != 2 {
		return uuid.Nil, errors.New("invalid format login:password")
	}
	login := parts[0]
	password := parts[1]

	ctx := context.Background()
	user, err := authS.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return uuid.Nil, err
	}
	// Проверяем пароль (в реальном приложении сравниваем хеши!)
	if user.Password != password {
		return uuid.Nil, errors.New("invalid password")
	}
	return user.ID, nil
}
