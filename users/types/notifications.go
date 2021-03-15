package types

import (
	"time"
)

type UploadNote struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`

	Subject   string    `bson:"subj" json:"subject"`
	Body      string    `bson:"body" json:"body"`
	Type      string    `bson:"tp" json:"type"`
	UploadID  string    `bson:"upId" json:"uploadId"`
	UploadKey string    `bson:"upKey" json:"firstName"`
	AlbumID   string    `bson:"albId" json:"albumId"`
	UploadAt  time.Time `bson:"upAt" json:"lastName"`
}

type FriendNote struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`

	Subject  string `bson:"subj" json:"subject"`
	Body     string `bson:"body" json:"body"`
	Type     string `bson:"tp" json:"type"`
	FriendID string `bson:"frndId" json:"friendId"`
}

type CommentNote struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	FriendID        string   `bson:"frndId,omitempty" json:"friendId,omitempty"`
	FriendFirstName string   `bson:"frndFirst,omitempty" json:"friendFirstName,omitempty"`
	FriendUserName  string   `bson:"frndUsr,omitempty" json:"friendUserName,omitempty"`
	AlbumID         string   `bson:"albumId" json:"albumId"`
	AlbumName       string   `bson:"albumNm" json:"albumName"`
	SwingIDs        []string `bson:"swingIds" json:"swingIds"`
	NumComments     int      `bson:"numComms" json:"numComments"`
}

type AlbumUserTagNote struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`

	TaggerID        string `bson:"tgrId,omitempty" json:"taggerId,omitempty"`
	TaggerFirstName string `bson:"ftgrFirst,omitempty" json:"taggerFirstName,omitempty"`
	TaggerUserName  string `bson:"tgrUsr,omitempty" json:"taggerUserName,omitempty"`
	AlbumID         string `bson:"albumId" json:"albumId"`
	AlbumName       string `bson:"albumNm" json:"albumName"`
	SwingID         string `bson:"swingId,omitempty" json:"swingId,omitempty"`
}
