package model

import (
	"context"
	"strconv"

	"github.com/version-1/golang-rest-api/internal/model/entity"
)

type PostModel struct {
	m Executor
}

func (p PostModel) Query() Query {
	return Query{m: p}
}

func (p PostModel) Table() string {
	return "posts"
}

func (p PostModel) Fields() []string {
	return []string{"user_id", "title", "content"}
}

func (p PostModel) Entity() entity.EntityScanner {
	return &entity.Post{}
}

func (p PostModel) Relationships() map[string]Relationship {
	return map[string]Relationship{
		"tags": ManyToMany[*entity.Tag]{
			build: func() *entity.Tag {
				return &entity.Tag{}
			},
			To:        "tags",
			Through:   "tag_posts",
			SourceKey: "tag_id",
			TargetKey: "post_id",
		},
		"user": OneToOne[*entity.User]{
			build: func() *entity.User {
				return &entity.User{}
			},
			To:        "users",
			TargetKey: "id",
		},
	}
}

func (p PostModel) Find(ctx context.Context, id int) (*entity.Post, error) {
	q := p.Query()
	m, err := find[*entity.Post](ctx, p.m, q, id)
	if err != nil {
		return m, err
	}

	tags, err := withOneByID[*entity.Tag](ctx, p.m, q, "tags", []string{strconv.Itoa(m.ID)})
	if err != nil {
		return m, err
	}

	m.Tags = tags

	users, err := withOneByID[*entity.User](ctx, p.m, q, "user", []string{strconv.Itoa(m.UserID)})
	if err != nil {
		return m, err
	}

	if len(users) > 0 {
		m.User = users[0]
	}

	return m, err
}

func (p PostModel) FindAll(ctx context.Context, limit int) ([]*entity.Post, error) {
	q := p.Query()
	list, err := findAll[*entity.Post](ctx, p.m, q, limit)
	if err != nil {

		return list, err
	}

	b1 := newBatchLoader[[]*entity.Tag]()
	for _, item := range list {
		b1.Add(strconv.Itoa(item.ID), func(v []*entity.Tag) {
			item.Tags = v
		})
	}
	b1.SetBatch(func(keys []string) (map[string][]*entity.Tag, error) {
		return withManyByIDs[*entity.Tag](ctx, p.m, q, "tags", keys)
	})

	if err := b1.Load(); err != nil {
		return list, err
	}

	b2 := newBatchLoader[[]*entity.User]()
	for _, item := range list {
		b2.Add(strconv.Itoa(item.UserID), func(v []*entity.User) {
			if len(v) > 0 {
				item.User = v[0]
			}
		})
	}
	b2.SetBatch(func(keys []string) (map[string][]*entity.User, error) {
		return withManyByIDs[*entity.User](ctx, p.m, q, "user", keys)
	})

	if err := b2.Load(); err != nil {
		return list, err
	}

	return list, nil
}

func (p PostModel) Create(ctx context.Context, userID int, title, content string) (*entity.Post, error) {
	return create[*entity.Post](ctx, p.m, Query{m: p}, userID, title, content)
}

func (p PostModel) Update(ctx context.Context, id int, title, content string) (*entity.Post, error) {
	return update[*entity.Post](ctx, p.m, p.Query(), id, title, content)
}

func (p PostModel) Delete(ctx context.Context, id int) (*entity.Post, error) {
	post, err := destroy[*entity.Post](ctx, p.m, Query{m: p}, id)
	if err != nil {
		return post, err
	}

	if _, err = p.m.ExecContext(ctx, "DELETE FROM tag_posts WHERE post_id = $1", id); err != nil {
		return post, err
	}

	return post, err
}
