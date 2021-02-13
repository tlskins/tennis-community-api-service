package types

import (
	"errors"
)

type GetSwingUploadURLReq struct {
	FileName string `json:"fileName"`
}

func (r GetSwingUploadURLReq) Validate() error {
	if len(r.FileName) == 0 {
		return errors.New("Missing file name")
	}
	return nil
}

type CreateSwingUploadReq struct {
	OriginalURL         string   `bson:"origUrl" json:"originalURL"`
	AlbumName           string   `bson:"albNm" json:"albumName"`
	IsPublic            bool     `bson:"public" json:"isPublic"`
	IsViewableByFriends bool     `bson:"frndView" json:"isViewableByFriends"`
	FriendIDs           []string `bson:"frndIds" json:"friendIds"`
}

func (r CreateSwingUploadReq) Validate() error {
	if r.OriginalURL == "" {
		return errors.New("Missing original url")
	}
	return nil
}

type UploadClipEvent struct {
	ResponsePayload struct {
		StatusCode int `json:"statusCode"`
		Body       struct {
			Bucket       string            `json:"bucket"`
			UserID       string            `json:"userId"`
			UploadID     string            `json:"uploadId"`
			SourceLength float64           `json:"sourceLength"` // seconds
			SourceSize   int64             `json:"sourceSize"`   // bytes
			Outputs      []*ClipUploadMeta `json:"outputs"`
		} `json:"body"`
	} `json:"responsePayload"`
}

type ClipUploadMeta struct {
	Path         string `json:"path"`
	FileName     string `json:"fileName"`
	Number       int    `json:"number"`
	StartSeconds int    `json:"startSec"`
	EndSeconds   int    `json:"endSec"`
}

type UploadSwingEvent struct {
	ResponsePayload struct {
		StatusCode int               `json:"statusCode"`
		Headers    map[string]string `json:"headers"`
		Body       struct {
			UploadID string         `json:"uploadId"`
			UserID   string         `json:"userId"`
			Clip     int            `json:"clipNum"`
			Bucket   string         `json:"bucket"`
			Outputs  []*UploadSwing `json:"outputs"`
		} `json:"body"`
	} `json:"responsePayload"`
}

func (u UploadSwingEvent) Outputs() (videos []string, gifs []string, jpgs []string, txts []string) {
	videos = make([]string, len(u.ResponsePayload.Body.Outputs))
	gifs = make([]string, len(u.ResponsePayload.Body.Outputs))
	jpgs = make([]string, len(u.ResponsePayload.Body.Outputs))
	txts = make([]string, len(u.ResponsePayload.Body.Outputs))
	for i, swing := range u.ResponsePayload.Body.Outputs {
		videos[i] = swing.Video
		gifs[i] = swing.Gif
		jpgs[i] = swing.Jpg
		txts[i] = swing.Txt
	}
	return
}

type UploadSwing struct {
	Video string `json:"video"`
	Gif   string `json:"gif"`
	Jpg   string `json:"jpg"`
	Txt   string `json:"txt"`
}
