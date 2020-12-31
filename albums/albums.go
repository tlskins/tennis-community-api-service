package albums

import (
	"context"
	"net/http"
	"time"

	t "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/enums"
	api "github.com/tennis-community-api-service/pkg/lambda"
)

func (u *AlbumsService) CreateAlbum(ctx context.Context, userID, uploadKey string, swingVideoUrls []string) (resp api.Response, err error) {
	now := time.Now()
	swingVids := make([]*t.SwingVideo, len(swingVideoUrls))
	for i, url := range swingVideoUrls {
		swingVids[i] = &t.SwingVideo{
			ID:        i,
			CreatedAt: now,
			VideoURL:  url,
			Status:    enums.SwingVideoStatusCreated,
			Order:     i,
		}
	}

	album := &t.Album{
		Name:        uploadKey,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
		Status:      enums.AlbumStatusCreated,
		SwingVideos: swingVids,
	}
	newAlbum, err := u.Store.CreateAlbum(album)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return api.Success(newAlbum, http.StatusCreated)
}
