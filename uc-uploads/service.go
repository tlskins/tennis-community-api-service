package uploads

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"

	alb "github.com/tennis-community-api-service/albums"
	"github.com/tennis-community-api-service/pkg/auth"
	"github.com/tennis-community-api-service/pkg/email"
	api "github.com/tennis-community-api-service/pkg/lambda"
	up "github.com/tennis-community-api-service/uploads"
	usr "github.com/tennis-community-api-service/users"
)

// deploy

type UCService struct {
	up          *up.UploadsService
	alb         *alb.AlbumsService
	usr         *usr.UsersService
	jwt         *auth.JWTService
	emailClient *email.EmailClient
	Resp        *api.Responder
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
	usersDBName := os.Getenv("USERS_DB_NAME")
	usersDBHost := os.Getenv("USERS_DB_HOST")
	usersDBUser := os.Getenv("USERS_DB_USER")
	usersDBPwd := os.Getenv("USERS_DB_PWD")
	jwtKeyPath := os.Getenv("JWT_KEY_PATH")
	jwtSecretPath := os.Getenv("JWT_SECRET_PATH")
	fromEmail := os.Getenv("FROM_EMAIL")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PWD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	cdnURL := os.Getenv("CDN_URL")

	var upSvc *up.UploadsService
	if upSvc, err = up.Init(uploadsDBName, uploadsDBHost, uploadsDBUser, uploadsDBPwd, cdnURL); err != nil {
		return
	}

	var albSvc *alb.AlbumsService
	if albSvc, err = alb.Init(albumsDBName, albumsDBHost, albumsDBUser, albumsDBPwd); err != nil {
		return
	}

	var usrSvc *usr.UsersService
	usrSvc, err = usr.Init(usersDBName, usersDBHost, usersDBUser, usersDBPwd)
	if err != nil {
		return
	}

	// Email
	var emailClient *email.EmailClient
	emailClient, err = email.NewEmailClient(fromEmail, smtpHost, smtpPort, smtpUser, smtpPassword)
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
		up:          upSvc,
		alb:         albSvc,
		usr:         usrSvc,
		jwt:         jwt,
		emailClient: emailClient,
		Resp:        responder,
	}
	return
}
