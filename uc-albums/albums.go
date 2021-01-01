package uploads

import (
	"context"
	"net/http"

	// "github.com/davecgh/go-spew/spew"

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
