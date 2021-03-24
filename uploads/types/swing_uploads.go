package types

import (
	"time"

	"github.com/tennis-community-api-service/pkg/enums"
)

type SwingUpload struct {
	ID                  string                  `bson:"_id" json:"id"`
	CreatedAt           time.Time               `bson:"crAt" json:"createdAt"`
	UpdatedAt           time.Time               `bson:"updAt" json:"updatedAt"`
	UploadKey           string                  `bson:"upKey" json:"uploadKey"`
	UserID              string                  `bson:"usrId" json:"userId"`
	Status              enums.SwingUploadStatus `bson:"status" json:"status"`
	OriginalURL         string                  `bson:"origUrl" json:"originalURL"`
	AlbumID             string                  `bson:"albId,omitempty" json:"albumId,omitempty"`
	AlbumName           string                  `bson:"albNm" json:"albumName"`
	IsPublic            bool                    `bson:"public" json:"isPublic"`
	IsViewableByFriends bool                    `bson:"frndView" json:"isViewableByFriends"`
	FriendIDs           []string                `bson:"frndIds" json:"friendIds"`
	ClipVideos          []*UploadClipVideo      `bson:"clipVids" json:"clipVideos"`
	SwingVideos         []*UploadSwingVideo     `bson:"swingVids" json:"swingVideos"`
	ProcessedClips      []int                   `bson:"procClips" json:"processedClips"`
}

func (a SwingUpload) IsFinal() bool {
	if len(a.ClipVideos) == 0 {
		// clips havent finished processing yet
		return false
	}
	return len(a.ProcessedClips) == len(a.ClipVideos)
}

type UpdateSwingUpload struct {
	UploadKey string    `bson:"-" json:"uploadKey"`
	UserID    string    `bson:"usrId" json:"userId"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Status              *enums.SwingUploadStatus `bson:"status,omitempty" json:"status,omitempty"`
	OriginalURL         *string                  `bson:"origUrl,omitempty" json:"originalURL,omitempty"`
	AlbumID             *string                  `bson:"albId,omitempty" json:"albumId,omitempty"`
	AlbumName           *string                  `bson:"albNm,omitempty" json:"albumName,omitempty"`
	IsPublic            *bool                    `bson:"public,omitempty" json:"isPublic,omitempty"`
	IsViewableByFriends *bool                    `bson:"frndView,omitempty" json:"isViewableByFriends,omitempty"`
	FriendIDs           *[]string                `bson:"frndIds,omitempty" json:"friendIds,omitempty"`
	ClipVideos          *[]*UploadClipVideo      `bson:"clipVids,omitempty" json:"clipVideos,omitempty"`
	SwingVideos         *[]*UploadSwingVideo     `bson:"swingVids,omitempty" json:"swingVideos,omitempty"`
}

type UploadClipVideo struct {
	ID           int       `bson:"id" json:"id"`
	CreatedAt    time.Time `bson:"crAt" json:"createdAt"`
	ClipURL      string    `bson:"clipUrl" json:"clipURL"`
	FileName     string    `bson:"fileNm" json:"fileName"`
	StartSeconds int       `bson:"st" json:"startSec"`
	EndSeconds   int       `bson:"end" json:"endSec"`
}

type UploadSwingVideo struct {
	ID               string    `bson:"id" json:"id"`
	CreatedAt        time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt        time.Time `bson:"updAt" json:"updatedAt"`
	TimestampSeconds int       `bson:"timeSt" json:"timestampSecs"`
	Frames           int       `bson:"frames" json:"frames"`
	ClipID           int       `bson:"clipId" json:"clipId"`
	SwingID          int       `bson:"swId" json:"swingId"`
	CutURL           string    `bson:"cutUrl" json:"cutURL"`
	GifURL           string    `bson:"gifUrl" json:"gifURL"`
	JpgURL           string    `bson:"jpgUrl" json:"jpgURL"`
}

type UpdateUploadSwingVideo struct {
	ID        string    `bson:"-" json:"id"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`
	ClipID    *int      `bson:"clipId,omitempty" json:"clipId,omitempty"`
	CutURL    *string   `bson:"cutUrl,omitempty" json:"cutURL,omitempty"`
}
