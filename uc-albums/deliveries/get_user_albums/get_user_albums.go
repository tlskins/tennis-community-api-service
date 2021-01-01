package main

import (
	"fmt"
	"net/http"

	api "github.com/tennis-community-api-service/pkg/lambda"
	alb "github.com/tennis-community-api-service/uc-albums"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ucAlb, err := alb.Init()
	if err != nil {
		fmt.Printf(err.Error())
	}
	api.CheckError(http.StatusInternalServerError, err)
	lambda.Start(ucAlb.GetUserAlbums)
}
