package uploads

import (
	"context"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/auth"
	"github.com/tennis-community-api-service/pkg/enums"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uc-uploads/types"
	uT "github.com/tennis-community-api-service/uploads/types"
	usrT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) GetRecentSwingUploads(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	uploads, err := u.up.GetRecentSwingUploads(ctx, claims.Subject)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(r, uploads, http.StatusCreated)
}

func (u *UCService) CreateSwingUpload(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	ctx, err = u.jwt.IncludeLambdaAuth(ctx, r)
	api.CheckError(http.StatusInternalServerError, err)
	req := &t.CreateSwingUploadReq{}
	api.ParseAndValidate(r, req)
	claims := auth.AuthorizedClaimsFromContext(ctx)
	upload, err := u.up.CreateSwingUpload(
		ctx,
		claims.Subject,
		req.OriginalURL,
		req.AlbumName,
		req.FriendIDs,
		req.IsPublic,
		req.IsViewableByFriends,
	)
	api.CheckError(http.StatusInternalServerError, err)
	return u.Resp.Success(r, upload, http.StatusCreated)
}

func (u *UCService) CreateUploadClipVideos(ctx context.Context, r *t.UploadClipEvent) (string, error) {
	upload, err := u.up.CreateUploadClipVideos(ctx, r.ResponsePayload.Body.Bucket, r.ResponsePayload.Body.Outputs)
	if err != nil {
		return "error", err
	}
	album, err := u.alb.CreateAlbumFromUpload(
		ctx,
		upload.UserID,
		upload.UploadKey,
		upload.AlbumName,
		upload.IsPublic,
		upload.IsViewableByFriends,
		upload.FriendIDs,
	)
	if err != nil {
		return "error", err
	}
	_, err = u.up.UpdateSwingUpload(ctx, &uT.UpdateSwingUpload{
		UploadKey: upload.UploadKey,
		UserID:    upload.UserID,
		UpdatedAt: time.Now(),
		AlbumID:   &album.ID,
	})
	if err != nil {
		return "error", err
	}
	return "success", nil
}

func (u *UCService) CreateUploadSwingVideos(ctx context.Context, r *t.UploadSwingEvent) (string, error) {
	videos, gifs, jpgs := r.Outputs()
	upload, swingUploads, err := u.up.CreateUploadSwingVideos(ctx, r.ResponsePayload.Body.Bucket, videos, gifs, jpgs)
	if err != nil {
		return "error CreateUploadSwingVideos", err
	}

	now := time.Now()
	swingVids := make([]*aT.SwingVideo, len(swingUploads))
	album, err := u.alb.AddVideosToAlbum(ctx, upload.UserID, upload.UploadKey, swingVids)
	if err != nil {
		return "error AddVideosToAlbum", err
	}

	// upload is finished
	if upload.IsFinal() {
		aStatus := enums.AlbumStatusCreated
		album, err := u.alb.UpdateAlbum(ctx, &aT.UpdateAlbum{
			ID:     album.ID,
			Status: &aStatus,
		})
		if err != nil {
			return "error UpdateAlbum", err
		}

		uStatus := enums.SwingUploadStatusFinished
		_, err = u.up.UpdateSwingUpload(ctx, &uT.UpdateSwingUpload{
			UploadKey: upload.UploadKey,
			UserID:    upload.UserID,
			Status:    &uStatus,
		})
		if err != nil {
			return "error UpdateSwingUpload", err
		}

		// notify user upload finished
		_, err = u.usr.AddUploadNotifications(ctx, upload.UserID, &usrT.UploadNote{
			ID:        uuid.NewV4().String(),
			CreatedAt: now,
			Subject:   fmt.Sprintf("Upload %s has finished processing", upload.UploadKey),
			Type:      "Upload Complete",
			UploadID:  upload.ID,
			UploadKey: upload.UploadKey,
			AlbumID:   album.ID,
			UploadAt:  upload.CreatedAt,
		})
		if err != nil {
			return "error AddUploadNotifications", err
		}

		// notify friends of album shared with them
		if len(album.FriendIDs) > 0 {
			user, err := u.usr.GetUser(ctx, album.UserID)
			if err != nil {
				return "error get user", err
			}

			err = u.usr.AddFriendNoteToUsers(ctx, album.FriendIDs, &usrT.FriendNote{
				CreatedAt: time.Now(),
				Subject:   fmt.Sprintf("%s has shared the album %s with you!", user.UserName, album.Name),
				Type:      "Shared Album",
			})
			if err != nil {
				fmt.Printf("error AddFriendNoteToUsers: %s\n", err.Error())
			}

			for _, friendID := range album.FriendIDs {
				friend, softErr := u.usr.GetUser(ctx, friendID)
				if softErr != nil {
					fmt.Printf("error getting friend: %s\n", softErr.Error())
				}
				softErr = u.emailClient.SendEmail(
					friend.Email,
					fmt.Sprintf("Tennis Community - %s Shared An Album With You!", user.UserName),
					fmt.Sprintf(`%s %s,
Your friend %s %s has has shared the album %s with you.
View At:
%s/albums/%s
					`, friend.FirstName, friend.LastName, user.FirstName, user.LastName, album.Name, u.Resp.Origin, album.ID),
				)
				if softErr != nil {
					fmt.Printf("error sending friend email: %s\n", softErr.Error())
				}
			}
		}
	}
	return "success", nil
}
