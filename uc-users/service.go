package users

import (
	"flag"
	"os"

	"github.com/joho/godotenv"

	"github.com/tennis-community-api-service/pkg/auth"
	usr "github.com/tennis-community-api-service/users"
)

type UCService struct {
	usr *usr.UsersService
	jwt *auth.JWTService
}

func Init() (svc *UCService, err error) {
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	usersDBName := os.Getenv("USERS_DB_NAME")
	usersDBHost := os.Getenv("USERS_DB_HOST")
	usersDBUser := os.Getenv("USERS_DB_USER")
	usersDBPwd := os.Getenv("USERS_DB_PWD")

	var usrSvc *usr.UsersService
	if usrSvc, err = usr.Init(usersDBName, usersDBHost, usersDBUser, usersDBPwd); err != nil {
		return
	}

	svc = &UCService{usr: usrSvc}
	return
}
