package types

import (
	"time"
)

type Comment struct {
	ID        string     `bson:"_id" json:"id"`
	ReplyID   string     `bson:"replyId" json:"replyId"`
	UserID    string     `bson:"userId" json:"userId"`
	CreatedAt time.Time  `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updAt" json:"updatedAt"`
	Frame     int        `bson:"frame,omitempty" json:"frame,omitempty"`
	Text      string     `bson:"txt" json:"text"`
	UserTags  []*UserTag `bson:"usrTags" json:"userTags"`

	// aggregated fields
	AlbumID string `bson:"albumId" json:"albumId,omitempty"`
	SwingID string `bson:"swingId" json:"swingId,omitempty"`
}

type UserTag struct {
	UserID    string `bson:"usrId" json:"userId"`
	UserName  string `bson:"usrNm" json:"userName"`
	FirstName string `bson:"fNm" json:"firstName"`
	LastName  string `bson:"lNm" json:"lastName"`
	Start     int    `bson:"st" json:"start"`
	End       int    `bson:"end" json:"end"`
}
