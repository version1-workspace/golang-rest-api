package model

import (
	"fmt"
	"strings"
)

type queryModel interface {
	Table() string
	Fields() []string
	Entity() entityScanner
}

type Query struct {
	m queryModel
}

func (q Query) Find(id int) string {
	return fmt.Sprintf("SELECT * FROM %s WHERE id = ?", q.m.Table())
}

func (q Query) FindAll(limit int) string {
	return fmt.Sprintf("SELECT * FROM %s limit %d", q.m.Table(), limit)
}

func (q Query) Create() string {
	fields := strings.Join(q.m.Fields(), ", ")
	placeholders := buildPlaceholders(len(q.m.Fields()))

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURN id", q.m.Table(), fields, placeholders)
}

func (q Query) Update() string {
	sets := []string{}
	for _, field := range q.m.Fields() {
		sets = append(sets, fmt.Sprintf("%s = ?", field))
	}

	return fmt.Sprintf("UPDATE %s SET %s, updated_at = NOW() WHERE id = ?", q.m.Table(), sets)
}

func (q Query) Delete() string {
	return fmt.Sprintf("DELETE FROM %s WHERE id = ?", q.m.Table())
}

func buildPlaceholders(n int) string {
	if n == 0 {
		return ""
	}

	return strings.Repeat("?, ", n-1) + "?"
}
