package model

import (
	"context"
	"database/sql"

	"github.com/version-1/golang-rest-api/internal/model/entity"
)

type TagModel struct {
	m Executor
}

func (t TagModel) Table() string {
	return "tags"
}

func (t TagModel) Fields() []string {
	return []string{"slug", "name"}
}

func (t TagModel) Entity() entity.EntityScanner {
	return &entity.Tag{}
}

func (t TagModel) Relationships() map[string]Relationship {
	return map[string]Relationship{}
}

func (u TagModel) Find(ctx context.Context, id int) (*entity.Tag, error) {
	return find[*entity.Tag](ctx, u.m, Query{m: u}, id)
}

func (u TagModel) FindAll(ctx context.Context) ([]*entity.Tag, error) {
	return findAll[*entity.Tag](ctx, u.m, Query{m: u}, 10)
}

func (u TagModel) Create(ctx context.Context, slug, name string) (*entity.Tag, error) {
	return create[*entity.Tag](ctx, u.m, Query{m: u}, slug, name)
}

func (u TagModel) Update(ctx context.Context, id int, name string) (*entity.Tag, error) {
	return update[*entity.Tag](ctx, u.m, Query{m: u}, id, name)
}

func (u TagModel) DetachAll(ctx context.Context, postID int) error {
	if _, err := u.m.ExecContext(ctx, "DELETE FROM tag_posts WHERE post_id = $1", postID); err != nil {
		return err
	}

	return nil
}

func (u TagModel) Attach(ctx context.Context, postID int, slug, name string) (*entity.Tag, error) {
	findQuery := "SELECT * FROM tags WHERE slug = $1 LIMIT 1"
	tag := &entity.Tag{}
	err := tag.Scan(u.m.QueryRowContext(ctx, findQuery, slug))
	if err != nil && err != sql.ErrNoRows {
		return tag, err
	}

	if err == sql.ErrNoRows {
		tag, err = u.Create(ctx, slug, name)
		if err != nil {
			return nil, err
		}
	}

	attachQuery := "INSERT INTO tag_posts (post_id, tag_id) VALUES($1, $2)"
	_, err = u.m.ExecContext(ctx, attachQuery, postID, tag.ID)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (u TagModel) UpsertBySlug(ctx context.Context, slug, name string, postID int) (*entity.Tag, error) {
	query := "SELECT * FROM tags WHERE slug = $1 LIMIT 1"
	tag := &entity.Tag{}
	err := tag.Scan(u.m.QueryRowContext(ctx, query, slug))
	if err != sql.ErrNoRows {
		return tag, err
	}

	if err == sql.ErrNoRows {
		return u.Create(ctx, slug, name)
	}

	return u.Update(ctx, tag.ID, name)
}

func (u TagModel) Delete(ctx context.Context, id int) (*entity.Tag, error) {
	return update[*entity.Tag](ctx, u.m, Query{m: u}, id)
}
