package types

import (
	"time"
)

type UserConfirmation struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`

	// normal user
	UserID string `bson:"usrId,omitempty" json:"userId,omitempty"`

	// invited user
	Email     string `bson:"em,omitempty" json:"email,omitempty"`
	InviterID string `bson:"inviter,omitempty" json:"inviterId,omitempty"`
	URL       string `bson:"url,omitempty" json:"url,omitempty"`
	FirstName string `bson:"fnm,omitempty" json:"firstName,omitempty"`
	LastName  string `bson:"lnm,omitempty" json:"lastName,omitempty"`
}
