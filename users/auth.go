package users

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/tennis-community-api-service/pkg/auth"
	t "github.com/tennis-community-api-service/users/types"
)

func (u *UsersService) ConfirmUser(_ context.Context, userID string) (resp *t.User, err error) {
	conf := true
	return u.Store.UpdateUser(&t.UpdateUser{
		ID:        userID,
		Confirmed: &conf,
	})
}

func (u *UsersService) SignIn(_ context.Context, email, pwd string) (user *t.User, err error) {
	if user, err = u.Store.GetUserByEmail(email); err != nil {
		return
	}
	if err = auth.ValidateCredentials(user.EncryptedPassword, pwd); err != nil {
		return nil, errors.Wrap(err, "error validating credentials")
	}

	now := time.Now()
	if user, err = u.Store.UpdateUser(&t.UpdateUser{ID: user.ID, LastLoggedIn: &now}); err != nil {
		return nil, errors.Wrap(err, "error updating user")
	}
	return
}
