package users

import (
	"context"
	"net/http"

	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-users/types"
	uT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) CreateUser(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	req := &t.SignInReq{}
	api.Parse(r, req)
	user := &uT.User{Email: req.Email}
	user.EncryptedPassword, err = auth.EncryptPassword(req.Password)
	api.CheckError(http.StatusUnprocessableEntity, err)
	usrResp, err := u.usr.CreateUser(ctx, user)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return api.Success(usrResp, http.StatusCreated)
}
