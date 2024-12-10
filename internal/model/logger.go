package model

import (
	"context"
	"database/sql"
	"log"
)

type queryLogger struct {
	ex     Executor
	logger *log.Logger
}

func (l queryLogger) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	l.logger.Printf("executing query: [ExecContext] %s", query)
	msg := []any{"args: [ExecContext]"}
	msg = append(msg, args...)
	l.logger.Println(msg...)
	return l.ex.ExecContext(ctx, query, args...)
}

func (l queryLogger) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	l.logger.Printf("executing query: [QueryContext] %s", query)
	msg := []any{"args: [QueryContext]"}
	msg = append(msg, args...)
	l.logger.Println(msg...)
	return l.ex.QueryContext(ctx, query, args...)
}

func (l queryLogger) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	l.logger.Printf("executing query: [QueryRowContext] %s", query)
	msg := []any{"args: [QueryRowContext]"}
	msg = append(msg, args...)
	l.logger.Println(msg...)
	return l.ex.QueryRowContext(ctx, query, args...)
}
