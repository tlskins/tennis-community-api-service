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

func (u *UsersService) UpdateUserProfile(_ context.Context, data *t.UpdateUserProfile) (resp *t.User, err error) {
	return u.Store.UpdateUserProfile(data)
}

func (u *UsersService) AddUploadNotifications(_ context.Context, id string, note *t.UploadNote) (resp *t.User, err error) {
	return u.Store.AddUploadNote(id, note)
}

func (u *UsersService) RemoveUploadNote(_ context.Context, userID, noteID string) (resp *t.User, err error) {
	return u.Store.RemoveUploadNote(userID, noteID)
}

func (u *UsersService) RemoveFriendNote(_ context.Context, userID, noteID string) (resp *t.User, err error) {
	return u.Store.RemoveFriendNote(userID, noteID)
}

func (u *UsersService) RemoveCommentNote(_ context.Context, userID, noteID string) (resp *t.User, err error) {
	return u.Store.RemoveCommentNote(userID, noteID)
}

func (u *UsersService) RemoveAlbumUserTagNote(_ context.Context, userID, noteID string) (resp *t.User, err error) {
	return u.Store.RemoveAlbumUserTagNote(userID, noteID)
}

func (u *UsersService) RecentUsers(_ context.Context, start, end time.Time, limit, offset int) ([]*t.User, error) {
	return u.Store.RecentUsers(start, end, limit, offset)
}
