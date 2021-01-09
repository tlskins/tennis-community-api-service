package users

import (
	"context"
	"time"

	t "github.com/tennis-community-api-service/users/types"
)

func (u *UsersService) GetUser(_ context.Context, userID string) (resp *t.User, err error) {
	return u.Store.GetUser(userID)
}

func (u *UsersService) GetUserByEmail(_ context.Context, email string) (resp *t.User, err error) {
	return u.Store.GetUserByEmail(email)
}

func (u *UsersService) CreateUser(_ context.Context, data *t.User) (resp *t.User, err error) {
	return u.Store.CreateUser(data)
}

func (u *UsersService) UpdateUser(_ context.Context, data *t.UpdateUser) (resp *t.User, err error) {
	return u.Store.UpdateUser(data)
}

func (u *UsersService) ClearUserNotifications(_ context.Context, id string, uploads, friends bool) (resp *t.User, err error) {
	now := time.Now()
	update := &t.UpdateUser{
		ID:        id,
		UpdatedAt: &now,
	}
	if uploads {
		empty := []*t.UploadNote{}
		update.UploadNotes = &empty
	}
	if friends {
		empty := []*t.FriendNote{}
		update.FriendNotes = &empty
	}
	return u.Store.UpdateUser(update)
}

func (u *UsersService) AddUploadNotifications(_ context.Context, id string, note *t.UploadNote) (resp *t.User, err error) {
	return u.Store.AddUploadNote(id, note)
}
