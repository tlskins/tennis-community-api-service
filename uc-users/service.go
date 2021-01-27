package users

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"

	"github.com/tennis-community-api-service/pkg/auth"
	"github.com/tennis-community-api-service/pkg/email"
	api "github.com/tennis-community-api-service/pkg/lambda"
	usr "github.com/tennis-community-api-service/users"
)

// deploy

type UCService struct {
	usr         *usr.UsersService
	jwt         *auth.JWTService
	emailClient *email.EmailClient
	hostName    string
	Resp        *api.Responder
}

func Init() (svc *UCService, err error) {
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	usersDBName := os.Getenv("USERS_DB_NAME")
	usersDBHost := os.Getenv("USERS_DB_HOST")
	usersDBUser := os.Getenv("USERS_DB_USER")
	usersDBPwd := os.Getenv("USERS_DB_PWD")
	jwtKeyPath := os.Getenv("JWT_KEY_PATH")
	jwtSecretPath := os.Getenv("JWT_SECRET_PATH")
	fromEmail := os.Getenv("FROM_EMAIL")
	emailPwd := os.Getenv("EMAIL_PWD")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PWD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	hostName := os.Getenv("API_HOST")
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	var usrSvc *usr.UsersService
	usrSvc, err = usr.Init(usersDBName, usersDBHost, usersDBUser, usersDBPwd)
	if err != nil {
		return
	}

	// Email
	var emailClient *email.EmailClient
	emailClient, err = email.NewEmailClient(fromEmail, emailPwd, smtpHost, smtpPort, smtpUser, smtpPassword)
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
	responder := &api.Responder{Origin: allowedOrigin}

	svc = &UCService{
		usr:         usrSvc,
		jwt:         jwt,
		emailClient: emailClient,
		hostName:    hostName,
		Resp:        responder,
	}
	return
}
