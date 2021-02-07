package types

import (
	"time"
)

type CommentFlag struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"upAt" json:"updatedAt"`

	CommentCreatedAt time.Time `bson:"comCrAt" json:"commentCreatedAt"`
	CommentID        string    `bson:"commId" json:"commentId"`
	CommenterID      string    `bson:"commerId" json:"commenterId"`
	FlaggerID        string    `bson:"flaggerId" json:"flaggerId"`
	AlbumID          string    `bson:"albId" json:"albumId"`
	SwingID          string    `bson:"swId,omitempty" json:"swingId,omitempty"`
	Text             string    `bson:"txt" json:"text"`
	Resolved         bool      `bson:"res" json:"resolved"`
}

type UpdateCommentFlag struct {
	ID        string    `bson:"_id" json:"id"`
	UpdatedAt time.Time `bson:"upAt" json:"updatedAt"`
	Resolved  bool      `bson:"res" json:"resolved"`
}

type AlbumFlag struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"upAt" json:"updatedAt"`

	AlbumCreatedAt time.Time `bson:"albCrAt" json:"albumCreatedAt"`
	AlbumUserID    string    `bson:"albUsrId" json:"albumUserId"`
	FlaggerID      string    `bson:"flaggerId" json:"flaggerId"`
	AlbumID        string    `bson:"albId" json:"albumId"`
	AlbumName      string    `bson:"albNm" json:"albumName"`
	Resolved       bool      `bson:"res" json:"resolved"`
}

type UpdateAlbumFlag struct {
	ID        string    `bson:"_id" json:"id"`
	UpdatedAt time.Time `bson:"upAt" json:"updatedAt"`
	Resolved  bool      `bson:"res" json:"resolved"`
}
