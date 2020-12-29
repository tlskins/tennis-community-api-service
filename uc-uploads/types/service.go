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
