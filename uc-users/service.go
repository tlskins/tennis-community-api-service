package users

import (
	"github.com/tennis-community-api-service/pkg/auth"
	usr "github.com/tennis-community-api-service/users"
)

type UCService struct {
	usr *usr.UsersService
	jwt *auth.JWTService
}

func Init(usersDBName, usersDBHost, usersMUser, usersDBPwd string) (svc *UCService, err error) {
	var usrSvc *usr.UsersService
	if usrSvc, err = usr.Init(usersDBName, usersDBHost, usersMUser, usersDBPwd); err != nil {
		return
	}

	svc = &UCService{usr: usrSvc}
	return
}
