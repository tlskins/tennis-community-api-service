package types

import (
	"github.com/tennis-community-api-service/pkg/enums"
	"time"
)

type Album struct {
	ID          string            `bson:"_id" json:"id"`
	UserID      string            `bson:"userId" json:"userId"`
	UploadKey   string            `bson:"upKey" json:"uploadKey"`
	Name        string            `bson:"nm" json:"name"`
	Clips       int               `bson:"clips" json:"clips"`
	CreatedAt   time.Time         `bson:"crAt" json:"createdAt"`
	UpdatedAt   time.Time         `bson:"updAt" json:"updatedAt"`
	Status      enums.AlbumStatus `bson:"status" json:"status"`
	Tags        []string          `bson:"tags" json:"tags"`
	SwingVideos []*SwingVideo     `bson:"swingVids" json:"swingVideos"`
}

func (a Album) IsFinal() bool {
	finalMap := make(map[int]bool)
	for _, vid := range a.SwingVideos {
		finalMap[vid.Clip] = true
	}
	isFinal := true
	for _, clipFinal := range finalMap {
		if !clipFinal {
			isFinal = false
		}
	}
	return isFinal
}

type UpdateAlbum struct {
	ID        string    `bson:"-" json:"id"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	UserID      *string            `bson:"userId,omitempty" json:"userId,omitempty"`
	UploadKey   *string            `bson:"upKey,omitempty" json:"uploadKey,omitempty"`
	Name        *string            `bson:"nm,omitempty" json:"name,omitempty"`
	Clips       *int               `bson:"clips,omitempty" json:"clips,omitempty"`
	Status      *enums.AlbumStatus `bson:"status,omitempty" json:"status,omitempty"`
	Tags        *[]string          `bson:"tags,omitempty" json:"tags,omitempty"`
	SwingVideos *[]*SwingVideo     `bson:"swingVids,omitempty" json:"swingVideos,omitempty"`
}

type SwingVideo struct {
	ID              int                    `bson:"id" json:"id"`
	CreatedAt       time.Time              `bson:"crAt" json:"createdAt"`
	Clip            int                    `bson:"clip" json:"clip"`
	Swing           int                    `bson:"swing" json:"swing"`
	VideoURL        string                 `bson:"vidUrl" json:"videoURL"`
	ContactImageURL string                 `bson:"ctcImgUrl" json:"contactImageURL"`
	Status          enums.SwingVideoStatus `bson:"status" json:"status"`
	Tags            []string               `bson:"tags,omitempty" json:"tags,omitempty"`
}
