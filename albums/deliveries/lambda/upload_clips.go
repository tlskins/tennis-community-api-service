package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tennis-community-api-service/albums"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
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

func HandleRequest(ctx context.Context, event UploadClipsRequest) (string, error) {
	fmt.Println("begin handle request")
	usecase, err := albums.Init()
	if err != nil {
		return "", err
	}
	body := event.ResponsePayload.Body
	usecase.CreateUploadClipVideos(ctx, body.Bucket, body.Outputs)
	return "success", nil
}

func main() {
	lambda.Start(HandleRequest)
}
