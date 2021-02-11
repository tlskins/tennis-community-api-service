package uploads

import (
	"context"
	"strings"
	"time"

	"github.com/tennis-community-api-service/pkg/enums"
	t "github.com/tennis-community-api-service/uploads/types"
)

func (u *UploadsService) GetRecentSwingUploads(ctx context.Context, userId string) (uploads []*t.SwingUpload, err error) {
	return u.Store.GetRecentSwingUploads(userId)
}

func (u *UploadsService) CreateSwingUpload(ctx context.Context, userId, originalURL, albumName string, friendIds []string, isPublic, isFriends bool) (resp *t.SwingUpload, err error) {
	now := time.Now()
	paths := strings.Split(originalURL, "/")
	return u.Store.CreateSwingUpload(&t.SwingUpload{
		CreatedAt:           now,
		UpdatedAt:           now,
		UploadKey:           paths[len(paths)-2],
		UserID:              userId,
		Status:              enums.SwingUploadStatusOriginal,
		OriginalURL:         originalURL,
		AlbumName:           albumName,
		IsPublic:            isPublic,
		IsViewableByFriends: isFriends,
		FriendIDs:           friendIds,
	})
}

func (u *UploadsService) UpdateSwingUpload(_ context.Context, data *t.UpdateSwingUpload) (upload *t.SwingUpload, err error) {
	return u.Store.UpdateSwingUpload(data)
}
