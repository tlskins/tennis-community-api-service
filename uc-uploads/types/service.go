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

type SwingStorageEvent struct {
	ResponsePayload struct {
		StatusCode int `json:"statusCode"`
		Body       struct {
			Bucket  string   `json:"bucket"`
			Outputs []string `json:"outputs"`
		} `json:"body"`
	} `json:"responsePayload"`
}
