package types

import (
	"github.com/tennis-community-api-service/pkg/enums"
	"time"
)

type SwingUpload struct {
	ID          string                  `bson:"_id" json:"id"`
	CreatedAt   time.Time               `bson:"crAt" json:"createdAt"`
	UpdatedAt   time.Time               `bson:"updAt" json:"updatedAt"`
	UploadKey   string                  `bson:"upKey" json:"uploadKey"`
	UserID      string                  `bson:"usrId" json:"userId"`
	Status      enums.SwingUploadStatus `bson:"status" json:"status"`
	OriginalURL string                  `bson:"origUrl" json:"originalURL"`
	AlbumID     string                  `bson:"albId,omitempty" json:"albumId,omitempty"`
	ClipVideos  []*UploadClipVideo      `bson:"clipVids" json:"clipVideos"`
	SwingVideos []*UploadSwingVideo     `bson:"swingVids" json:"swingVideos"`
}

func (a SwingUpload) IsFinal() bool {
	finalMap := make(map[int]bool)
	for _, vid := range a.SwingVideos {
		finalMap[vid.ClipID] = true
	}
	isFinal := true
	for _, clipFinal := range finalMap {
		if !clipFinal {
			isFinal = false
		}
	}
	return isFinal
}

type UpdateSwingUpload struct {
	UploadKey string    `bson:"-" json:"uploadKey"`
	UserID    string    `bson:"usrId" json:"userId"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Status      *enums.SwingUploadStatus `bson:"status,omitempty" json:"status,omitempty"`
	OriginalURL *string                  `bson:"origUrl,omitempty" json:"originalURL,omitempty"`
	AlbumID     *string                  `bson:"albId,omitempty" json:"albumId,omitempty"`
	ClipVideos  *[]*UploadClipVideo      `bson:"clipVids,omitempty" json:"clipVideos,omitempty"`
	SwingVideos *[]*UploadSwingVideo     `bson:"swingVids,omitempty" json:"swingVideos,omitempty"`
}

type UploadClipVideo struct {
	ID        int       `bson:"id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	ClipURL   string    `bson:"clipUrl" json:"clipURL"`
}

type UploadSwingVideo struct {
	ID            string    `bson:"id" json:"id"`
	CreatedAt     time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt     time.Time `bson:"updAt" json:"updatedAt"`
	ClipID        int       `bson:"clipId" json:"clipId"`
	SwingID       int       `bson:"swId" json:"swingId"`
	CutURL        string    `bson:"cutUrl" json:"cutURL"`
	TranscodedURL string    `bson:"tranUrl" json:"transcodedURL"`
}

type UpdateUploadSwingVideo struct {
	ID            string    `bson:"-" json:"id"`
	UpdatedAt     time.Time `bson:"updAt" json:"updatedAt"`
	ClipID        *int      `bson:"clipId,omitempty" json:"clipId,omitempty"`
	CutURL        *string   `bson:"cutUrl,omitempty" json:"cutURL,omitempty"`
	TranscodedURL *string   `bson:"tranUrl,omitempty" json:"transcodedURL,omitempty"`
}
