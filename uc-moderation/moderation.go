package ucalbums

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	mT "github.com/tennis-community-api-service/moderation/types"
	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-moderation/types"
)

func (u *UCService) CreateCommentFlag(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	now := time.Now()
	req := mT.CommentFlag{
		FlaggerID: claims.Subject,
		CreatedAt: now,
		UpdatedAt: now,
	}
	api.Parse(r, &req)
	flag, err := u.mod.CreateCommentFlag(ctx, &req)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(r, flag, http.StatusOK)
}

func (u *UCService) CreateAlbumFlag(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	now := time.Now()
	req := mT.AlbumFlag{
		FlaggerID: claims.Subject,
		CreatedAt: now,
		UpdatedAt: now,
	}
	api.Parse(r, &req)
	flag, err := u.mod.CreateAlbumFlag(ctx, &req)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(r, flag, http.StatusOK)
}

func (u *UCService) UpdateCommentFlag(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	if !claims.IsAdmin {
		api.CheckError(http.StatusUnauthorized, errors.New("Must be admin"))
	}
	req := &mT.UpdateCommentFlag{ID: r.PathParameters["id"], UpdatedAt: time.Now()}
	api.Parse(r, req)
	flag, err := u.mod.UpdateCommentFlag(ctx, req)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, flag, http.StatusOK)
}

func (u *UCService) UpdateAlbumFlag(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	if !claims.IsAdmin {
		api.CheckError(http.StatusUnauthorized, errors.New("Must be admin"))
	}
	req := &mT.UpdateAlbumFlag{ID: r.PathParameters["id"], UpdatedAt: time.Now()}
	api.Parse(r, req)
	flag, err := u.mod.UpdateAlbumFlag(ctx, req)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, flag, http.StatusOK)
}

func (u *UCService) RecentFlaggedComments(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	if !claims.IsAdmin {
		api.CheckError(http.StatusUnauthorized, errors.New("Must be admin"))
	}
	req := &t.RecentFlaggedCommentsReq{}
	api.Parse(r, req)
	limit, err := strconv.Atoi(req.Limit)
	api.CheckError(http.StatusUnprocessableEntity, err)
	offset, err := strconv.Atoi(req.Offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	var resolved *bool
	if req.Resolved != nil {
		res := *req.Resolved == "true"
		resolved = &res
	}
	flags, err := u.mod.RecentFlaggedComments(ctx, req.Start, req.End, resolved, limit, offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, flags, http.StatusOK)
}

func (u *UCService) RecentFlaggedAlbums(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	if !claims.IsAdmin {
		api.CheckError(http.StatusUnauthorized, errors.New("Must be admin"))
	}
	req := &t.RecentFlaggedAlbumsReq{}
	api.Parse(r, req)
	limit, err := strconv.Atoi(req.Limit)
	api.CheckError(http.StatusUnprocessableEntity, err)
	offset, err := strconv.Atoi(req.Offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	var resolved *bool
	if req.Resolved != nil {
		res := *req.Resolved == "true"
		resolved = &res
	}
	flags, err := u.mod.RecentFlaggedAlbums(ctx, req.Start, req.End, resolved, limit, offset)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(r, flags, http.StatusOK)
}
