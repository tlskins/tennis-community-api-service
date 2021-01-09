package types

import (
	"errors"
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

func (r CreateUserReq) Validate() error {
	if len(r.UserName) < 3 {
		return errors.New("Username must be at least 3 characters long")
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
	return nil
}

type ClearNotificationsReq struct {
	Uploads bool `json:"uploads"`
	Friends bool `json:"friends"`
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
