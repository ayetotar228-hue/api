package user

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type User struct {
	ID        int
	Email     string
	CreatedAt time.Time
}

func NewUser(email string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}

	return &User{
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}, nil
}
