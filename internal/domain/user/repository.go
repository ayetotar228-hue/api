package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetAll(ctx context.Context) ([]User, error)
}
