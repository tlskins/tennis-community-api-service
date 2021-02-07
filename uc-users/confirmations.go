package users

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tennis-community-api-service/pkg/auth"
	"github.com/tennis-community-api-service/pkg/enums"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-users/types"
	uT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) sendWelcomeEmail(ctx context.Context, r *api.Request, email, confirmationID string) (err error) {
	return u.emailClient.SendEmail(
		email,
		"Welcome to Hive Tennis",
		fmt.Sprintf(`
Welcome to Hive Tennis!
Please follow this link to confirm your account:

%s
		`, fmt.Sprintf("%s/?confirmation=%s", u.Resp.Origin(r), confirmationID)),
	)
}

func (u *UCService) CreateUser(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	req := &t.CreateUserReq{}
	api.ParseAndValidate(r, req)

	// create user
	now := time.Now()
	user := &uT.User{
		CreatedAt:  now,
		UpdatedAt:  now,
		UserName:   req.UserName,
		Email:      req.Email,
		LowerEmail: strings.ToLower(req.Email),
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Status:     enums.UserStatusPending,
	}
	user.EncryptedPassword, err = auth.EncryptPassword(req.Password)
	api.CheckError(http.StatusUnprocessableEntity, err)
	usrResp, err := u.usr.CreateUser(ctx, user)
	api.CheckError(http.StatusUnprocessableEntity, err)

	// create confirmation
	conf, err := u.usr.CreateConfirmation(ctx, &uT.UserConfirmation{
		CreatedAt: now,
		UserID:    usrResp.ID,
	})
	api.CheckError(http.StatusUnprocessableEntity, err)

	err = u.sendWelcomeEmail(ctx, r, user.Email, conf.ID)
	api.CheckError(http.StatusUnprocessableEntity, err)

	return u.Resp.Success(r, usrResp, http.StatusCreated)
}

func (u *UCService) InviteUser(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	req := &t.UserInvitationReq{}
	api.Parse(r, req)

	inviter, err := u.usr.GetUser(ctx, req.InviterID)
	api.CheckError(http.StatusNotFound, err)
	if _, err := u.usr.GetUserByEmail(ctx, req.Email); err == nil {
		api.CheckError(http.StatusUnprocessableEntity, fmt.Errorf("User already exists for email %s", req.Email))
	}

	// create confirmation
	conf, err := u.usr.CreateConfirmation(ctx, &uT.UserConfirmation{
		CreatedAt: time.Now(),
		Email:     req.Email,
		InviterID: req.InviterID,
		URL:       req.URL,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	api.CheckError(http.StatusUnprocessableEntity, err)

	err = u.emailClient.SendEmail(
		conf.Email,
		fmt.Sprintf("Hive Tennis - You have been invited to review %s's tennis videos!", inviter.Name()),
		fmt.Sprintf(`
Welcome to Hive Tennis!
%s Has invited you to review their tennis swings on Hive Tennis:

Please follow the link below:
%s

-Hive Tennis
		`, inviter.Name(), fmt.Sprintf("%s/%s?confirmation=%s", u.Resp.Origin(r), conf.URL, conf.ID)),
	)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, nil, http.StatusCreated)
}

func (u *UCService) GetUserConfirmation(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	id := r.PathParameters["id"]
	conf, err := u.usr.GetConfirmation(ctx, id)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, conf, http.StatusOK)
}

func (u *UCService) Confirm(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	req := &t.UserConfirmationReq{}
	api.Parse(r, req)
	conf, err := u.usr.GetConfirmation(ctx, req.ID)
	api.CheckError(http.StatusNotFound, err)

	now := time.Now()
	var newUser *uT.User
	if conf.UserID != "" {
		// normal confirmation
		user, err := u.usr.GetUser(ctx, conf.UserID)
		api.CheckError(http.StatusUnprocessableEntity, err)
		status := enums.UserStatusCreated
		newUser, err = u.usr.UpdateUser(ctx, &uT.UpdateUser{
			ID:        user.ID,
			Status:    &status,
			UpdatedAt: now,
		})
		api.CheckError(http.StatusUnprocessableEntity, err)
	} else if conf.Email != "" {
		inviter, err := u.usr.GetUser(ctx, conf.InviterID)
		api.CheckError(http.StatusUnprocessableEntity, err)
		// create invited user
		user := &uT.User{
			CreatedAt:  now,
			UpdatedAt:  now,
			UserName:   req.UserName,
			Email:      conf.Email,
			LowerEmail: strings.ToLower(conf.Email),
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Status:     enums.UserStatusCreated,
			FriendIds:  []string{inviter.ID},
		}
		user.EncryptedPassword, err = auth.EncryptPassword(req.Password)
		api.CheckError(http.StatusUnprocessableEntity, err)
		newUser, err = u.usr.CreateUser(ctx, user)
		api.CheckError(http.StatusUnprocessableEntity, err)
		inviter.FriendIds = append(inviter.FriendIds, newUser.ID)
		_, err = u.usr.UpdateUser(ctx, &uT.UpdateUser{
			ID:        inviter.ID,
			UpdatedAt: now,
			FriendIds: &inviter.FriendIds,
		})
		api.CheckError(http.StatusUnprocessableEntity, err)
	}

	err = u.usr.DeleteUserConfirmations(ctx, newUser.ID, newUser.Email)
	api.CheckError(http.StatusUnprocessableEntity, err)
	authToken, err := u.jwt.GenAccessToken(newUser)
	api.CheckError(http.StatusInternalServerError, err)
	newUser.AuthToken = authToken

	return u.Resp.Success(r, newUser, http.StatusOK)
}
