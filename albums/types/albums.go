package types

import (
	"errors"
	"sort"
	"time"

	"github.com/tennis-community-api-service/pkg/enums"
)

type Album struct {
	ID        string            `bson:"_id" json:"id"`
	UserID    string            `bson:"userId" json:"userId"`
	UploadKey string            `bson:"upKey,omitempty" json:"uploadKey,omitempty"`
	Name      string            `bson:"nm" json:"name"`
	CreatedAt time.Time         `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time         `bson:"updAt" json:"updatedAt"`
	Status    enums.AlbumStatus `bson:"status" json:"status"`

	Tags        []string      `bson:"tags" json:"tags"`
	SwingVideos []*SwingVideo `bson:"swingVids" json:"swingVideos"`
	Comments    []*Comment    `bson:"cmnts" json:"comments"`

	IsPublic            bool     `bson:"public" json:"isPublic"`
	IsViewableByFriends bool     `bson:"frndView" json:"isViewableByFriends"`
	FriendIDs           []string `bson:"frndIds" json:"friendIds"`

	HomeApproved bool `bson:"home" json:"homeApproved"`
	IsPro        bool `bson:"pro" json:"pro"`

	// file upload only
	SourceLength float64 `bson:"srcLen,omitempty" json:"sourceLength,omitempty"` // seconds
	SourceSize   int64   `bson:"srcSize,omitempty" json:"sourceSize,omitempty"`  // bytes
}

func (a *Album) CalculateMetrics() {
	// sort swings by upload key
	sort.SliceStable(a.SwingVideos, func(i, j int) bool {
		return a.SwingVideos[i].UploadKey < a.SwingVideos[j].UploadKey
	})
	// sort swings by timestamp
	sort.SliceStable(a.SwingVideos, func(i, j int) bool {
		return a.SwingVideos[i].TimestampSeconds < a.SwingVideos[j].TimestampSeconds
	})
	// calculate rallies
	lastTime := 0
	currRally := 1
	rallyWindow := 6 // seconds
	for i, swing := range a.SwingVideos {
		if i != 0 && swing.TimestampSeconds-lastTime >= rallyWindow {
			currRally++
		}
		swing.Rally = currRally
		lastTime = swing.TimestampSeconds
	}
}

type UpdateAlbum struct {
	ID        string    `bson:"-" json:"id"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Name        *string            `bson:"nm,omitempty" json:"name,omitempty"`
	Status      *enums.AlbumStatus `bson:"status,omitempty" json:"status,omitempty"`
	Tags        *[]string          `bson:"tags,omitempty" json:"tags,omitempty"`
	SwingVideos *[]*SwingVideo     `bson:"swingVids,omitempty" json:"swingVideos,omitempty"`
	Comments    *[]*Comment        `bson:"cmnts,omitempty" json:"comments,omitempty"`

	IsPublic            *bool     `bson:"public,omitempty" json:"isPublic,omitempty"`
	IsViewableByFriends *bool     `bson:"frndView,omitempty" json:"isViewableByFriends,omitempty"`
	FriendIDs           *[]string `bson:"frndIds,omitempty" json:"friendIds,omitempty"`

	HomeApproved *bool `bson:"home,omitempty" json:"homeApproved,omitempty"`
	IsPro        *bool `bson:"pro,omitempty" json:"pro,omitempty"`

	// file upload only
	SourceLength *float64 `bson:"srcLen,omitempty" json:"sourceLength,omitempty"` // seconds
	SourceSize   *int64   `bson:"srcSize,omitempty" json:"sourceSize,omitempty"`  // bytes
}

func (a UpdateAlbum) Validate() error {
	if a.Name != nil && *a.Name == "" {
		return errors.New("Name cannot be blank")
	}
	return nil
}
