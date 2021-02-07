package types

import (
	"errors"
	"time"

	"github.com/tennis-community-api-service/pkg/enums"
)

type Album struct {
	ID                  string            `bson:"_id" json:"id"`
	UserID              string            `bson:"userId" json:"userId"`
	UploadKey           string            `bson:"upKey,omitempty" json:"uploadKey,omitempty"`
	Name                string            `bson:"nm" json:"name"`
	CreatedAt           time.Time         `bson:"crAt" json:"createdAt"`
	UpdatedAt           time.Time         `bson:"updAt" json:"updatedAt"`
	Status              enums.AlbumStatus `bson:"status" json:"status"`
	Tags                []string          `bson:"tags" json:"tags"`
	SwingVideos         []*SwingVideo     `bson:"swingVids" json:"swingVideos"`
	HomeApproved        bool              `bson:"home" json:"homeApproved"`
	IsPublic            bool              `bson:"public" json:"isPublic"`
	IsViewableByFriends bool              `bson:"frndView" json:"isViewableByFriends"`
	FriendIDs           []string          `bson:"frndIds" json:"friendIds"`
	Comments            []*Comment        `bson:"cmnts" json:"comments"`
}

type UpdateAlbum struct {
	ID        string    `bson:"-" json:"id"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Name                *string            `bson:"nm,omitempty" json:"name,omitempty"`
	Status              *enums.AlbumStatus `bson:"status,omitempty" json:"status,omitempty"`
	Tags                *[]string          `bson:"tags,omitempty" json:"tags,omitempty"`
	SwingVideos         *[]*SwingVideo     `bson:"swingVids,omitempty" json:"swingVideos,omitempty"`
	HomeApproved        *bool              `bson:"home,omitempty" json:"homeApproved,omitempty"`
	IsPublic            *bool              `bson:"public,omitempty" json:"isPublic,omitempty"`
	IsViewableByFriends *bool              `bson:"frndView,omitempty" json:"isViewableByFriends,omitempty"`
	FriendIDs           *[]string          `bson:"frndIds,omitempty" json:"friendIds,omitempty"`
	Comments            *[]*Comment        `bson:"cmnts,omitempty" json:"comments,omitempty"`
}

func (a UpdateAlbum) Validate() error {
	if a.Name != nil && *a.Name == "" {
		return errors.New("Name cannot be blank")
	}
	return nil
}
