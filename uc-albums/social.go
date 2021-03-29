package ucalbums

import (
	"context"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/auth"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-albums/types"
	uT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) shareAlbum(ctx context.Context, r *api.Request, album *aT.Album, newFriendIDs []string) (err error) {
	user, err := u.usr.GetUser(ctx, album.UserID)
	if err != nil {
		return err
	}
	if len(newFriendIDs) > 0 {
		err = u.usr.AddFriendNoteToUsers(ctx, newFriendIDs, &uT.FriendNote{
			CreatedAt: time.Now(),
			Subject:   fmt.Sprintf("%s has shared the album %s with you!", user.UserName, album.Name),
			Type:      "Shared Album",
		})
		if err != nil {
			return
		}

		for _, friendID := range newFriendIDs {
			friend, softErr := u.usr.GetUser(ctx, friendID)
			if softErr != nil {
				fmt.Printf("error getting friend: %s\n", softErr.Error())
			}
			softErr = u.emailClient.SendEmail(
				friend.Email,
				fmt.Sprintf("Hive Tennis - %s Shared An Album With You!", user.UserName),
				fmt.Sprintf(`
%s %s,
Your friend %s %s has shared the album "%s" with you.

Link: %s/albums/%s
				`, friend.FirstName, friend.LastName, user.FirstName, user.LastName, album.Name, u.Resp.Origin(r.Headers), album.ID),
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

	// add comment to album in store
	now := time.Now()
	comment := &aT.Comment{
		ReplyID:   req.ReplyID,
		UserID:    req.UserID,
		CreatedAt: now,
		UpdatedAt: now,
		Frame:     req.Frame,
		Text:      req.Text,
		UserTags:  req.UserTags,
	}
	var album *aT.Album
	if req.SwingID == "" {
		album, err = u.alb.PostCommentToAlbum(ctx, req.AlbumID, comment)
	} else {
		album, err = u.alb.PostCommentToSwing(ctx, req.AlbumID, req.SwingID, comment)
	}
	api.CheckError(http.StatusUnprocessableEntity, err)

	// notifications
	poster, err := u.usr.GetUser(ctx, claims.Subject)
	api.CheckError(http.StatusUnprocessableEntity, err)

	// notify tagged users
	for _, tag := range comment.UserTags {
		if tag.UserID == album.UserID || comment.UserID == tag.UserID {
			continue
		}
		tagUser, err := u.usr.GetUser(ctx, tag.UserID)
		api.CheckError(http.StatusUnprocessableEntity, err)
		notes := append(tagUser.AlbumUserTagNotes, &uT.AlbumUserTagNote{
			ID:              uuid.NewV4().String(),
			CreatedAt:       now,
			TaggerID:        poster.ID,
			TaggerFirstName: poster.FirstName,
			TaggerUserName:  poster.UserName,
			AlbumID:         album.ID,
			AlbumName:       album.Name,
			SwingID:         req.SwingID,
		})
		_, err = u.usr.UpdateUser(ctx, &uT.UpdateUser{
			ID:                tagUser.ID,
			AlbumUserTagNotes: &notes,
		})
		api.CheckError(http.StatusUnprocessableEntity, err)

		softErr := u.emailClient.SendEmail(
			tagUser.Email,
			fmt.Sprintf("Hive Tennis - %s Tagged you in the album %s!", poster.UserName, album.Name),
			fmt.Sprintf(`
%s %s,
Your friend %s %s has tagged you in a comment on the album "%s".

Link: %s/albums/%s
	`, tagUser.FirstName, tagUser.LastName, poster.FirstName, poster.LastName, album.Name, u.Resp.Origin(r.Headers), album.ID),
		)
		if softErr != nil {
			fmt.Printf("error sending album user email: %s\n", softErr.Error())
		}
	}

	// notify album user if not commenter
	if album.UserID != "" && album.UserID != claims.Subject {
		albumUser, err := u.usr.GetUser(ctx, album.UserID)
		api.CheckError(http.StatusUnprocessableEntity, err)
		albumUser.AddCommentNote(poster, album.ID, album.Name, req.SwingID)
		albumUser, err = u.usr.UpdateUser(ctx, &uT.UpdateUser{
			ID:           albumUser.ID,
			CommentNotes: &albumUser.CommentNotes,
		})
		api.CheckError(http.StatusUnprocessableEntity, err)

		// record poster's recent comments
		note := poster.AddMyRecentComment(album.ID, album.Name, req.SwingID)
		poster, err = u.usr.UpdateUser(ctx, &uT.UpdateUser{
			ID:               poster.ID,
			MyRecentComments: &poster.MyRecentComments,
		})
		api.CheckError(http.StatusUnprocessableEntity, err)

		// only send email notif if havent notified of their posting recently on this album
		if note.NumComments == 1 {
			softErr := u.emailClient.SendEmail(
				albumUser.Email,
				fmt.Sprintf("Hive Tennis - %s Commented on your album %s!", poster.UserName, album.Name),
				fmt.Sprintf(`
%s %s,
Your friend %s %s has commented on the album "%s".

Link: %s/albums/%s
		`, albumUser.FirstName, albumUser.LastName, poster.FirstName, poster.LastName, album.Name, u.Resp.Origin(r.Headers), album.ID),
			)
			if softErr != nil {
				fmt.Printf("error sending album user email: %s\n", softErr.Error())
			}
		}
	}

	return u.Resp.Success(r.Headers, album, http.StatusOK)
}
