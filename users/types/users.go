package types

import (
	"github.com/tennis-community-api-service/pkg/enums"
	"time"
)

type User struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Email     string           `bson:"em" json:"email"`
	FirstName string           `bson:"fnm" json:"firstName"`
	LastName  string           `bson:"lnm" json:"lastName"`
	Status    enums.UserStatus `bson:"status" json:"status"`

	// auth
	EncryptedPassword string     `bson:"pwd" json:"-"`
	LastLoggedIn      *time.Time `bson:"lstLogIn" json:"lastLoggedIn"`
	AuthToken         string     `bson:"-" json:"authToken"`

	// permissions
	AllowSuggestions   bool `bson:"allowSug" json:"allowSuggestions"`
	AllowFlagging      bool `bson:"allowFlag" json:"allowFlagging"`
	AllowPublicAlbums  bool `bson:"allowPubAlb" json:"allowPublicAlbums"`
	AllowAlbumComments bool `bson:"allowAlbCom" json:"allowAlbumComments"`
	AllowVideoComments bool `bson:"allowVidCom" json:"allowVideoComments"`

	// social
	FriendIds      []string         `bson:"friendIds" json:"friendIds"`
	FriendRequests []*FriendRequest `bson:"frndReqs" json:"friendRequests"`

	// notifications
	UploadNotes []*UploadNote `bson:"upNotes" json:"uploadNotifications"`
}

func (u User) GetAuthables() (id, email string, conf bool) {
	return u.ID, u.Email, u.Status != enums.UserStatusPending
}

type UpdateUser struct {
	ID string `bson:"-" json:"id"`

	UpdatedAt *time.Time        `bson:"updAt,omitempty" json:"updatedAt,omitempty"`
	Email     *string           `bson:"em,omitempty" json:"email,omitempty"`
	FirstName *string           `bson:"fnm,omitempty" json:"firstName,omitempty"`
	LastName  *string           `bson:"lnm,omitempty" json:"lastName,omitempty"`
	Status    *enums.UserStatus `bson:"status,omitempty" json:"status,omitempty"`

	// auth
	EncryptedPassword *string    `bson:"pwd,omitempty" json:"-"`
	LastLoggedIn      *time.Time `bson:"lstLogIn,omitempty" json:"lastLoggedIn,omitempty"`
	AuthToken         *string    `bson:"-" json:"authToken,omitempty"`

	// permissions
	AllowSuggestions   *bool `bson:"allowSug,omitempty" json:"allowSuggestions,omitempty"`
	AllowFlagging      *bool `bson:"allowFlag,omitempty" json:"allowFlagging,omitempty"`
	AllowPublicAlbums  *bool `bson:"allowPubAlb,omitempty" json:"allowPublicAlbums,omitempty"`
	AllowAlbumComments *bool `bson:"allowAlbCom,omitempty" json:"allowAlbumComments,omitempty"`
	AllowVideoComments *bool `bson:"allowVidCom,omitempty" json:"allowVideoComments,omitempty"`

	// social
	FriendIds      *[]string         `bson:"friendIds,omitempty" json:"friendIds,omitempty"`
	FriendRequests *[]*FriendRequest `bson:"frndReqs,omitempty" json:"friendRequests,omitempty"`

	// notifications
	UploadNotes *[]*UploadNote `bson:"upNotes,omitempty" json:"uploadNotifications,omitempty"`
}
