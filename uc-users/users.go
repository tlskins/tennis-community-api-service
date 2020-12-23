package users

import (
	"context"
	"net/http"
	"time"

	"github.com/tennis-community-api-service/pkg/auth"
	"github.com/tennis-community-api-service/pkg/enums"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-users/types"
	uT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) CreateUser(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	req := &t.CreateUserReq{}
	api.ParseAndValidate(r, req)

	now := time.Now()
	user := &uT.User{
		Email:     req.Email,
		Status:    enums.UserStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	user.EncryptedPassword, err = auth.EncryptPassword(req.Password)
	api.CheckError(http.StatusUnprocessableEntity, err)

	usrResp, err := u.usr.CreateUser(ctx, user)
	api.CheckError(http.StatusUnprocessableEntity, err)

	return api.Success(usrResp, http.StatusCreated)
}
