package uploads

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"

	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-uploads/types"
)

func (u *UCService) GetRecentSwingUploads(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	uploads, err := u.up.GetRecentSwingUploads(ctx, claims.Subject)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(r.Headers, uploads, http.StatusCreated)
}

func (u *UCService) CreateSwingUpload(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	req := &t.CreateSwingUploadReq{}
	api.ParseAndValidate(r, req)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	spew.Dump(r)
	fmt.Printf("OriginalURL=%s\n", req.OriginalURL)
	upload, err := u.up.CreateSwingUpload(
		ctx,
		claims.Subject,
		req.OriginalURL,
		req.AlbumName,
		req.FriendIDs,
		req.IsPublic,
		req.IsViewableByFriends,
	)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(r.Headers, upload, http.StatusCreated)
}
