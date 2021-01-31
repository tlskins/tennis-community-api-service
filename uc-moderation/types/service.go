package types

import (
	"time"
)

type RecentFlaggedCommentsReq struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Resolved *string   `json:"resolved,omitempty"`
	Limit    string    `json:"limit"`
	Offset   string    `json:"offset"`
}

type RecentFlaggedAlbumsReq struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Resolved *string   `json:"resolved,omitempty"`
	Limit    string    `json:"limit"`
	Offset   string    `json:"offset"`
}
