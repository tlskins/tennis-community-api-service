package main

import (
	"net/http"

	api "github.com/tennis-community-api-service/pkg/lambda"
	mod "github.com/tennis-community-api-service/uc-moderation"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ucMod, err := mod.Init()
	api.CheckError(http.StatusInternalServerError, err)
	lambda.Start(ucMod.Resp.HandleRequest(ucMod.RecentFlaggedComments))
}
