package main

import (
	"flag"
	"net/http"
	"os"

	api "github.com/tennis-community-api-service/pkg/lambda"
	usr "github.com/tennis-community-api-service/uc-users"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func main() {

	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	usersDBName := os.Getenv("USERS_DB_NAME")
	usersDBHost := os.Getenv("USERS_HOST")
	usersDBUser := os.Getenv("USERS_USER")
	usersDBPwd := os.Getenv("USERS_PWD")

	var ucUsr *usr.UCService
	ucUsr, err := usr.Init(usersDBName, usersDBHost, usersDBUser, usersDBPwd)
	api.CheckError(http.StatusInternalServerError, err)

	handler := api.HandleRequest(ucUsr.CreateUser)
	lambda.Start(handler)
}
