package model

import (
	"context"
	"database/sql"
	"log"
)

type Model struct {
	db     *sql.DB
	logger *log.Logger
}

func New(db *sql.DB) Model {
	return Model{db: db, logger: &log.Logger{}}
}

func (m Model) User() *UserModel {
	return &UserModel{m: m}
}

func (m Model) Post() *PostModel {
	return &PostModel{m: m}
}

func (m Model) Tag() *TagModel {
	return &TagModel{m: m}
}

func (m Model) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m Model) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m Model) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return m.db.QueryRowContext(ctx, query, args...)
}

type scanner interface {
	Scan(dest ...interface{}) error
}

type entityScanner interface {
	Scan(scanner) error
}
