package model

import (
	"fmt"
	"strings"

	"github.com/version-1/golang-rest-api/internal/model/entity"
)

type schema interface {
	Table() string
	Fields() []string
	Entity() entity.EntityScanner
	Relationships() map[string]Relationship
}

type Relationship interface {
	Query() string
	Build() entity.EntityScanner
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

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURN id", q.m.Table(), fields, placeholders)
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
