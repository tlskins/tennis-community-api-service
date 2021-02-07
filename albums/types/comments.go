package types

import (
	"time"
)

type Comment struct {
	ID        string    `bson:"_id" json:"id"`
	ReplyID   string    `bson:"replyId" json:"replyId"`
	UserID    string    `bson:"userId" json:"userId"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`
	Frame     int       `bson:"frame,omitempty" json:"frame,omitempty"`
	Text      string    `bson:"txt" json:"text"`

	// aggregated fields
	AlbumID string `bson:"albumId" json:"albumId,omitempty"`
	SwingID string `bson:"swingId" json:"swingId,omitempty"`
}
