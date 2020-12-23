package types

import (
	"github.com/tennis-community-api-service/pkg/enums"
	"time"
)

type SwingUpload struct {
	ID          string                  `bson:"_id" json:"id"`
	CreatedAt   time.Time               `bson:"crAt" json:"createdAt"`
	UpdatedAt   time.Time               `bson:"updAt" json:"updatedAt"`
	UserID      string                  `bson:"usrId" json:"userId"`
	Status      enums.SwingUploadStatus `bson:"status" json:"status"`
	OriginalURL string                  `bson:"origUrl" json:"originalURL"`
	ClipVideos  []*UploadClipVideo      `bson:"clipVids" json:"clipVideos"`
	SwingVideos []*UploadSwingVideo     `bson:"swingVids" json:"swingVideos"`
}

type UploadClipVideo struct {
	ID        int       `bson:"id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	ClipURL   string    `bson:"clipUrl" json:"clipURL"`
}

type UploadSwingVideo struct {
	ID            int       `bson:"id" json:"id"`
	CreatedAt     time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt     time.Time `bson:"updAt" json:"updatedAt"`
	ClipID        int       `bson:"clipId" json:"clipId"`
	CutURL        string    `bson:"cutUrl" json:"cutURL"`
	TranscodedURL string    `bson:"tranUrl" json:"transcodedURL"`
}

type UpdateUploadSwingVideo struct {
	ID            int       `bson:"-" json:"id"`
	UpdatedAt     time.Time `bson:"updAt" json:"updatedAt"`
	ClipID        *int      `bson:"clipId,omitempty" json:"clipId,omitempty"`
	CutURL        *string   `bson:"cutUrl,omitempty" json:"cutURL,omitempty"`
	TranscodedURL *string   `bson:"tranUrl,omitempty" json:"transcodedURL,omitempty"`
}
