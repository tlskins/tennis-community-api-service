package main

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-lambda-go/events"
)

type SampleEvent struct {
	ResponsePayload struct {
		StatusCode int `json:"statusCode"`
		Body       struct {
			Bucket  string   `json:"bucket"`
			Outputs []string `json:"outputs"`
		} `json:"body"`
	} `json:"responsePayload"`
}

func HandleRequest(ctx context.Context, event map[string]interface{}) (string, error) {
	fmt.Println("begin handle request")
	spew.Dump(event)
	return fmt.Sprintf("%+v", event), nil
}

func main() {
	lambda.Start(HandleRequest)
}
