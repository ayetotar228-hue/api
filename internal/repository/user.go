package repository

import (
	"api/internal/models"
	"api/internal/storage"
	"context"
)

type UserRepository struct {
	storage storage.Storager
}

func NewUserRepository(storage storage.Storager) *UserRepository {
	return &UserRepository{storage: storage}
}

func (r *UserRepository) Create(ctx context.Context, email string) (*models.User, error) {
	db := r.storage.GetDB()

	var user models.User
	err := db.QueryRowContext(
		ctx,
		"INSERT INTO users (email) VALUES ($1) RETURNING id, email, created_at",
		email,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	db := r.storage.GetDB()

	rows, err := db.QueryContext(ctx, "SELECT id, email, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}
