package storage

import (
	"github.com/jmoiron/sqlx"
)

type SQLxStorage struct {
	db *sqlx.DB
}

func NewSQLxStorage(db *sqlx.DB) *SQLxStorage {
	return &SQLxStorage{db: db}
}

func (s *SQLxStorage) GetDB() DBTX {
	return s.db.DB
}

func (s *SQLxStorage) Close() error {
	return s.db.Close()
}
