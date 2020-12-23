package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type Request events.APIGatewayProxyRequest

type Response events.APIGatewayProxyResponse

var corsHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
}

func Parse(req *Request, out interface{}) {
	if err := json.Unmarshal([]byte(req.Body), out); err != nil {
		Abort(http.StatusUnprocessableEntity, err)
	}
}

// Fail returns an internal server error with the error message
func Fail(msg string, status int) (Response, error) {
	e := make(map[string]string, 0)
	e["message"] = msg

	// We don't need to worry about this error,
	// as we're controlling the input.
	body, _ := json.Marshal(e)

	return Response{
		Body:       string(body),
		Headers:    corsHeaders,
		StatusCode: status,
	}, nil
}

// Success returns a valid response
func Success(data interface{}, status int) (Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return Fail(err.Error(), http.StatusInternalServerError)
	}

	return Response{
		Body:       string(body),
		Headers:    corsHeaders,
		StatusCode: status,
	}, nil
}

func SuccessWithCookie(data interface{}, status int, cookieStr string) (Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return Fail(err.Error(), http.StatusInternalServerError)
	}

	corsHeaders["Set-Cookie"] = fmt.Sprintf("%s; Secure; SameSite=None; Path=/", cookieStr)

	return Response{
		Body:       string(body),
		StatusCode: status,
		Headers:    corsHeaders,
	}, nil
}

func MatchesRoute(pattern, method string, req *Request) bool {
	if req.HTTPMethod != method {
		return false
	}
	pPaths := strings.Split(pattern, "/")
	rPaths := strings.Split(req.Path, "/")
	if len(pPaths) != len(rPaths) {
		return false
	}
	for i, rPath := range rPaths {
		pPath := pPaths[i]
		isWc := strings.ContainsAny(pPath, "{}")
		if !isWc && pPath != rPath {
			log.Println("matches route false loop")
			return false
		}
	}
	return true
}

func GetPathParam(param string, req *Request) (string, error) {
	res, ok := req.PathParameters[param]
	if !ok {
		return "", fmt.Errorf("Param %s not found", param)
	}
	return res, nil
}

func HandleRequest(handle func(context.Context, *Request) (Response, error)) func(context.Context, *Request) (Response, error) {
	return func(ctx context.Context, req *Request) (resp Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				if e, ok := r.(Error); ok {
					resp, err = Fail(e.String(), e.Code)
				} else if e, ok := r.(error); ok {
					resp, err = Fail(e.Error(), http.StatusInternalServerError)
				} else {
					resp, err = Fail("unknown error", http.StatusInternalServerError)
				}
			}
		}()
		resp, err = handle(ctx, req)
		return
	}
}

// WriteJSON(w, err.Code, map[string]interface{}{"message": err.String()})
// } else if err, ok := r.(error); ok {
// 	WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
// } else {
// 	WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": "unknown error"})
