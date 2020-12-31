package types

import (
	"github.com/tennis-community-api-service/pkg/enums"
	"time"
)

type Album struct {
	ID          string            `bson:"_id" json:"id"`
	UserID      string            `bson:"userId" json:"userId"`
	Name        string            `bson:"nm" json:"name"`
	CreatedAt   time.Time         `bson:"crAt" json:"createdAt"`
	UpdatedAt   time.Time         `bson:"updAt" json:"updatedAt"`
	Status      enums.AlbumStatus `bson:"status" json:"status"`
	Tags        []string          `bson:"tags" json:"tags"`
	SwingVideos []*SwingVideo     `bson:"swingVids" json:"swingVideos"`
}

type SwingVideo struct {
	ID              int                    `bson:"id" json:"id"`
	CreatedAt       time.Time              `bson:"crAt" json:"createdAt"`
	VideoURL        string                 `bson:"vidUrl" json:"videoURL"`
	ContactImageURL string                 `bson:"ctcImgUrl" json:"contactImageURL"`
	Status          enums.SwingVideoStatus `bson:"status" json:"status"`
	Tags            []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Order           int                    `bson:"order,omitempty" json:"order,omitempty"`
}
