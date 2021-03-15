package types

import (
	"errors"
	"regexp"
	"time"

	uT "github.com/tennis-community-api-service/users/types"
)

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserReq struct {
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var invalidUserNameChars = []string{"@", "#", " ", "	", "%", "/", `\\`}

func (r CreateUserReq) Validate() error {
	if len(r.UserName) < 3 {
		return errors.New("Username must be at least 3 characters long")
	}
	if len(r.UserName) > 20 {
		return errors.New("Username must be shorter than 20 characters")
	}
	if len(r.FirstName) == 0 {
		return errors.New("Missing first name")
	}
	if len(r.LastName) == 0 {
		return errors.New("Missing last name")
	}
	if len(r.Email) == 0 {
		return errors.New("Missing email")
	}
	if len(r.Password) == 0 {
		return errors.New("Missing password")
	}
	for _, char := range invalidUserNameChars {
		re := regexp.MustCompile(char)
		if re.Match([]byte(r.UserName)) {
			return errors.New(`Username cannot include any of the following special characters: spaces, tabs, @, #, %, \, /`)
		}
	}
	return nil
}

type UpdateUserProfileReq uT.UpdateUserProfile

func (r UpdateUserProfileReq) Validate() error {
	if len(r.UserName) < 3 {
		return errors.New("Username must be at least 3 characters long")
	}
	if len(r.UserName) > 20 {
		return errors.New("Username must be shorter than 20 characters")
	}
	if len(r.FirstName) == 0 {
		return errors.New("Missing first name")
	}
	if len(r.LastName) == 0 {
		return errors.New("Missing last name")
	}
	for _, char := range invalidUserNameChars {
		re := regexp.MustCompile(char)
		if re.Match([]byte(r.UserName)) {
			return errors.New(`Username cannot include any of the following special characters: spaces, tabs, @, #, %, \, /`)
		}
	}

	return nil
}

type RemoveNotificationReq struct {
	UploadNoteID       string `json:"uploadNotificationId"`
	FriendNoteID       string `json:"friendNotificationId"`
	CommentNoteID      string `json:"commentNotificationId"`
	AlbumUserTagNoteID string `json:"albumUserTagNotificationId"`
}

func (r RemoveNotificationReq) Validate() error {
	count := 0
	if len(r.UploadNoteID) > 0 {
		count++
	}
	if len(r.FriendNoteID) > 0 {
		count++
	}
	if len(r.CommentNoteID) > 0 {
		count++
	}
	if len(r.AlbumUserTagNoteID) > 0 {
		count++
	}
	if count > 1 {
		return errors.New("Can only remove one notification at a time")
	}
	return nil
}

type AcceptFriendReq struct {
	Accept bool `json:"accept"`
}

type SearchFriendsReq struct {
	IDs    *[]string `json:"ids"`
	Search *string   `json:"search"`
	Offset int       `json:"offset"`
	Limit  int       `json:"limit"`
}

type RecentUsersReq struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Offset string    `json:"offset"`
	Limit  string    `json:"limit"`
}

type UserConfirmationReq struct {
	ID        string `json:"id"`
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserInvitationReq struct {
	Email     string `json:"email"`
	InviterID string `json:"inviterId"`
	URL       string `json:"url"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
