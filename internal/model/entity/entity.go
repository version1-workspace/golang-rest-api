package entity

import "time"

type Scanner interface {
	Scan(dest ...any) error
}

type EntityScanner interface {
	Scan(Scanner, ...any) error
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Scan(rows Scanner, extra ...any) error {
	_args := []any{&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt}
	_args = append(_args, extra...)
	return rows.Scan(_args...)
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Tags []*Tag `json:"tags"`
	User *User  `json:"user"`
}

func (u *Post) Scan(rows Scanner, extra ...any) error {
	list := []any{&u.ID, &u.UserID, &u.Title, &u.Content, &u.CreatedAt, &u.UpdatedAt}
	list = append(list, extra...)
	return rows.Scan(list...)
}

type Tag struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Tag) Scan(rows Scanner, extra ...any) error {
	_args := []any{&t.ID, &t.Slug, &t.Name, &t.CreatedAt, &t.UpdatedAt}
	_args = append(_args, extra...)
	return rows.Scan(_args...)
}

type PostBody struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Tags    []TagBody `json:"tags"`
}

type TagBody struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}
