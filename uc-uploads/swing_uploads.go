package uploads

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	// "github.com/davecgh/go-spew/spew"

	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-uploads/types"
)

// func (u *UCService) GetSwingUploadURL(ctx context.Context, r *api.Request) (resp api.Response, err error) {
// 	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
// 	api.CheckError(http.StatusInternalServerError, err)
// 	req := &t.GetSwingUploadURLReq{}
// 	api.ParseAndValidate(r, req)
// 	claims := auth.AuthorizedClaimsFromContext(ctx)
// 	return u.up.GetSwingUploadURL(ctx, claims.Subject, req.FileName)
// }

func (u *UCService) CreateSwingUpload(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	fmt.Println("after include lambda")
	req := &t.CreateSwingUploadReq{}
	api.ParseAndValidate(r, req)
	fmt.Println("after parse and validate")
	claims := auth.AuthorizedClaimsFromContext(ctx)
	fmt.Println("after AuthorizedClaimsFromContext")
	return u.up.CreateSwingUpload(ctx, claims.Subject, req.OriginalURL)
}

func (u *UCService) CreateUploadClipVideos(ctx context.Context, r *t.SwingStorageEvent) (string, error) {
	_, err := u.up.CreateUploadClipVideos(ctx, r.ResponsePayload.Body.Bucket, r.ResponsePayload.Body.Outputs)
	if err != nil {
		return "error", err
	}
	return "success", nil
}

func (u *UCService) CreateUploadSwingVideos(ctx context.Context, r *t.SwingStorageEvent) (string, error) {
	upload, err := u.up.CreateUploadSwingVideos(ctx, r.ResponsePayload.Body.Bucket, r.ResponsePayload.Body.Outputs)
	if err != nil {
		return "error", err
	}
	swingClips, finished := upload.SwingClips()
	if finished {
		swingVids := []string{}
		for _, vids := range swingClips {
			for _, vid := range vids {
				swingVids = append(swingVids, strings.Replace(vid, "tmp/", "", 1))
			}
		}
		u.alb.CreateAlbum(ctx, upload.UserID, upload.UploadKey, swingVids)
	}
	return "success", nil
}
