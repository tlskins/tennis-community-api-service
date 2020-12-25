package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
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
const AccessTokenKey = "accessToken"

// RPCAccessTokenKey - cookie and context key for rpc authorization
const RPCAccessTokenKey = "rpcAccessToken"

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

func (j *JWTService) IncludeLambdaAuth(ctx context.Context, req *events.APIGatewayProxyRequest) (context.Context, error) {
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

func ClaimsFromContext(ctx context.Context) (authorized bool, claims *CustomClaims) {
	if ctx.Value(AccessTokenKey) != nil {
		return true, ctx.Value(AccessTokenKey).(*CustomClaims)
	}
	if ctx.Value(RPCAccessTokenKey) != nil {
		return true, nil
	}
	return false, nil
}
