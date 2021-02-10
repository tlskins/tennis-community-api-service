package users

import (
	"context"
	"net/http"

	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-users/types"
)

func (u *UCService) SendFriendRequest(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	friendID := r.PathParameters["friendId"]
	err = u.usr.SendFriendRequest(ctx, claims.Subject, friendID)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r.Headers, nil, http.StatusOK)
}

func (u *UCService) AcceptFriendRequest(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := &t.AcceptFriendReq{}
	api.Parse(r, req)
	reqID := r.PathParameters["reqId"]
	user, err := u.usr.AcceptFriendRequest(ctx, claims.Subject, reqID, req.Accept)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r.Headers, user, http.StatusOK)
}

func (u *UCService) Unfriend(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	friendID := r.PathParameters["friendId"]
	err = u.usr.Unfriend(ctx, claims.Subject, friendID)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r.Headers, nil, http.StatusOK)
}

func (u *UCService) SearchFriends(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	req := &t.SearchFriendsReq{}
	api.Parse(r, req)
	friends, err := u.usr.SearchFriends(ctx, req.Search, req.IDs, req.Limit, req.Offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r.Headers, friends, http.StatusOK)
}
