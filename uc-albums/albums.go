package ucalbums

import (
	"context"
	"errors"
	"net/http"
	"time"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-albums/types"
)

func (u *UCService) GetAlbums(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	req := t.SearchAlbumsReq{}
	if authorized, claims := auth.ClaimsFromContext(ctx); authorized {
		req.UserID = claims.Subject
	}
	api.Parse(r, &req)

	albumResp := t.AlbumsResp{
		LastRequestAt: time.Now(),
		MyAlbums:      []*aT.Album{},
		FriendsAlbums: []*aT.Album{},
		PublicAlbums:  []*aT.Album{},
	}
	if req.UserID != "" {
		albumResp.MyAlbums, err = u.alb.GetUserAlbums(ctx, req.UserID)
		api.CheckError(http.StatusInternalServerError, err)
		if !req.ExcludeFriends {
			albumResp.FriendsAlbums, err = u.alb.GetFriendsAlbums(ctx, req.UserID)
			api.CheckError(http.StatusInternalServerError, err)
		}
	}
	if !req.ExcludePublic {
		albumResp.PublicAlbums, err = u.alb.GetPublicAlbums(ctx)
		api.CheckError(http.StatusInternalServerError, err)
	}

	return u.Resp.Success(r, albumResp, http.StatusOK)
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
