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

func (u *UCService) GetSwingUploadURL(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	fmt.Println(r.QueryStringParameters)
	fmt.Println(r.QueryStringParameters["fileName"])
	req := &t.GetSwingUploadURLReq{}
	api.ParseAndValidate(r, req)
	spew.Dump(req)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	return u.up.GetSwingUploadURL(ctx, claims.Subject, req.FileName)
}

func (u *UCService) CreateSwingUpload(ctx context.Context, r *api.Request) (api.Response, error) {
	return u.up.CreateSwingUpload(ctx, r)
}
