package types

import (
	"errors"
	"time"

	"github.com/tennis-community-api-service/pkg/enums"
)

type Album struct {
	ID            string            `bson:"_id" json:"id"`
	UserID        string            `bson:"userId" json:"userId"`
	UploadKey     string            `bson:"upKey" json:"uploadKey"`
	Name          string            `bson:"nm" json:"name"`
	Clips         int               `bson:"clips" json:"clips"`
	CreatedAt     time.Time         `bson:"crAt" json:"createdAt"`
	UpdatedAt     time.Time         `bson:"updAt" json:"updatedAt"`
	Status        enums.AlbumStatus `bson:"status" json:"status"`
	Tags          []string          `bson:"tags" json:"tags"`
	SwingVideoIDs []*SwingVideo     `bson:"swingVids" json:"swingVideos"`
}

type UpdateAlbum struct {
	ID        string    `bson:"-" json:"id"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Name        *string            `bson:"nm,omitempty" json:"name,omitempty"`
	Clips       *int               `bson:"clips,omitempty" json:"clips,omitempty"`
	Status      *enums.AlbumStatus `bson:"status,omitempty" json:"status,omitempty"`
	Tags        *[]string          `bson:"tags,omitempty" json:"tags,omitempty"`
	SwingVideos *[]*SwingVideo     `bson:"swingVids,omitempty" json:"swingVideos,omitempty"`
}

func (a UpdateAlbum) Validate() error {
	if a.Name != nil && *a.Name == "" {
		return errors.New("Name cannot be blank")
	}
	return nil
}
