package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/davecgh/go-spew/spew"
)

type Request events.APIGatewayProxyRequest

type Response events.APIGatewayProxyResponse

func Parse(req *Request, out interface{}) {
	var jsonBytes []byte
	var err error
	if req.Body == "" {
		if jsonBytes, err = json.Marshal(req.QueryStringParameters); err != nil {
			Abort(http.StatusUnprocessableEntity, err)
		}
	} else {
		jsonBytes = []byte(req.Body)
	}
	if err = json.Unmarshal(jsonBytes, out); err != nil {
		Abort(http.StatusUnprocessableEntity, err)
	}
}

type Validator interface {
	Validate() error
}

func ParseAndValidate(req *Request, out Validator) {
	Parse(req, out)
	if err := out.Validate(); err != nil {
		Abort(http.StatusUnprocessableEntity, err)
	}
}

func GetPathParam(param string, req *Request) (string, error) {
	res, ok := req.PathParameters[param]
	if !ok {
		return "", fmt.Errorf("Param %s not found", param)
	}
	return res, nil
}

type Responder struct {
	Origin string
}

func corsHeaders(req *Request, originsStr string) map[string]string {
	spew.Dump(req.Headers["origin"])
	spew.Dump(originsStr)

	origins := strings.Split(originsStr, ",")
	respOrigin := origins[0]
	for _, origin := range origins {
		if origin == req.Headers["origin"] {
			respOrigin = origin
		}
	}
	spew.Dump(respOrigin)

	return map[string]string{
		"Access-Control-Allow-Origin":      respOrigin,
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "OPTIONS,POST,GET",
		"Access-Control-Allow-Headers":     "Content-Type",
	}
}

// Fail returns an internal server error with the error message
func (r Responder) Fail(req *Request, msg string, status int) (Response, error) {
	e := make(map[string]string, 0)
	e["message"] = msg

	// We don't need to worry about this error,
	// as we're controlling the input.
	body, _ := json.Marshal(e)

	return Response{
		Body:       string(body),
		Headers:    corsHeaders(req, r.Origin),
		StatusCode: status,
	}, nil
}

// Success returns a valid response
func (r Responder) Success(req *Request, data interface{}, status int) (Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return r.Fail(req, err.Error(), http.StatusInternalServerError)
	}

	return Response{
		Body:       string(body),
		Headers:    corsHeaders(req, r.Origin),
		StatusCode: status,
	}, nil
}

func (responder Responder) HandleRequest(handle func(context.Context, *Request) (Response, error)) func(context.Context, *Request) (Response, error) {
	return func(ctx context.Context, req *Request) (resp Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				if e, ok := r.(Error); ok {
					resp, err = responder.Fail(req, e.String(), e.Code)
				} else if e, ok := r.(error); ok {
					resp, err = responder.Fail(req, e.Error(), http.StatusInternalServerError)
				} else {
					resp, err = responder.Fail(req, "unknown error", http.StatusInternalServerError)
				}
			}
		}()
		resp, err = handle(ctx, req)
		return
	}
}
