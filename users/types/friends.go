package types

import (
	"time"
)

type Friend struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"crAt" json:"createdAt"`

	UserName  string `bson:"usrNm" json:"userName"`
	FirstName string `bson:"fnm" json:"firstName"`
	LastName  string `bson:"lnm" json:"lastName"`
}

type FriendRequest struct {
	ID         string    `bson:"_id" json:"id"`
	CreatedAt  time.Time `bson:"crAt" json:"createdAt"`
	FromUserID string    `bson:"from" json:"fromUserId"`
	ToUserID   string    `bson:"to" json:"toUserId"`
}
