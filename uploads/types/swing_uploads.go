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
	ClipVideos  []*UploadClipVideo      `bson:"clipVids" json:"clipVideos"`
	SwingVideos []*UploadSwingVideo     `bson:"swingVids" json:"swingVideos"`
}

func (u SwingUpload) SwingClips() (out map[int][]string, finished bool) {
	out = make(map[int][]string)
	finished = true
	for _, clip := range u.ClipVideos {
		out[clip.ID] = []string{}
		for _, swing := range u.SwingVideos {
			if swing.ClipID == clip.ID {
				out[clip.ID] = append(out[clip.ID], swing.CutURL)
			}
		}
		if len(out[clip.ID]) == 0 {
			finished = false
		}
	}
	return
}

type UpdateSwingUpload struct {
	UploadKey   string                   `bson:"-" json:"uploadKey"`
	UpdatedAt   time.Time                `bson:"updAt,omitempty" json:"updatedAt,omitempty"`
	UserID      *string                  `bson:"usrId,omitempty" json:"userId,omitempty"`
	Status      *enums.SwingUploadStatus `bson:"status,omitempty" json:"status,omitempty"`
	OriginalURL *string                  `bson:"origUrl,omitempty" json:"originalURL,omitempty"`
	ClipVideos  *[]*UploadClipVideo      `bson:"clipVids,omitempty" json:"clipVideos,omitempty"`
	SwingVideos *[]*UploadSwingVideo     `bson:"swingVids,omitempty" json:"swingVideos,omitempty"`
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
	SwingID       int       `bson:"swId" json:"swingId"`
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
