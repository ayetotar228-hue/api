package repository

import (
	"api/internal/domain/user"
	"api/pkg/storage"
	"context"
)

type UserRepository struct {
	storage storage.Storager
}

func NewUserRepository(storage storage.Storager) *UserRepository {
	return &UserRepository{storage: storage}
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	db := r.storage.GetDB()

	err := db.QueryRowContext(
		ctx,
		"INSERT INTO users (email) VALUES ($1) RETURNING id, email, created_at",
		user.Email,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	return err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]user.User, error) {
	db := r.storage.GetDB()

	rows, err := db.QueryContext(ctx, "SELECT id, email, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var user user.User
		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}
