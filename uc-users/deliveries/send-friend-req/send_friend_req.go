package main

import (
	"net/http"

	api "github.com/tennis-community-api-service/pkg/lambda"
	usr "github.com/tennis-community-api-service/uc-users"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ucUsr, err := usr.Init()
	api.CheckError(http.StatusInternalServerError, err)
	handler := api.HandleRequest(ucUsr.SendFriendRequest)
	lambda.Start(handler)
}
