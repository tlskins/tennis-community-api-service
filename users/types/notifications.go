package types

import (
	"time"
)

type UploadNote struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Subject   string    `bson:"subj" json:"subject"`
	Body      string    `bson:"body" json:"body"`
	Type      string    `bson:"tp" json:"type"`
	UploadID  string    `bson:"upId" json:"uploadId"`
	UploadKey string    `bson:"upKey" json:"firstName"`
	AlbumID   string    `bson:"albId" json:"albumId"`
	UploadAt  time.Time `bson:"upAt" json:"lastName"`
}
