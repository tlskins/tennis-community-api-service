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
	return nil
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

	albumUser, err := u.usr.GetUser(ctx, album.UserID)
	api.CheckError(http.StatusUnprocessableEntity, err)
	friend, err := u.usr.GetUser(ctx, claims.Subject)
	api.CheckError(http.StatusUnprocessableEntity, err)

	noteFound := false
	notes := albumUser.CommentNotes
	for _, note := range notes {
		if note.AlbumID == req.AlbumID && note.FriendID == friend.ID {
			noteFound = true
			note.NumComments++
			if req.SwingID != "" {
				note.SwingIDs = append(note.SwingIDs, req.SwingID)
			}
			break
		}
	}
	if !noteFound {
		note := &uT.CommentNote{
			ID:              uuid.NewV4().String(),
			CreatedAt:       time.Now(),
			FriendID:        claims.Subject,
			FriendFirstName: friend.FirstName,
			FriendUserName:  friend.UserName,
			AlbumID:         req.AlbumID,
			AlbumName:       album.Name,
			NumComments:     1,
		}
		if req.SwingID != "" {
			note.SwingIDs = []string{req.SwingID}
		}
		notes = append(notes, note)
	}
	albumUser, err = u.usr.UpdateUser(ctx, &uT.UpdateUser{
		ID:           albumUser.ID,
		CommentNotes: &notes,
	})
	api.CheckError(http.StatusUnprocessableEntity, err)

	return u.Resp.Success(album, http.StatusOK)
}
