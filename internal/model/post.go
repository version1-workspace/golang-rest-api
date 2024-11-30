package model

import (
	"context"
	"time"
)

type PostModel struct {
	m Model
}

func (p PostModel) Table() string {
	return "posts"
}

func (p PostModel) Fields() []string {
	return []string{"user_id", "title", "content"}
}

func (p PostModel) Entity() entityScanner {
	return &PostEntity{}
}

var _ entityScanner = &PostEntity{}

type PostEntity struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *PostEntity) Scan(rows scanner) error {
	return rows.Scan(&u.ID, &u.UserID, &u.Title, &u.Content, &u.CreatedAt, &u.UpdatedAt)
}

func (p PostModel) Find(ctx context.Context, id int) (*PostEntity, error) {
	return find[*PostEntity](ctx, p.m, Query{m: p}, id)
}

func (p PostModel) FindAll(ctx context.Context) ([]*PostEntity, error) {
	return findAll[*PostEntity](ctx, p.m, Query{m: p}, 10)
}

func (p PostModel) Create(ctx context.Context, userID int, title, content string) (*PostEntity, error) {
	return create[*PostEntity](ctx, p.m, Query{m: p}, userID, title, content)
}

func (p PostModel) Update(ctx context.Context, id, userID int, title, content string) (*PostEntity, error) {
	return update[*PostEntity](ctx, p.m, Query{m: p}, userID, title, content, id)
}

func (p PostModel) Delete(ctx context.Context, id int) (*PostEntity, error) {
	return destroy[*PostEntity](ctx, p.m, Query{m: p}, id)
}
