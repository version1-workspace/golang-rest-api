package model

import (
	"context"
	"time"
)

type TagModel struct {
	m Model
}

type TagEntity struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *TagEntity) Scan(rows scanner) error {
	return rows.Scan(&t.ID, &t.Name, &t.PostID, &t.CreatedAt, &t.UpdatedAt)
}

func (t TagModel) Table() string {
	return "tags"
}

func (t TagModel) Fields() []string {
	return []string{"name", "post_id"}
}

func (t TagModel) Entity() entityScanner {
	return &TagEntity{}
}

func (u TagModel) Find(ctx context.Context, id int) (*TagEntity, error) {
	return find[*TagEntity](ctx, u.m, Query{m: u}, id)
}

func (u TagModel) FindAll(ctx context.Context) ([]*TagEntity, error) {
	return findAll[*TagEntity](ctx, u.m, Query{m: u}, 10)
}

func (u TagModel) Create(ctx context.Context, name string, postID int) (*TagEntity, error) {
	return create[*TagEntity](ctx, u.m, Query{m: u}, name, postID)
}

func (u TagModel) Update(ctx context.Context, id int, name string) (*TagEntity, error) {
	return update[*TagEntity](ctx, u.m, Query{m: u}, id, name)
}

func (u TagModel) Delete(ctx context.Context, id int) (*TagEntity, error) {
	return update[*TagEntity](ctx, u.m, Query{m: u}, id)
}
