package model

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/version-1/golang-rest-api/internal/model/entity"
)

type Executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func find[V entity.EntityScanner](ctx context.Context, ex Executor, qm Query, id int) (V, error) {
	e := qm.m.Entity()
	row := ex.QueryRowContext(ctx, qm.Find(id), id)
	if err := e.Scan(row); err != nil {
		if err == sql.ErrNoRows {
			return e.(V), NewErrorNotFound(qm.m.Table(), id)
		}
		return e.(V), err
	}

	return e.(V), nil
}

func findAll[V entity.EntityScanner](ctx context.Context, ex Executor, qm Query, id int) ([]V, error) {
	list := []V{}
	rows, err := ex.QueryContext(ctx, qm.FindAll(10))
	if err != nil {
		return list, err
	}

	for rows.Next() {
		item := qm.m.Entity()
		if err := item.Scan(rows); err != nil {
			return list, err
		}
		list = append(list, item.(V))
	}

	return list, nil
}

func create[V entity.EntityScanner](ctx context.Context, ex Executor, qm Query, args ...any) (V, error) {
	id := 0
	var res V
	if err := ex.QueryRowContext(ctx, qm.Create(), args...).Scan(&id); err != nil {
		return res, err
	}

	return find[V](ctx, ex, qm, id)
}

func update[V entity.EntityScanner](ctx context.Context, ex Executor, qm Query, id int, args ...any) (V, error) {
	var res V
	_args := []any{id}
	_args = append(_args, args...)
	if _, err := ex.ExecContext(ctx, qm.Update(), _args...); err != nil {
		return res, err
	}

	return find[V](ctx, ex, qm, id)
}

func destroy[V entity.EntityScanner](ctx context.Context, ex Executor, qm Query, id int) (V, error) {
	var res V
	if _, err := ex.ExecContext(ctx, qm.Delete(), id); err != nil {
		return res, err
	}

	return find[V](ctx, ex, qm, id)
}

func withOneByID[E entity.EntityScanner](ctx context.Context, m Executor, qm Query, with string, keys []string) ([]E, error) {
	list := []E{}
	rel, err := qm.With(with)
	if err != nil {
		return list, err
	}

	rows, err := m.QueryContext(ctx, rel.Query(), pq.Array(keys))
	if err != nil {
		return list, err
	}

	for rows.Next() {
		item := RelationEntity[E]{entity: rel.Build().(E)}
		if err := item.Scan(rows); err != nil {
			return list, err
		}
		list = append(list, item.Entity().(E))
	}

	return list, nil
}

func withManyByIDs[E entity.EntityScanner](ctx context.Context, ex Executor, qm Query, with string, keys []string) (map[string][]E, error) {
	maps := map[string][]E{}
	rel, err := qm.With(with)
	if err != nil {
		return maps, err
	}

	rows, err := ex.QueryContext(ctx, rel.Query(), pq.Array(keys))
	if err != nil {
		return maps, err
	}

	for rows.Next() {
		item := RelationEntity[E]{entity: rel.Build().(E)}
		if err := item.Scan(rows); err != nil {
			return maps, err
		}
		maps[item.SourceKey()] = append(maps[item.SourceKey()], item.Entity().(E))
	}

	return maps, nil
}
