package types

import (
	"time"
)

type RecentFlaggedCommentsReq struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Resolved *bool     `json:"resolved,omitempty"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}

type RecentFlaggedAlbumsReq struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Resolved *bool     `json:"resolved,omitempty"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}
