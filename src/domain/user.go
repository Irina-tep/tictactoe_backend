package domain

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Login    string
	Password string
	Symbol   int
}

func NewUser(login, password string) (User, error) {
	if err := validateLogin(login); err != nil {
		return User{}, err
	}

	if err := validatePassword(password); err != nil {
		return User{}, err
	}

	user :=
		User{
			ID:       uuid.New(),
			Login:    login,
			Password: password, // В реальном приложении нужно хешировать!
		}
	return user, nil
}

func validateLogin(login string) error {
	if len(login) < 3 || len(login) > 20 {
		return errors.New("Invalid login")
	}
	expected, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, login)
	if !expected {
		return errors.New("Invalid login")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("Invalid password")
	}
	return nil
}
