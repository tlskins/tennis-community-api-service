package uploads

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"

	alb "github.com/tennis-community-api-service/albums"
	"github.com/tennis-community-api-service/pkg/auth"
	up "github.com/tennis-community-api-service/uploads"
)

type UCService struct {
	up  *up.UploadsService
	alb *alb.AlbumsService
	jwt *auth.JWTService
}

func Init() (svc *UCService, err error) {
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	uploadsDBName := os.Getenv("UPLOADS_DB_NAME")
	uploadsDBHost := os.Getenv("UPLOADS_DB_HOST")
	uploadsDBUser := os.Getenv("UPLOADS_DB_USER")
	uploadsDBPwd := os.Getenv("UPLOADS_DB_PWD")
	albumsDBName := os.Getenv("ALBUMS_DB_NAME")
	albumsDBHost := os.Getenv("ALBUMS_DB_HOST")
	albumsDBUser := os.Getenv("ALBUMS_DB_USER")
	albumsDBPwd := os.Getenv("ALBUMS_DB_PWD")
	jwtKeyPath := os.Getenv("JWT_KEY_PATH")
	jwtSecretPath := os.Getenv("JWT_SECRET_PATH")

	var upSvc *up.UploadsService
	if upSvc, err = up.Init(uploadsDBName, uploadsDBHost, uploadsDBUser, uploadsDBPwd); err != nil {
		return
	}

	var albSvc *alb.AlbumsService
	if albSvc, err = alb.Init(albumsDBName, albumsDBHost, albumsDBUser, albumsDBPwd); err != nil {
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
		up:  upSvc,
		alb: albSvc,
		jwt: jwt,
	}
	return
}
