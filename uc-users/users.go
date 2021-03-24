package users

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/tennis-community-api-service/pkg/auth"
	"github.com/tennis-community-api-service/pkg/enums"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-users/types"
	uT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) UpdateUserProfile(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := t.UpdateUserProfileReq{ID: claims.Subject}
	api.ParseAndValidate(r, &req)

	usrReq := uT.UpdateUserProfile(req)
	user, err := u.usr.UpdateUserProfile(ctx, &usrReq)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r.Headers, user, http.StatusOK)
}

func (u *UCService) SignIn(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	req := &t.SignInReq{}
	api.Parse(r, req)

	user, err := u.usr.GetUserByEmail(ctx, req.Email)
	api.CheckError(http.StatusNotFound, err)
	if user.Status == enums.UserStatusPending {
		api.Abort(http.StatusForbidden, "Please confirm email before signing in")
	}
	err = auth.ValidateCredentials(user.EncryptedPassword, req.Password)
	api.CheckError(http.StatusUnauthorized, err, "Incorrect Password")
	authToken, err := u.jwt.GenAccessToken(user)
	api.CheckError(http.StatusInternalServerError, err)
	now := time.Now()
	user, err = u.usr.UpdateUser(ctx, &uT.UpdateUser{
		ID:           user.ID,
		LastLoggedIn: &now,
	})
	api.CheckError(http.StatusUnprocessableEntity, err)
	user.AuthToken = authToken

	return u.Resp.Success(r.Headers, user, http.StatusOK)
}

func (u *UCService) GetUser(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	user, err := u.usr.GetUser(ctx, claims.Subject)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r.Headers, user, http.StatusOK)
}

func (u *UCService) RemoveUserNotification(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := &t.RemoveNotificationReq{}
	api.ParseAndValidate(r, req)

	var user *uT.User
	if len(req.UploadNoteID) > 0 {
		user, err = u.usr.RemoveUploadNote(ctx, claims.Subject, req.UploadNoteID)
		api.CheckError(http.StatusUnprocessableEntity, err)
	} else if len(req.FriendNoteID) > 0 {
		user, err = u.usr.RemoveFriendNote(ctx, claims.Subject, req.FriendNoteID)
		api.CheckError(http.StatusUnprocessableEntity, err)
	} else if len(req.CommentNoteID) > 0 {
		user, err = u.usr.RemoveCommentNote(ctx, claims.Subject, req.CommentNoteID)
		api.CheckError(http.StatusUnprocessableEntity, err)
	} else if len(req.AlbumUserTagNoteID) > 0 {
		user, err = u.usr.RemoveAlbumUserTagNote(ctx, claims.Subject, req.AlbumUserTagNoteID)
		api.CheckError(http.StatusUnprocessableEntity, err)
	}
	return u.Resp.Success(r.Headers, user, http.StatusOK)
}

func (u *UCService) RecentUsers(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	if !claims.IsAdmin {
		api.CheckError(http.StatusUnauthorized, errors.New("Must be admin"))
	}
	req := &t.RecentUsersReq{}
	api.Parse(r, req)
	limit, err := strconv.Atoi(req.Limit)
	api.CheckError(http.StatusUnprocessableEntity, err)
	offset, err := strconv.Atoi(req.Offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	users, err := u.usr.RecentUsers(ctx, req.Start, req.End, limit, offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r.Headers, users, http.StatusOK)
}
