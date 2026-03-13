package repository

import (
	"api/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Select(&users, "SELECT id, email, created_at FROM users")
	return users, err
}
