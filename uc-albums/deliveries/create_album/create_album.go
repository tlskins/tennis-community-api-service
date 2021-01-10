package main

import (
	"net/http"

	api "github.com/tennis-community-api-service/pkg/lambda"
	alb "github.com/tennis-community-api-service/uc-albums"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ucAlb, err := alb.Init()
	api.CheckError(http.StatusInternalServerError, err)
	handler := api.HandleRequest(ucAlb.CreateAlbum)
	lambda.Start(handler)
}
