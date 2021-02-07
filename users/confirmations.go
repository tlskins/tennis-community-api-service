package users

import (
	"context"

	t "github.com/tennis-community-api-service/users/types"
)

func (u *UsersService) GetConfirmation(_ context.Context, ID string) (conf *t.UserConfirmation, err error) {
	return u.Store.GetConfirmation(ID)
}

func (u *UsersService) CreateConfirmation(_ context.Context, data *t.UserConfirmation) (conf *t.UserConfirmation, err error) {
	return u.Store.CreateConfirmation(data)
}

func (u *UsersService) DeleteUserConfirmations(_ context.Context, userID, email string) (err error) {
	return u.Store.DeleteUserConfirmations(userID, email)
}
