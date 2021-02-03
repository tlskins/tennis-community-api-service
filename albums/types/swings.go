package types

import (
	"time"

	"github.com/tennis-community-api-service/pkg/enums"
)

type SwingVideo struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	UserID    string                 `bson:"userId" json:"userId"`
	UploadKey string                 `bson:"upKey" json:"uploadKey"`
	Name      string                 `bson:"nm" json:"name"`
	Clip      int                    `bson:"clip" json:"clip"`
	Swing     int                    `bson:"swing" json:"swing"`
	VideoURL  string                 `bson:"vidUrl" json:"videoURL"`
	GifURL    string                 `bson:"gifUrl" json:"gifURL"`
	JpgURL    string                 `bson:"jpgUrl" json:"jpgURL"`
	Status    enums.SwingVideoStatus `bson:"status" json:"status"`
	Tags      []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Comments  []*Comment             `bson:"cmnts" json:"comments"`
}

type UpdateSwingVideo struct {
	ID        string    `bson:"_id" json:"id"`
	AlbumID   string    `bson:"-" json:"albumId"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	UserID    *string                 `bson:"userId,omitempty" json:"userId,omitempty"`
	UploadKey *string                 `bson:"upKey,omitempty" json:"uploadKey,omitempty"`
	Name      *string                 `bson:"nm,omitempty" json:"name,omitempty"`
	Clip      *int                    `bson:"clip,omitempty" json:"clip,omitempty"`
	Swing     *int                    `bson:"swing,omitempty" json:"swing,omitempty"`
	VideoURL  *string                 `bson:"vidUrl,omitempty" json:"videoURL,omitempty"`
	GifURL    *string                 `bson:"gifUrl,omitempty" json:"gifURL,omitempty"`
	JpgURL    *string                 `bson:"jpgUrl,omitempty" json:"jpgURL,omitempty"`
	Status    *enums.SwingVideoStatus `bson:"status,omitempty" json:"status,omitempty"`
	Tags      *[]string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Comments  *[]*Comment             `bson:"cmnts,omitempty" json:"comments,omitempty"`
}
