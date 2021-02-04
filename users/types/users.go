package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/tennis-community-api-service/pkg/enums"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	UserName   string           `bson:"usrNm" json:"userName"`
	Email      string           `bson:"em" json:"email"`
	LowerEmail string           `bson:"lowEm" json:"lowerEmail"`
	FirstName  string           `bson:"fnm" json:"firstName"`
	LastName   string           `bson:"lnm" json:"lastName"`
	Status     enums.UserStatus `bson:"status" json:"status"`
	IconNumber int              `bson:"icon" json:"iconNumber"`

	// auth
	EncryptedPassword string     `bson:"pwd" json:"-"`
	LastLoggedIn      *time.Time `bson:"lstLogIn" json:"lastLoggedIn"`
	AuthToken         string     `bson:"-" json:"authToken"`

	// permissions
	IsAdmin             bool `bson:"admin" json:"isAdmin"`
	DisableComments     bool `bson:"noComms" json:"disableComments"`
	DisablePublicAlbums bool `bson:"noPubAlb" json:"disablePublicAlbums"`
	Warnings            int  `bson:"warns" json:"warnings"`

	// social
	FriendIds      []string         `bson:"friendIds" json:"friendIds"`
	FriendRequests []*FriendRequest `bson:"frndReqs" json:"friendRequests"`

	// notifications
	UploadNotes      []*UploadNote  `bson:"upNotes" json:"uploadNotifications"`
	FriendNotes      []*FriendNote  `bson:"frndNotes" json:"friendNotifications"`
	CommentNotes     []*CommentNote `bson:"commentNotes" json:"commentNotifications"`
	MyRecentComments []*CommentNote `bson:"recComms" json:"myRecentComments"`

	// profile
	IsPublic    bool     `bson:"pub" json:"isPublic"`
	RightHanded *bool    `bson:"rgtHnd" json:"rightHanded"`
	BirthYear   *int     `bson:"birthYr" json:"birthYear"`
	Gender      *string  `bson:"gender" json:"gender"`
	USTALevel   *float64 `bson:"ustaLvl" json:"ustaLevel"`
	State       *string  `bson:"state" json:"state"`
	City        *string  `bson:"city" json:"city"`
	Longitude   *float64 `bson:"lng" json:"longitude"`
	Latitude    *float64 `bson:"lat" json:"latitue"`
}

func (u User) Name() string {
	name := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	return strings.Trim(name, " ")
}

func (u User) GetAuthables() (id, email string, conf, isAdmin bool) {
	return u.ID, u.Email, u.Status != enums.UserStatusPending, u.IsAdmin
}

func (u *User) AddCommentNote(friend *User, albumID, albumName, swingID string) {
	noteFound := false
	for _, note := range u.CommentNotes {
		if note.AlbumID == albumID && note.FriendID == friend.ID {
			noteFound = true
			note.NumComments++
			note.UpdatedAt = time.Now()
			if swingID != "" {
				note.SwingIDs = append(note.SwingIDs, swingID)
			}
			break
		}
	}
	if !noteFound {
		now := time.Now()
		note := &CommentNote{
			ID:              uuid.NewV4().String(),
			CreatedAt:       now,
			UpdatedAt:       now,
			FriendID:        friend.ID,
			FriendFirstName: friend.FirstName,
			FriendUserName:  friend.UserName,
			AlbumID:         albumID,
			AlbumName:       albumName,
			NumComments:     1,
		}
		if swingID != "" {
			note.SwingIDs = []string{swingID}
		}
		u.CommentNotes = append(u.CommentNotes, note)
	}
}

func (u *User) AddMyRecentComment(albumID, albumName, swingID string) *CommentNote {
	now := time.Now()
	cutoff := now.Add(time.Duration(-3) * time.Hour)
	recentCutoff := now.AddDate(0, 0, -14)
	recentNotes := []*CommentNote{}
	var note *CommentNote

	// keep recent and find existing note
	for _, n := range u.MyRecentComments {
		if n.AlbumID == albumID && n.CreatedAt.After(cutoff) {
			n.NumComments++
			n.UpdatedAt = time.Now()
			if swingID != "" {
				n.SwingIDs = append(n.SwingIDs, swingID)
			}
			note = n
		}
		if n.CreatedAt.After(recentCutoff) {
			recentNotes = append(recentNotes, n)
		}
	}

	// create note if new comment
	if note == nil {
		note = &CommentNote{
			ID:          uuid.NewV4().String(),
			CreatedAt:   now,
			UpdatedAt:   now,
			AlbumID:     albumID,
			AlbumName:   albumName,
			NumComments: 1,
		}
		if swingID != "" {
			note.SwingIDs = []string{swingID}
		}
		recentNotes = append(recentNotes, note)
	}

	u.MyRecentComments = recentNotes
	return note
}

type UpdateUser struct {
	ID        string    `bson:"-" json:"id"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	Email  *string           `bson:"em,omitempty" json:"email,omitempty"`
	Status *enums.UserStatus `bson:"status,omitempty" json:"status,omitempty"`

	// auth
	EncryptedPassword *string    `bson:"pwd,omitempty" json:"-"`
	LastLoggedIn      *time.Time `bson:"lstLogIn,omitempty" json:"lastLoggedIn,omitempty"`
	AuthToken         *string    `bson:"-" json:"authToken,omitempty"`

	// permissions
	DisableComments     bool `bson:"noComms,omitempty" json:"disableComments,omitempty"`
	DisablePublicAlbums bool `bson:"noPubAlb,omitempty" json:"disablePublicAlbums,omitempty"`
	Warnings            int  `bson:"warns,omitempty" json:"warnings,omitempty"`

	// social
	FriendIds      *[]string         `bson:"friendIds,omitempty" json:"friendIds,omitempty"`
	FriendRequests *[]*FriendRequest `bson:"frndReqs,omitempty" json:"friendRequests,omitempty"`

	// notifications
	UploadNotes      *[]*UploadNote  `bson:"upNotes,omitempty" json:"uploadNotifications,omitempty"`
	FriendNotes      *[]*FriendNote  `bson:"frndNotes,omitempty" json:"friendNotifications,omitempty"`
	CommentNotes     *[]*CommentNote `bson:"commentNotes,omitempty" json:"commentNotifications,omitempty"`
	MyRecentComments *[]*CommentNote `bson:"recComms,omitempty" json:"myRecentComments,omitempty"`
}

type UpdateUserProfile struct {
	ID        string    `bson:"-" json:"id"`
	UpdatedAt time.Time `bson:"updAt" json:"updatedAt"`

	UserName   string `bson:"usrNm" json:"userName"`
	FirstName  string `bson:"fnm" json:"firstName"`
	LastName   string `bson:"lnm" json:"lastName"`
	IconNumber int    `bson:"icon" json:"iconNumber"`

	// profile
	IsPublic    bool     `bson:"pub" json:"isPublic"`
	RightHanded *bool    `bson:"rgtHnd" json:"rightHanded"`
	BirthYear   *int     `bson:"birthYr" json:"birthYear"`
	Gender      *string  `bson:"gender" json:"gender"`
	USTALevel   *float64 `bson:"ustaLvl" json:"ustaLevel"`
	State       *string  `bson:"state" json:"state"`
	City        *string  `bson:"city" json:"city"`
	Longitude   *float64 `bson:"lng" json:"longitude"`
	Latitude    *float64 `bson:"lat" json:"latitue"`
}
