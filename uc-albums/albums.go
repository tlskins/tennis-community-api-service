package ucalbums

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-albums/types"
)

func (u *UCService) SearchAlbums(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	req := t.SearchAlbumsReq{}
	_, claims := auth.ClaimsFromContext(ctx)
	api.Parse(r, &req)

	var userID, friendID string
	if req.My {
		userID = claims.Subject
	}
	if req.Friends {
		friendID = claims.Subject
	}
	albums, err := u.alb.SearchAlbums(ctx, userID, friendID, req.Public, req.HomeApproved, req.Limit, req.Offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, albums, http.StatusOK)
}

func (u *UCService) CreateAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := t.CreateAlbumReq{UserID: claims.Subject}
	api.ParseAndValidate(r, &req)
	albumReq := aT.Album(req)
	album, err := u.alb.CreateAlbum(ctx, &albumReq)
	api.CheckError(http.StatusInternalServerError, err)
	err = u.shareAlbum(ctx, album)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(r, album, http.StatusOK)
}

func (u *UCService) GetAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	id := r.PathParameters["id"]
	album, err := u.alb.GetAlbum(ctx, id)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(r, album, http.StatusOK)
}

func (u *UCService) DeleteAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	id := r.PathParameters["id"]
	album, err := u.alb.GetAlbum(ctx, id)
	if album.UserID != claims.Subject {
		api.CheckError(http.StatusUnauthorized, errors.New("Unable to delete an album that does not belong to you"))
	}
	u.alb.DeleteAlbum(ctx, id)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(r, nil, http.StatusOK)
}

func (u *UCService) UpdateAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := &t.UpdateAlbumReq{}
	api.ParseAndValidate(r, req)
	req.UpdateAlbum.ID = r.PathParameters["id"]
	album, err := u.alb.GetAlbum(ctx, req.ID)
	api.CheckError(http.StatusNotFound, err)
	if album.UserID != claims.Subject {
		panic(errors.New("Cannot edit another user's album"))
	}
	album, err = u.alb.UpdateAlbum(ctx, req.UpdateAlbum)
	api.CheckError(http.StatusUnprocessableEntity, err)
	if req.ShareAlbum {
		err = u.shareAlbum(ctx, album)
		api.CheckError(http.StatusUnprocessableEntity, err)
	}
	return u.Resp.Success(r, album, http.StatusOK)
}

func (u *UCService) RecentAlbums(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	if !claims.IsAdmin {
		api.CheckError(http.StatusUnauthorized, errors.New("Must be admin"))
	}
	req := &t.RecentAlbumsReq{}
	api.Parse(r, req)
	limit, err := strconv.Atoi(req.Limit)
	api.CheckError(http.StatusUnprocessableEntity, err)
	offset, err := strconv.Atoi(req.Offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	albums, err := u.alb.RecentAlbums(ctx, req.Start, req.End, limit, offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, albums, http.StatusOK)
}

func (u *UCService) RecentAlbumComments(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	if !claims.IsAdmin {
		api.CheckError(http.StatusUnauthorized, errors.New("Must be admin"))
	}
	req := &t.RecentAlbumCommentsReq{}
	api.Parse(r, req)
	limit, err := strconv.Atoi(req.Limit)
	api.CheckError(http.StatusUnprocessableEntity, err)
	offset, err := strconv.Atoi(req.Offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	comments, err := u.alb.RecentAlbumComments(ctx, req.Start, req.End, limit, offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, comments, http.StatusOK)
}

func (u *UCService) RecentSwingComments(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	if !claims.IsAdmin {
		api.CheckError(http.StatusUnauthorized, errors.New("Must be admin"))
	}
	req := &t.RecentSwingCommentsReq{}
	api.Parse(r, req)
	limit, err := strconv.Atoi(req.Limit)
	api.CheckError(http.StatusUnprocessableEntity, err)
	offset, err := strconv.Atoi(req.Offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	comments, err := u.alb.RecentSwingComments(ctx, req.Start, req.End, limit, offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, comments, http.StatusOK)
}
