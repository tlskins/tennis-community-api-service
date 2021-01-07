package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	api "github.com/tennis-community-api-service/pkg/lambda"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	key                      *rsa.PublicKey
	secret                   *rsa.PrivateKey
	accessExpirationMinutes  int
	refreshExpirationMinutes int
}

type JWTServiceConfig struct {
	Key                      []byte
	Secret                   []byte
	AccessExpirationMinutes  int
	RefreshExpirationMinutes int
}

// AccessTokenKey - cookie and context key for jwt claims authorization
const AccessTokenKey = "authToken"

func NewJWTService(config JWTServiceConfig) (*JWTService, error) {
	var jKey *rsa.PublicKey
	var err error
	if config.Key != nil {
		if jKey, err = jwt.ParseRSAPublicKeyFromPEM(config.Key); err != nil {
			return nil, err
		}
	}
	var jSecret *rsa.PrivateKey
	if config.Secret != nil {
		if jSecret, err = jwt.ParseRSAPrivateKeyFromPEM(config.Secret); err != nil {
			return nil, err
		}
	}

	return &JWTService{
		key:                      jKey,
		secret:                   jSecret,
		accessExpirationMinutes:  config.AccessExpirationMinutes,
		refreshExpirationMinutes: config.RefreshExpirationMinutes,
	}, nil
}

func (j *JWTService) AccessExpirationMinutes() int {
	return j.accessExpirationMinutes
}

func (j *JWTService) RefreshExpirationMinutes() int {
	return j.refreshExpirationMinutes
}

func (j *JWTService) IncludeLambdaAuth(ctx context.Context, req *api.Request) (context.Context, error) {
	authVal := ""
	if val, ok := req.Headers["Cookie"]; ok {
		authVal = strings.ReplaceAll(val, fmt.Sprintf("%s=", AccessTokenKey), "")
	} else if val, ok := req.Headers["Authorization"]; ok {
		authVal = strings.TrimPrefix(val, "Bearer ")
	}
	if authVal != "" {
		var err error
		var claims *CustomClaims
		if claims, err = j.Decode(authVal); err != nil {
			return ctx, err
		}
		ctx = context.WithValue(ctx, AccessTokenKey, claims)
	}
	return ctx, nil
}

func (j *JWTService) IncludeLambdaWSAuth(ctx context.Context, req *api.WebsocketRequest) (context.Context, error) {
	authVal := req.QueryStringParameters["Authorization"]
	var err error
	var claims *CustomClaims
	if claims, err = j.Decode(authVal); err != nil {
		return ctx, err
	}
	ctx = context.WithValue(ctx, AccessTokenKey, claims)
	return ctx, nil
}

func ClaimsFromContext(ctx context.Context) (authorized bool, claims *CustomClaims) {
	if ctx.Value(AccessTokenKey) != nil {
		return true, ctx.Value(AccessTokenKey).(*CustomClaims)
	}
	return false, nil
}

func AuthorizedClaimsFromContext(ctx context.Context) *CustomClaims {
	authorized, claims := ClaimsFromContext(ctx)
	if !authorized {
		panic(errors.New("Unauthorized"))
	}
	return claims
}
