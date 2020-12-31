package uploads

import (
	"context"
	"net/http"
	"strings"
	"time"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/auth"
	"github.com/tennis-community-api-service/pkg/enums"
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
	req := &t.CreateSwingUploadReq{}
	api.ParseAndValidate(r, req)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	return u.up.CreateSwingUpload(ctx, claims.Subject, req.OriginalURL)
}

func (u *UCService) CreateUploadClipVideos(ctx context.Context, r *t.SwingStorageEvent) (string, error) {
	upload, err := u.up.CreateUploadClipVideos(ctx, r.ResponsePayload.Body.Bucket, r.ResponsePayload.Body.Outputs)
	if err != nil {
		return "error", err
	}
	_, err = u.alb.CreateAlbum(ctx, upload.UserID, upload.UploadKey, len(upload.ClipVideos))
	if err != nil {
		return "error", err
	}
	return "success", nil
}

func (u *UCService) CreateUploadSwingVideos(ctx context.Context, r *t.SwingStorageEvent) (string, error) {
	upload, swings, err := u.up.CreateUploadSwingVideos(ctx, r.ResponsePayload.Body.Bucket, r.ResponsePayload.Body.Outputs)
	if err != nil {
		return "error", err
	}

	now := time.Now()
	swingVids := make([]*aT.SwingVideo, len(swings))
	for i, swing := range swings {
		swingVids[i] = &aT.SwingVideo{
			ID:        i,
			CreatedAt: now,
			Clip:      swing.ClipID,
			Swing:     swing.SwingID,
			VideoURL:  strings.Replace(swing.CutURL, "tmp/", "", 1),
			Status:    enums.SwingVideoStatusCreated,
		}
	}
	album, err := u.alb.AddVideosToAlbum(ctx, upload.UserID, upload.UploadKey, swingVids)
	if err != nil {
		return "error", err
	}

	if album.IsFinal() {
		status := enums.AlbumStatusCreated
		_, err = u.alb.UpdateAlbum(ctx, &aT.UpdateAlbum{
			ID:     album.ID,
			Status: &status,
		})
	}
	return "success", nil
}
