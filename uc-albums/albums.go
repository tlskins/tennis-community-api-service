package uploads

import (
	"context"
	"errors"
	"net/http"

	// "github.com/davecgh/go-spew/spew"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
)

func (u *UCService) GetUserAlbums(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	albums, err := u.alb.GetUserAlbums(ctx, claims.Subject)
	api.CheckError(http.StatusInternalServerError, err)
	return api.Success(albums, http.StatusOK)
}

func (u *UCService) GetAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	id := r.PathParameters["id"]
	album, err := u.alb.GetAlbum(ctx, id)
	api.CheckError(http.StatusInternalServerError, err)
	return api.Success(album, http.StatusOK)
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
	return api.Success(album, http.StatusOK)
}
