package model

import (
	"context"
	"database/sql"
	"log"

	"github.com/version-1/golang-rest-api/internal/model/entity"
)

type Model struct {
	db     *sql.DB
	logger *log.Logger
}

func New(db *sql.DB) *Model {
	return &Model{db: db, logger: log.Default()}
}

func (m Model) User(tx ...Executor) *UserModel {
	if len(tx) > 0 {
		return &UserModel{m: tx[0]}
	}
	return &UserModel{m: m}
}

func (m Model) Post(tx ...Executor) *PostModel {
	if len(tx) > 0 {
		return &PostModel{m: tx[0]}
	}
	return &PostModel{m: m}
}

func (m Model) Tag(tx ...Executor) *TagModel {
	if len(tx) > 0 {
		return &TagModel{m: tx[0]}
	}
	return &TagModel{m: m}
}

func (m Model) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	m.logger.Printf("executing query: [ExecContext] %s", query)
	msg := []any{"args: [ExecContext]"}
	msg = append(msg, args...)
	m.logger.Println(msg...)
	return m.db.ExecContext(ctx, query, args...)
}

func (m Model) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	m.logger.Printf("executing query: [QueryContext] %s", query)
	msg := []any{"args: [QueryContext]"}
	msg = append(msg, args...)
	m.logger.Println(msg...)
	return m.db.QueryContext(ctx, query, args...)
}

func (m Model) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	m.logger.Printf("executing query: [QueryRowContext] %s", query)
	msg := []any{"args: [QueryRowContext]"}
	msg = append(msg, args...)
	m.logger.Println(msg...)
	return m.db.QueryRowContext(ctx, query, args...)
}

func (m Model) Transaction(ctx context.Context, cb func(tx Executor) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			m.logger.Printf("recovered from panic: %v", r)
			tx.Rollback()
		}
	}()
	if err := cb(tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

type relationScanner interface {
	entity.EntityScanner
	SourceKey() string
	Entity() entity.EntityScanner
}

type RelationEntity[V entity.EntityScanner] struct {
	entity    V
	sourceKey string
}

type key struct {
	v any
}

func (r *RelationEntity[V]) Scan(rows entity.Scanner, extra ...any) error {
	sourceKey := ""
	err := r.entity.Scan(rows, &sourceKey)
	if err != nil {
		return err
	}

	r.sourceKey = sourceKey
	return nil
}

func (r RelationEntity[V]) SourceKey() string {
	return r.sourceKey
}

func (r RelationEntity[V]) Entity() entity.EntityScanner {
	return r.entity
}
