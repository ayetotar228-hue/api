package storage

import (
	"context"
	"database/sql"
)

type MockStorage struct {
	DB *MockDBTX
}

func NewMockStorage() *MockStorage {
	return &MockStorage{DB: &MockDBTX{}}
}

func (m *MockStorage) GetDB() DBTX {
	return m.DB
}

func (m *MockStorage) Close() error {
	return nil
}

type MockDBTX struct {
	ExecContextFunc     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContextFunc    func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContextFunc func(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func (m *MockDBTX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if m.ExecContextFunc != nil {
		return m.ExecContextFunc(ctx, query, args...)
	}
	return nil, nil
}

func (m *MockDBTX) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if m.QueryContextFunc != nil {
		return m.QueryContextFunc(ctx, query, args...)
	}
	return nil, nil
}

func (m *MockDBTX) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if m.QueryRowContextFunc != nil {
		return m.QueryRowContextFunc(ctx, query, args...)
	}
	return nil
}
