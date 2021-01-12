package uploads

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	// "github.com/davecgh/go-spew/spew"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-albums/types"
	uT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) GetAlbums(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := t.SearchAlbumsReq{UserID: claims.Subject}
	api.Parse(r, &req)

	albumResp := t.AlbumsResp{
		LastRequestAt: time.Now(),
		MyAlbums:      []*aT.Album{},
		FriendsAlbums: []*aT.Album{},
		PublicAlbums:  []*aT.Album{},
	}
	albumResp.MyAlbums, err = u.alb.GetUserAlbums(ctx, claims.Subject)
	api.CheckError(http.StatusInternalServerError, err)
	if !req.ExcludePublic {
		albumResp.PublicAlbums, err = u.alb.GetPublicAlbums(ctx)
		api.CheckError(http.StatusInternalServerError, err)
	}
	if !req.ExcludeFriends {
		albumResp.FriendsAlbums, err = u.alb.GetFriendsAlbums(ctx, claims.Subject)
		api.CheckError(http.StatusInternalServerError, err)
	}
	return u.Resp.Success(albumResp, http.StatusOK)
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

	user, err := u.usr.GetUser(ctx, claims.Subject)
	api.CheckError(http.StatusInternalServerError, err)
	if len(album.FriendIDs) > 0 {
		err = u.usr.AddFriendNoteToUsers(ctx, album.FriendIDs, &uT.FriendNote{
			CreatedAt: time.Now(),
			Subject:   fmt.Sprintf("%s has shared the album %s with you!", user.UserName, album.Name),
		})
		api.CheckError(http.StatusInternalServerError, err)
	}
	return u.Resp.Success(album, http.StatusOK)
}

func (u *UCService) GetAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	id := r.PathParameters["id"]
	album, err := u.alb.GetAlbum(ctx, id)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(album, http.StatusOK)
}

func (u *UCService) UpdateAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := &aT.UpdateAlbum{ID: r.PathParameters["id"]}
	api.ParseAndValidate(r, req)
	album, err := u.alb.GetAlbum(ctx, req.ID)
	api.CheckError(http.StatusNotFound, err)
	if album.UserID != claims.Subject {
		panic(errors.New("Cannot edit another user's album"))
	}
	album, err = u.alb.UpdateAlbum(ctx, req)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(album, http.StatusOK)
}
