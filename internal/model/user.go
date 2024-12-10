package model

import (
	"context"

	"github.com/version-1/golang-rest-api/internal/model/entity"
)

const DummyUserID = 1

type UserModel struct {
	m Executor
}

func (u UserModel) Query() Query {
	return Query{m: u}
}

func (u UserModel) Table() string {
	return "users"
}

func (u UserModel) Fields() []string {
	return []string{"username", "email"}
}

func (u UserModel) Entity() entity.EntityScanner {
	return &entity.User{}
}

func (u UserModel) Relationships() map[string]Relationship {
	return map[string]Relationship{}
}

func (u UserModel) Find(ctx context.Context, id int) (*entity.User, error) {
	return find[*entity.User](ctx, u.m, Query{m: u}, id)
}

func (u UserModel) FindAll(ctx context.Context) ([]*entity.User, error) {
	return findAll[*entity.User](ctx, u.m, Query{m: u}, 10)
}

func (u UserModel) Create(ctx context.Context, username, email string) (*entity.User, error) {
	return create[*entity.User](ctx, u.m, Query{m: u}, username, email)
}

func (u UserModel) Update(ctx context.Context, id int, username, email string) (*entity.User, error) {
	return update[*entity.User](ctx, u.m, Query{m: u}, id, username, email)
}

func (u UserModel) Delete(ctx context.Context, id int) (*entity.User, error) {
	return update[*entity.User](ctx, u.m, Query{m: u}, id)
}
