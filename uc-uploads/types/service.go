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
	OriginalURL string `bson:"origUrl" json:"originalURL"`
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
			Bucket  string   `json:"bucket"`
			Outputs []string `json:"outputs"`
		} `json:"body"`
	} `json:"responsePayload"`
}

type UploadSwingEvent struct {
	ResponsePayload struct {
		StatusCode int `json:"statusCode"`
		Body       struct {
			Bucket  string         `json:"bucket"`
			Outputs []*UploadSwing `json:"outputs"`
		} `json:"body"`
	} `json:"responsePayload"`
}

func (u UploadSwingEvent) Outputs() (videos []string, gifs []string, jpgs []string) {
	videos = make([]string, len(u.ResponsePayload.Body.Outputs))
	gifs = make([]string, len(u.ResponsePayload.Body.Outputs))
	jpgs = make([]string, len(u.ResponsePayload.Body.Outputs))
	for i, swing := range u.ResponsePayload.Body.Outputs {
		videos[i] = swing.Video
		gifs[i] = swing.Gif
		jpgs[i] = swing.Jpg
	}
	return
}

type UploadSwing struct {
	Video string `json:"video"`
	Gif   string `json:"gif"`
	Jpg   string `json:"jpg"`
}
