package main

import (
	"fmt"
	"net/http"

	api "github.com/tennis-community-api-service/pkg/lambda"
	up "github.com/tennis-community-api-service/uc-uploads"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ucUp, err := up.Init()
	if err != nil {
		fmt.Printf(err.Error())
	}
	api.CheckError(http.StatusInternalServerError, err)
	lambda.Start(ucUp.CreateUploadSwingVideos)
}
