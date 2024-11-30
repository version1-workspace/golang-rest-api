package model

import (
	"context"
	"time"
)

type UserModel struct {
	m Model
}

var _ entityScanner = &UserEntity{}

type UserEntity struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserEntity) Scan(rows scanner) error {
	return rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)
}

func (u UserModel) Table() string {
	return "users"
}

func (u UserModel) Fields() []string {
	return []string{"username", "email"}
}

func (u UserModel) Entity() entityScanner {
	return &UserEntity{}
}

func (u UserModel) Find(ctx context.Context, id int) (*UserEntity, error) {
	return find[*UserEntity](ctx, u.m, Query{m: u}, id)
}

func (u UserModel) FindAll(ctx context.Context) ([]*UserEntity, error) {
	return findAll[*UserEntity](ctx, u.m, Query{m: u}, 10)
}

func (u UserModel) Create(ctx context.Context, username, email string) (*UserEntity, error) {
	return create[*UserEntity](ctx, u.m, Query{m: u}, username, email)
}

func (u UserModel) Update(ctx context.Context, id int, username, email string) (*UserEntity, error) {
	return update[*UserEntity](ctx, u.m, Query{m: u}, id, username, email)
}

func (u UserModel) Delete(ctx context.Context, id int) (*UserEntity, error) {
	return update[*UserEntity](ctx, u.m, Query{m: u}, id)
}
