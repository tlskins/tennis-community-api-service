package uploads

import (
	up "github.com/tennis-community-api-service/uploads"
)

type UCService struct {
	up *up.UploadsService
}

func Init(uploadsDBName, uploadsDBHost, uploadsMUser, uploadsDBPwd string) (svc *UCService, err error) {
	var upSvc *up.UploadsService
	if upSvc, err = up.Init(uploadsDBName, uploadsDBHost, uploadsMUser, uploadsDBPwd); err != nil {
		return
	}

	svc = &UCService{up: upSvc}
	return
}
