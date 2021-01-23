package types

import (
	"github.com/tennis-community-api-service/pkg/enums"
	"time"
)

type SwingVideo struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	UserID    string                 `bson:"userId" json:"userId"`
	UploadKey string                 `bson:"upKey" json:"uploadKey"`
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
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Status   *enums.SwingVideoStatus `bson:"status,omitempty" json:"status,omitempty"`
	Tags     *[]string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Comments *[]*Comment             `bson:"cmnts,omitempty" json:"comments,omitempty"`
}
