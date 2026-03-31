package repository

import (
	"api/internal/storage"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	store := storage.NewSQLxStorage(sqlxDB)

	repo := NewUserRepository(store)

	createdAt, _ := time.Parse("2006-01-02", "2024-01-01")

	rows := sqlmock.NewRows([]string{"id", "email", "created_at"}).
		AddRow(1, "test@example.com", createdAt)

	mock.ExpectQuery("SELECT id, email, created_at FROM users").WillReturnRows(rows)

	users, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "test@example.com", users[0].Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	store := storage.NewSQLxStorage(sqlxDB)
	repo := NewUserRepository(store)

	mock.ExpectQuery("SELECT id, email, created_at FROM users").
		WillReturnError(sqlmock.ErrCancelled)

	users, err := repo.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}
