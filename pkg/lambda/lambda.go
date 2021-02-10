package lambda

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/aws/aws-lambda-go/events"
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

type Responder struct {
	OriginStr string
}

func (r Responder) corsHeaders(headers map[string]string) map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      r.Origin(headers),
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "OPTIONS,POST,GET",
		"Access-Control-Allow-Headers":     "Content-Type",
	}
}

func (r Responder) Origin(headers map[string]string) (origin string) {
	origins := strings.Split(r.OriginStr, ",")
	origin = origins[0]
	for _, str := range origins {
		if str == headers["origin"] {
			origin = str
		}
	}
	return
}

// Fail returns an internal server error with the error message
func (r Responder) Fail(headers map[string]string, msg string, status int) (Response, error) {
	e := make(map[string]string, 0)
	e["message"] = msg

	// We don't need to worry about this error,
	// as we're controlling the input.
	body, _ := json.Marshal(e)

	return Response{
		Body:       string(body),
		Headers:    r.corsHeaders(headers),
		StatusCode: status,
	}, nil
}

// Success returns a valid response
func (r Responder) Success(headers map[string]string, data interface{}, status int) (Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return r.Fail(headers, err.Error(), http.StatusInternalServerError)
	}

	return Response{
		Body:       string(body),
		Headers:    r.corsHeaders(headers),
		StatusCode: status,
	}, nil
}

func (responder Responder) HandleRequest(handle func(context.Context, *Request) (Response, error)) func(context.Context, *Request) (Response, error) {
	return func(ctx context.Context, req *Request) (resp Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				if e, ok := r.(Error); ok {
					resp, err = responder.Fail(req.Headers, e.String(), e.Code)
				} else if e, ok := r.(error); ok {
					resp, err = responder.Fail(req.Headers, e.Error(), http.StatusInternalServerError)
				} else {
					resp, err = responder.Fail(req.Headers, "unknown error", http.StatusInternalServerError)
				}
			}
		}()
		resp, err = handle(ctx, req)
		return
	}
}
