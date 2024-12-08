package model

import (
	"fmt"

	"github.com/version-1/golang-rest-api/internal/model/entity"
)

type ManyToMany[T entity.EntityScanner] struct {
	build     func() T
	To        string
	Through   string
	SourceKey string
	TargetKey string
}

func (r ManyToMany[T]) Query() string {
	return fmt.Sprintf("SELECT %s.*, %s.%s FROM %s INNER JOIN %s ON %s.%s = %s.id WHERE %s.%s = any($1)",
		r.To,
		r.Through,
		r.TargetKey,
		r.To,
		r.Through,
		r.Through,
		r.SourceKey,
		r.To,
		r.Through,
		r.TargetKey,
	)
}

func (b ManyToMany[T]) Build() entity.EntityScanner {
	return b.build()
}

type OneToMany[T entity.EntityScanner] struct {
	build     func() T
	To        string
	TargetKey string
}

func (r OneToMany[T]) Query() string {
	return fmt.Sprintf("SELECT %s.* FROM %s WHERE %s.%s = any($1)",
		r.To,
		r.To,
		r.To,
		r.TargetKey,
	)
}

func (r OneToMany[T]) Build() entity.EntityScanner {
	return r.build()
}

type OneToOne[T entity.EntityScanner] struct {
	build     func() T
	To        string
	SourceKey string
	TargetKey string
}

func (r OneToOne[T]) Query() string {
	return fmt.Sprintf("SELECT %s.*, %s.%s FROM %s WHERE %s.%s = any($1) LIMIT 1",
		r.To,
		r.To,
		r.TargetKey,
		r.To,
		r.To,
		r.TargetKey,
	)
}

func (r OneToOne[T]) Build() entity.EntityScanner {
	return r.build()
}
