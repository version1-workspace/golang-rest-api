package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/version-1/golang-rest-api/internal/model/entity"
)

type schema interface {
	Table() string
	Fields() []string
	Entity() entity.EntityScanner
	Relationships() map[string]Relationship
}

type Query struct {
	m schema
}

func (q Query) Find(id int) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE id = $1", q.m.Table())
}

func (q Query) FindAll(limit int) string {
	return fmt.Sprintf("SELECT * FROM %s limit %d", q.m.Table(), limit)
}

func (q Query) With(name string) (Relationship, error) {
	relationship, ok := q.m.Relationships()[name]
	if !ok {
		return nil, fmt.Errorf("relationship %s not found", name)
	}
	return relationship, nil
}

func (q Query) Create() string {
	fields := strings.Join(q.m.Fields(), ", ")
	placeholders := buildPlaceholders(len(q.m.Fields()))

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", q.m.Table(), fields, placeholders)
}

func (q Query) Update() string {
	sets := []string{}
	for i, field := range q.m.Fields() {
		placeholder := fmt.Sprintf("$%d", i+2)
		sets = append(sets, fmt.Sprintf("%s = %s", field, placeholder))
	}

	return fmt.Sprintf("UPDATE %s SET %s, updated_at = NOW() WHERE id = $1", q.m.Table(), sets)
}

func (q Query) Delete() string {
	return fmt.Sprintf("DELETE FROM %s WHERE id = $1", q.m.Table())
}

func buildPlaceholders(n int) string {
	if n == 0 {
		return ""
	}

	placeholders := []string{}
	for i := 0; i < n; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}

	return strings.Join(placeholders, ", ")
}

type Relationship interface {
	Query() string
	Build() entity.EntityScanner
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
