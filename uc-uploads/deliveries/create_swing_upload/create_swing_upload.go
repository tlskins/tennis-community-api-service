package main

import (
	"net/http"

	api "github.com/tennis-community-api-service/pkg/lambda"
	up "github.com/tennis-community-api-service/uc-uploads"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ucUp, err := up.Init()
	api.CheckError(http.StatusInternalServerError, err)
	handler := api.HandleRequest(ucUp.CreateSwingUpload)
	lambda.Start(handler)
}
