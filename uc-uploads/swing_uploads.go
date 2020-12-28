package uploads

import (
	"context"

	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-uploads/types"
)

func (u *UCService) GetSwingUploadURL(ctx context.Context, r *api.Request) (api.Response, error) {
	req := &t.GetSwingUploadURLReq{}
	api.Parse(r, req)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	return u.up.GetSwingUploadURL(ctx, claims.Subject, req.FileName)
}

func (u *UCService) CreateSwingUpload(ctx context.Context, r *api.Request) (api.Response, error) {
	return u.up.CreateSwingUpload(ctx, r)
}
