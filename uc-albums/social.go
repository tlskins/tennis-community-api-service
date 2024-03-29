package ucalbums

import (
	"context"
	"fmt"
	"net/http"
	"time"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-albums/types"
	uT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) shareAlbum(ctx context.Context, album *aT.Album) (err error) {
	user, err := u.usr.GetUser(ctx, album.UserID)
	if err != nil {
		return err
	}
	if len(album.FriendIDs) > 0 {
		err = u.usr.AddFriendNoteToUsers(ctx, album.FriendIDs, &uT.FriendNote{
			CreatedAt: time.Now(),
			Subject:   fmt.Sprintf("%s has shared the album %s with you!", user.UserName, album.Name),
			Type:      "Shared Album",
		})
		if err != nil {
			return
		}

		for _, friendID := range album.FriendIDs {
			friend, softErr := u.usr.GetUser(ctx, friendID)
			if softErr != nil {
				fmt.Printf("error getting friend: %s\n", softErr.Error())
			}
			softErr = u.emailClient.SendEmail(
				friend.Email,
				fmt.Sprintf("Hive Tennis - %s Shared An Album With You!", user.UserName),
				fmt.Sprintf(`
%s %s,
Your friend %s %s has shared the album %s with you.
View At
%s/albums/%s
				`, friend.FirstName, friend.LastName, user.FirstName, user.LastName, album.Name, u.Resp.Origin, album.ID),
			)
			if softErr != nil {
				fmt.Printf("error sending friend email: %s\n", softErr.Error())
			}
		}
	}
	return nil
}

func (u *UCService) PostComment(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	req := &t.PostCommentReq{UserID: claims.Subject, AlbumID: r.PathParameters["id"]}
	api.ParseAndValidate(r, req)
	now := time.Now()

	// add comment to album in store
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

	// only send notifications if post to another user's album
	if album.UserID != claims.Subject {
		// add comment notification to user in store
		albumUser, err := u.usr.GetUser(ctx, album.UserID)
		api.CheckError(http.StatusUnprocessableEntity, err)
		friend, err := u.usr.GetUser(ctx, claims.Subject)
		api.CheckError(http.StatusUnprocessableEntity, err)
		albumUser.AddCommentNote(friend, album.ID, album.Name, req.SwingID)
		albumUser, err = u.usr.UpdateUser(ctx, &uT.UpdateUser{
			ID:           albumUser.ID,
			CommentNotes: &albumUser.CommentNotes,
		})
		api.CheckError(http.StatusUnprocessableEntity, err)

		// send email if new comment on album
		note := friend.AddMyRecentComment(album.ID, album.Name, req.SwingID)
		albumUser, err = u.usr.UpdateUser(ctx, &uT.UpdateUser{
			ID:               friend.ID,
			MyRecentComments: &friend.MyRecentComments,
		})
		api.CheckError(http.StatusUnprocessableEntity, err)
		if note.NumComments == 1 {
			softErr := u.emailClient.SendEmail(
				friend.Email,
				fmt.Sprintf("Hive Tennis - %s Commented on your album %s!", friend.UserName, album.Name),
				fmt.Sprintf(`
%s %s,
Your friend %s %s has commented on the album %s.
View At
%s/albums/%s
		`, albumUser.FirstName, albumUser.LastName, friend.FirstName, friend.LastName, album.Name, u.Resp.Origin, album.ID),
			)
			if softErr != nil {
				fmt.Printf("error sending album user email: %s\n", softErr.Error())
			}
		}
	}

	return u.Resp.Success(r, album, http.StatusOK)
}
