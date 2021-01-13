package types

import (
	"github.com/tennis-community-api-service/pkg/enums"
	"time"
)

type SwingUpload struct {
	ID                  string                  `bson:"_id" json:"id"`
	CreatedAt           time.Time               `bson:"crAt" json:"createdAt"`
	UpdatedAt           time.Time               `bson:"updAt" json:"updatedAt"`
	UploadKey           string                  `bson:"upKey" json:"uploadKey"`
	UserID              string                  `bson:"usrId" json:"userId"`
	Status              enums.SwingUploadStatus `bson:"status" json:"status"`
	OriginalURL         string                  `bson:"origUrl" json:"originalURL"`
	AlbumID             string                  `bson:"albId,omitempty" json:"albumId,omitempty"`
	AlbumName           string                  `bson:"albNm" json:"albumName"`
	IsPublic            bool                    `bson:"public" json:"isPublic"`
	IsViewableByFriends bool                    `bson:"frndView" json:"isViewableByFriends"`
	FriendIDs           []string                `bson:"frndIds" json:"friendIds"`
	ClipVideos          []*UploadClipVideo      `bson:"clipVids" json:"clipVideos"`
	SwingVideos         []*UploadSwingVideo     `bson:"swingVids" json:"swingVideos"`
}

func (a SwingUpload) IsFinal() bool {
	clipVidsCount := 0
	swingMap := map[int]bool{}
	for _, swing := range a.SwingVideos {
		if !swingMap[swing.ClipID] {
			swingMap[swing.ClipID] = true
			clipVidsCount++
		}
	}
	return clipVidsCount == len(a.ClipVideos)
}

type UpdateSwingUpload struct {
	UploadKey string    `bson:"-" json:"uploadKey"`
	UserID    string    `bson:"usrId" json:"userId"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Status              *enums.SwingUploadStatus `bson:"status,omitempty" json:"status,omitempty"`
	OriginalURL         *string                  `bson:"origUrl,omitempty" json:"originalURL,omitempty"`
	AlbumID             *string                  `bson:"albId,omitempty" json:"albumId,omitempty"`
	AlbumName           *string                  `bson:"albNm,omitempty" json:"albumName,omitempty"`
	IsPublic            *bool                    `bson:"public,omitempty" json:"isPublic,omitempty"`
	IsViewableByFriends *bool                    `bson:"frndView,omitempty" json:"isViewableByFriends,omitempty"`
	FriendIDs           *[]string                `bson:"frndIds,omitempty" json:"friendIds,omitempty"`
	ClipVideos          *[]*UploadClipVideo      `bson:"clipVids,omitempty" json:"clipVideos,omitempty"`
	SwingVideos         *[]*UploadSwingVideo     `bson:"swingVids,omitempty" json:"swingVideos,omitempty"`
}

type UploadClipVideo struct {
	ID        int       `bson:"id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	ClipURL   string    `bson:"clipUrl" json:"clipURL"`
}

type UploadSwingVideo struct {
	ID        string    `bson:"id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`
	ClipID    int       `bson:"clipId" json:"clipId"`
	SwingID   int       `bson:"swId" json:"swingId"`
	CutURL    string    `bson:"cutUrl" json:"cutURL"`
	GifURL    string    `bson:"gifUrl" json:"gifURL"`
	JpgURL    string    `bson:"jpgUrl" json:"jpgURL"`
}

type UpdateUploadSwingVideo struct {
	ID        string    `bson:"-" json:"id"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`
	ClipID    *int      `bson:"clipId,omitempty" json:"clipId,omitempty"`
	CutURL    *string   `bson:"cutUrl,omitempty" json:"cutURL,omitempty"`
}
