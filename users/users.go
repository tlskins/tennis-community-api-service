package users

import (
	"context"

	t "github.com/tennis-community-api-service/users/types"
)

func (u *UsersService) GetUser(_ context.Context, userID string) (resp *t.User, err error) {
	return u.Store.GetUser(userID)
}

func (u *UsersService) CreateUser(_ context.Context, data *t.User) (resp *t.User, err error) {
	return u.Store.CreateUser(data)
}

func (u *UsersService) UpdateUser(_ context.Context, data *t.UpdateUser) (resp *t.User, err error) {
	return u.Store.UpdateUser(data)
}
