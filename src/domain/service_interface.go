package domain

import (
	"github.com/google/uuid"
)

type Interface interface {
	NextTurn(currentGame CurrentGame) (int, int, error)
	CheckField(currentGame CurrentGame, expectedGame CurrentGame) (bool, error)
	CheckFinish(currentGame CurrentGame) (int, int, bool)
}

// GameServiceInterface - структура, реализующая интерфейс Interface
type GameServiceInterface struct{}

// serService - интерфейс для работы с пользователями
type UserService interface {
	Registration(login, password string) (bool, error)
	// Authorization(login, password string) (uuid.UUID, error)
	Authenticate(authHeader string) (uuid.UUID, error)
	GetUserByID(id uuid.UUID) (*User, error)
}

// // сервис авторизации, который реализует интерфейс UserService
// type AutorizationService struct {
// }

// // Метод регистрации, который принимает SignUpRequest и возвращает факт успешной регистрации;
// func (*UserService) Registration(singUp SignUpRequest) bool {

// }

// // Метод авторизации, который принимает в заголовке логин и пароль в виде base64(login:password) и возвращает UUID пользователя.
// func (*UserService) Autorization(login, password base64.Encoding) uuid.UUID {

// }
