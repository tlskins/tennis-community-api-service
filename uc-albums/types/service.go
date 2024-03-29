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

type UpdateAlbumReq struct {
	*aT.UpdateAlbum
	ShareAlbum bool `json:"shareAlbum"`
}

type SearchAlbumsReq struct {
	UserID       string  `json:"userId"`
	HomeApproved *string `json:"homeApproved"`
}

type AlbumsResp struct {
	LastRequestAt time.Time   `json:"lastRequestAt"`
	MyAlbums      []*aT.Album `json:"myAlbums"`
	FriendsAlbums []*aT.Album `json:"friendsAlbums"`
	PublicAlbums  []*aT.Album `json:"publicAlbums"`
}

type PostCommentReq struct {
	UserID  string `json:"userId"`
	AlbumID string `json:"albumId"`
	SwingID string `json:"swingId,omitempty"`
	ReplyID string `json:"replyId,omitempty"`
	Frame   int    `json:"frame,omitempty"`
	Text    string `json:"text"`
}

func (r PostCommentReq) Validate() error {
	if r.UserID == "" {
		return errors.New("Missing user")
	}
	if r.AlbumID == "" {
		return errors.New("Missing album")
	}
	if r.Text == "" {
		return errors.New("Missing comment text")
	}
	if len(r.Text) > 500 {
		return errors.New("Comment must be less than 500 characters")
	}
	return nil
}

type RecentAlbumsReq struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Limit  string    `json:"limit"`
	Offset string    `json:"offset"`
}

type RecentAlbumCommentsReq struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Limit  string    `json:"limit"`
	Offset string    `json:"offset"`
}

type RecentSwingCommentsReq struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Limit  string    `json:"limit"`
	Offset string    `json:"offset"`
}
