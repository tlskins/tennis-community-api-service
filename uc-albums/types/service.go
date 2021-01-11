package types

import (
	"errors"
	"time"

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

type SearchAlbumsReq struct {
	UserID         string `json:"userId"`
	ExcludeFriends bool   `json:"excludeFriends"`
	ExcludePublic  bool   `json:"excludePublic"`
}

type AlbumsResp struct {
	LastRequestAt time.Time   `json:"lastRequestAt"`
	MyAlbums      []*aT.Album `json:"myAlbums"`
	FriendsAlbums []*aT.Album `json:"friendsAlbums"`
	PublicAlbums  []*aT.Album `json:"publicAlbums"`
}
