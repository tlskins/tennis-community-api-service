package ucalbums

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"

	alb "github.com/tennis-community-api-service/albums"
	mod "github.com/tennis-community-api-service/moderation"
	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
)

// deploy

type UCService struct {
	alb  *alb.AlbumsService
	mod  *mod.ModerationService
	jwt  *auth.JWTService
	Resp *api.Responder
}

func Init() (svc *UCService, err error) {
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	albumsDBName := os.Getenv("ALBUMS_DB_NAME")
	albumsDBHost := os.Getenv("ALBUMS_DB_HOST")
	albumsDBUser := os.Getenv("ALBUMS_DB_USER")
	albumsDBPwd := os.Getenv("ALBUMS_DB_PWD")
	modDBName := os.Getenv("MODERATION_DB_NAME")
	modDBHost := os.Getenv("MODERATION_DB_HOST")
	modDBUser := os.Getenv("MODERATION_DB_USER")
	modDBPwd := os.Getenv("MODERATION_DB_PWD")
	jwtKeyPath := os.Getenv("JWT_KEY_PATH")
	jwtSecretPath := os.Getenv("JWT_SECRET_PATH")
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	var albSvc *alb.AlbumsService
	if albSvc, err = alb.Init(albumsDBName, albumsDBHost, albumsDBUser, albumsDBPwd); err != nil {
		return
	}

	var modSvc *mod.ModerationService
	modSvc, err = mod.Init(modDBName, modDBHost, modDBUser, modDBPwd)
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

	// init responder
	responder := &api.Responder{OriginStr: allowedOrigin}

	svc = &UCService{
		mod:  modSvc,
		alb:  albSvc,
		jwt:  jwt,
		Resp: responder,
	}
	return
}
