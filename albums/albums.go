package albums

import (
	"context"
	"time"

	t "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/enums"
)

func (u *AlbumsService) GetUserAlbums(ctx context.Context, userID string) ([]*t.Album, error) {
	return u.Store.GetAlbumsByUser(userID)
}

func (u *AlbumsService) GetAlbum(ctx context.Context, id string) (*t.Album, error) {
	return u.Store.GetAlbum(id)
}

func (u *AlbumsService) CreateAlbum(ctx context.Context, userID, uploadKey string, clips int) (*t.Album, error) {
	now := time.Now()
	return u.Store.CreateAlbum(&t.Album{
		Name:      uploadKey,
		UploadKey: uploadKey,
		UserID:    userID,
		Clips:     clips,
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusProcessing,
	})
}

func (u *AlbumsService) UpdateAlbum(ctx context.Context, data *t.UpdateAlbum) (*t.Album, error) {
	data.UpdatedAt = time.Now()
	return u.Store.UpdateAlbum(data)
}

func (u *AlbumsService) AddVideosToAlbum(ctx context.Context, userID, uploadKey string, swingVideos []*t.SwingVideo) (*t.Album, error) {
	return u.Store.AddVideosToAlbum(userID, uploadKey, swingVideos)
}
