package users

import (
	"context"
	"time"

	t "github.com/tennis-community-api-service/users/types"
)

func (u *UsersService) SendFriendRequest(_ context.Context, fromID, toID string) (err error) {
	req := &t.FriendRequest{
		CreatedAt:  time.Now(),
		FromUserID: fromID,
		ToUserID:   toID,
	}
	note := &t.FriendNote{
		CreatedAt: time.Now(),
		Subject:   "New Friend Request",
		Type:      "New Friend Request",
		FriendID:  fromID,
	}
	return u.Store.SendFriendRequest(req, note)
}

func (u *UsersService) AcceptFriendRequest(_ context.Context, acceptorID, reqID string, accept bool) (user *t.User, err error) {
	var note *t.FriendNote
	if accept {
		note = &t.FriendNote{
			CreatedAt: time.Now(),
			Subject:   "New Friend",
			Type:      "New Friend",
			FriendID:  acceptorID,
		}
	}

	return u.Store.AcceptFriendRequest(acceptorID, reqID, accept, note)
}

func (u *UsersService) Unfriend(_ context.Context, sourceID, targetID string) (err error) {
	return u.Store.Unfriend(sourceID, targetID)
}

func (u *UsersService) SearchFriends(_ context.Context, search *string, IDs *[]string, limit, offset int) ([]*t.Friend, error) {
	return u.Store.SearchFriends(search, IDs, limit, offset)
}
