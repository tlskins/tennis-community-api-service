package types

import (
	"errors"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/enums"
)

type CreateAlbumReq aT.Album

func (r CreateAlbumReq) Validate() error {
	if r.UserID == "" {
		return errors.New("Missing user id")
	}
	if r.Name == "" {
		return errors.New("Missing album name")
	}
	if r.Status != enums.AlbumStatusCreated {
		return errors.New("Album status must be created")
	}
	if len(r.SwingVideos) == 0 {
		return errors.New("No swings in album")
	}
	return nil
}
