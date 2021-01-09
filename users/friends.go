package users

import (
	"context"
	"time"

	t "github.com/tennis-community-api-service/users/types"
)

func (u *UsersService) SendFriendRequest(_ context.Context, fromID, toID string) (err error) {
	return u.Store.SendFriendRequest(&t.FriendRequest{
		CreatedAt:  time.Now(),
		FromUserID: fromID,
		ToUserID:   toID,
	})
}

func (u *UsersService) AcceptFriendRequest(_ context.Context, acceptorID, reqID string, accept bool) (user *t.User, err error) {
	return u.Store.AcceptFriendRequest(acceptorID, reqID, accept)
}

func (u *UsersService) Unfriend(_ context.Context, sourceID, targetID string) (err error) {
	return u.Store.Unfriend(sourceID, targetID)
}

func (u *UsersService) SearchFriends(_ context.Context, search *string, IDs *[]string, limit, offset int) ([]*t.Friend, error) {
	return u.Store.SearchFriends(search, IDs, limit, offset)
}
