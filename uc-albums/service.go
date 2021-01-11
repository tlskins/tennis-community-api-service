package uploads

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"

	alb "github.com/tennis-community-api-service/albums"
	"github.com/tennis-community-api-service/pkg/auth"
	usr "github.com/tennis-community-api-service/users"
)

type UCService struct {
	usr *usr.UsersService
	alb *alb.AlbumsService
	jwt *auth.JWTService
}

func Init() (svc *UCService, err error) {
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	albumsDBName := os.Getenv("ALBUMS_DB_NAME")
	albumsDBHost := os.Getenv("ALBUMS_DB_HOST")
	albumsDBUser := os.Getenv("ALBUMS_DB_USER")
	albumsDBPwd := os.Getenv("ALBUMS_DB_PWD")
	usersDBName := os.Getenv("USERS_DB_NAME")
	usersDBHost := os.Getenv("USERS_DB_HOST")
	usersDBUser := os.Getenv("USERS_DB_USER")
	usersDBPwd := os.Getenv("USERS_DB_PWD")
	jwtKeyPath := os.Getenv("JWT_KEY_PATH")
	jwtSecretPath := os.Getenv("JWT_SECRET_PATH")

	var albSvc *alb.AlbumsService
	if albSvc, err = alb.Init(albumsDBName, albumsDBHost, albumsDBUser, albumsDBPwd); err != nil {
		return
	}

	var usrSvc *usr.UsersService
	usrSvc, err = usr.Init(usersDBName, usersDBHost, usersDBUser, usersDBPwd)
	if err != nil {
		return
	}

	// Init jwt service
	jwtKey, err := ioutil.ReadFile(jwtKeyPath)
	if err != nil {
		return nil, err
	}
	jwtSecret, err := ioutil.ReadFile(jwtSecretPath)
	if err != nil {
		return nil, err
	}
	jwt, err := auth.NewJWTService(auth.JWTServiceConfig{
		Key:    jwtKey,
		Secret: jwtSecret,
	})
	if err != nil {
		return nil, err
	}

	svc = &UCService{
		usr: usrSvc,
		alb: albSvc,
		jwt: jwt,
	}
	return
}
