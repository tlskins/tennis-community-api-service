package types

import (
	"github.com/tennis-community-api-service/pkg/enums"
	"time"
)

type SwingVideo struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	UserID          string                 `bson:"userId" json:"userId"`
	UploadKey       string                 `bson:"upKey" json:"uploadKey"`
	Clip            int                    `bson:"clip" json:"clip"`
	Swing           int                    `bson:"swing" json:"swing"`
	VideoURL        string                 `bson:"vidUrl" json:"videoURL"`
	ContactImageURL string                 `bson:"ctcImgUrl" json:"contactImageURL"`
	Status          enums.SwingVideoStatus `bson:"status" json:"status"`
	Tags            []string               `bson:"tags,omitempty" json:"tags,omitempty"`
}

type UpdateSwingVideo struct {
	ID        string    `bson:"_id" json:"id"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Status *enums.SwingVideoStatus `bson:"status,omitempty" json:"status,omitempty"`
	Tags   *[]string               `bson:"tags,omitempty" json:"tags,omitempty"`
}
