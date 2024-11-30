package model

import (
	"context"
	"database/sql"
)

type executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func find[V entityScanner](ctx context.Context, ex executor, qm Query, id int) (V, error) {
	e := qm.m.Entity()
	row := ex.QueryRowContext(ctx, qm.Find(id), id)
	if err := e.Scan(row); err != nil {
		return e.(V), err
	}

	return e.(V), nil
}

func findAll[V entityScanner](ctx context.Context, ex executor, qm Query, id int) ([]V, error) {
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

func create[V entityScanner](ctx context.Context, ex executor, qm Query, args ...any) (V, error) {
	id := 0
	var res V
	if err := ex.QueryRowContext(ctx, qm.Create(), args...).Scan(&id); err != nil {
		return res, err
	}

	return find[V](ctx, ex, qm, id)
}

func update[V entityScanner](ctx context.Context, ex executor, qm Query, id int, args ...any) (V, error) {
	var res V
	_args := []any{id}
	_args = append(_args, args...)
	if _, err := ex.ExecContext(ctx, qm.Update(), _args...); err != nil {
		return res, err
	}

	return find[V](ctx, ex, qm, id)
}

func destroy[V entityScanner](ctx context.Context, ex executor, qm Query, id int) (V, error) {
	var res V
	if _, err := ex.ExecContext(ctx, qm.Delete(), id); err != nil {
		return res, err
	}

	return find[V](ctx, ex, qm, id)
}
