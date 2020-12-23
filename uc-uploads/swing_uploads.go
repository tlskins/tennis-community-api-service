package uploads

import (
	"context"

	t "github.com/tennis-community-api-service/uploads/types"
)

func (u *UCService) CreateSwingUpload(ctx context.Context, origURL, userID string) (*t.SwingUpload, error) {
	return u.up.CreateSwingUpload(ctx, origURL, userID)
}

func (u *UCService) CreateUploadClipVideos(ctx context.Context, bucket string, outputs []string) (*t.SwingUpload, error) {
	return u.up.CreateUploadClipVideos(ctx, bucket, outputs)
}

func (u *UCService) CreateUploadSwingVideos(ctx context.Context, bucket string, outputs []string) (*t.SwingUpload, error) {
	return u.up.CreateUploadSwingVideos(ctx, bucket, outputs)
}
