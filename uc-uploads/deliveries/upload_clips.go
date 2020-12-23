package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	up "github.com/tennis-community-api-service/uc-uploads"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

type UploadClipsRequest struct {
	ResponsePayload struct {
		StatusCode int `json:"statusCode"`
		Body       struct {
			Bucket  string   `json:"bucket"`
			Outputs []string `json:"outputs"`
		} `json:"body"`
	} `json:"responsePayload"`
}

func HandleRequest(ctx context.Context, event UploadClipsRequest) (resp string, err error) {
	fmt.Println("begin handle request")

	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	uploadsDBName := os.Getenv("UPLOADS_DB_NAME")
	uploadsDBHost := os.Getenv("UPLOADS_HOST")
	uploadsDBUser := os.Getenv("UPLOADS_USER")
	uploadsDBPwd := os.Getenv("UPLOADS_PWD")

	var ucUp *up.UCService
	ucUp, err = up.Init(
		uploadsDBName,
		uploadsDBHost,
		uploadsDBUser,
		uploadsDBPwd,
	)
	if err != nil {
		return
	}

	body := event.ResponsePayload.Body
	_, err = ucUp.CreateUploadClipVideos(ctx, body.Bucket, body.Outputs)
	return
}

func main() {
	lambda.Start(HandleRequest)
}
