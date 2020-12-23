package types

import (
	"github.com/tennis-community-api-service/pkg/enums"
	"time"
)

type Album struct {
	ID        string            `bson:"_id" json:"id"`
	Name      string            `bson:"nm" json:"name"`
	CreatedAt *time.Time        `bson:"crAt" json:"createdAt"`
	UpdatedAt *time.Time        `bson:"updAt" json:"updatedAt"`
	Status    enums.AlbumStatus `bson:"status" json:"status"`
	Date      *time.Time        `bson:"dt" json:"date"`
	Tags      []string          `bson:"tags" json:"tags"`
	// OriginalVideo string   `bson:"origVid" json:"-"`
	// ClipVideos    []string `bson:"clipVids" json:"-"`
	SwingVideos []string `bson:"swingVids" json:"swingVideos"`
}

type SwingVideo struct {
	ID              string                 `bson:"_id" json:"id"`
	CreatedAt       *time.Time             `bson:"crAt" json:"createdAt"`
	VideoURL        string                 `bson:"vidUrl" json:"videoURL"`
	ContactImageURL string                 `bson:"ctcImgUrl" json:"contactImageURL"`
	Status          enums.SwingVideoStatus `bson:"status" json:"status"`
	Tags            []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Order           int                    `bson:"order,omitempty" json:"order,omitempty"`
}
