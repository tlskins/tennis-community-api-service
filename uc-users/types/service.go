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
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (r CreateUserReq) Validate() error {
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
}

type FriendIdReq struct {
	FriendID string `json:"friendID"`
}

type AcceptFriendReq struct {
	Accept bool `json:"accept"`
}
