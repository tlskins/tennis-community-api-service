package uploads

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	// "github.com/davecgh/go-spew/spew"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-albums/types"
	uT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) GetAlbums(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := t.SearchAlbumsReq{UserID: claims.Subject}
	api.Parse(r, &req)

	albumResp := t.AlbumsResp{
		LastRequestAt: time.Now(),
		MyAlbums:      []*aT.Album{},
		FriendsAlbums: []*aT.Album{},
		PublicAlbums:  []*aT.Album{},
	}
	albumResp.MyAlbums, err = u.alb.GetUserAlbums(ctx, claims.Subject)
	api.CheckError(http.StatusInternalServerError, err)
	if !req.ExcludePublic {
		albumResp.PublicAlbums, err = u.alb.GetPublicAlbums(ctx)
		api.CheckError(http.StatusInternalServerError, err)
	}
	if !req.ExcludeFriends {
		albumResp.FriendsAlbums, err = u.alb.GetFriendsAlbums(ctx, claims.Subject)
		api.CheckError(http.StatusInternalServerError, err)
	}
	return u.Resp.Success(albumResp, http.StatusOK)
}

func (u *UCService) CreateAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := t.CreateAlbumReq{UserID: claims.Subject}
	api.ParseAndValidate(r, &req)
	albumReq := aT.Album(req)
	album, err := u.alb.CreateAlbum(ctx, &albumReq)
	api.CheckError(http.StatusInternalServerError, err)

	user, err := u.usr.GetUser(ctx, claims.Subject)
	api.CheckError(http.StatusInternalServerError, err)
	if len(album.FriendIDs) > 0 {
		err = u.usr.AddFriendNoteToUsers(ctx, album.FriendIDs, &uT.FriendNote{
			CreatedAt: time.Now(),
			Subject:   fmt.Sprintf("%s has shared the album %s with you!", user.UserName, album.Name),
			Type:      "Shared Album",
		})
		api.CheckError(http.StatusInternalServerError, err)

		for _, friendID := range album.FriendIDs {
			friend, softErr := u.usr.GetUser(ctx, friendID)
			if softErr != nil {
				fmt.Printf("error getting friend: %s\n", softErr.Error())
			}
			softErr = u.emailClient.SendEmail(
				friend.Email,
				fmt.Sprintf("Tennis Community - %s Shared An Album With You!", user.UserName),
				fmt.Sprintf(`
%s %s,
Your friend %s %s has has shared the album %s with you.
View At
%s/albums/%s
				`, friend.FirstName, friend.LastName, user.FirstName, user.LastName, album.Name, u.Resp.Origin, album.ID),
			)
			if softErr != nil {
				fmt.Printf("error sending friend email: %s\n", softErr.Error())
			}
		}
	}
	return u.Resp.Success(album, http.StatusOK)
}

func (u *UCService) GetAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	id := r.PathParameters["id"]
	album, err := u.alb.GetAlbum(ctx, id)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(album, http.StatusOK)
}

func (u *UCService) UpdateAlbum(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := &aT.UpdateAlbum{ID: r.PathParameters["id"]}
	api.ParseAndValidate(r, req)
	album, err := u.alb.GetAlbum(ctx, req.ID)
	api.CheckError(http.StatusNotFound, err)
	if album.UserID != claims.Subject {
		panic(errors.New("Cannot edit another user's album"))
	}
	album, err = u.alb.UpdateAlbum(ctx, req)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(album, http.StatusOK)
}

func (u *UCService) PostComment(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := &t.PostCommentReq{UserID: claims.Subject, AlbumID: r.PathParameters["id"]}
	api.ParseAndValidate(r, req)
	now := time.Now()

	var album *aT.Album
	if req.SwingID == "" {
		album, err = u.alb.PostCommentToAlbum(ctx, req.AlbumID, &aT.Comment{
			ReplyID:   req.ReplyID,
			UserID:    req.UserID,
			CreatedAt: now,
			UpdatedAt: now,
			Frame:     req.Frame,
			Text:      req.Text,
		})
	} else {
		album, err = u.alb.PostCommentToSwing(ctx, req.AlbumID, req.SwingID, &aT.Comment{
			ReplyID:   req.ReplyID,
			UserID:    req.UserID,
			CreatedAt: now,
			UpdatedAt: now,
			Frame:     req.Frame,
			Text:      req.Text,
		})
	}
	api.CheckError(http.StatusUnprocessableEntity, err)
	return u.Resp.Success(album, http.StatusOK)
}
