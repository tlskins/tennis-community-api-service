package uploads

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	uuid "github.com/satori/go.uuid"

	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/enums"
	t "github.com/tennis-community-api-service/uc-uploads/types"
	uT "github.com/tennis-community-api-service/uploads/types"
	usrT "github.com/tennis-community-api-service/users/types"
)

func (u *UCService) CreateUploadClipVideos(ctx context.Context, r *t.UploadClipEvent) (string, error) {
	now := time.Now()
	body := r.ResponsePayload.Body
	clips := make([]*uT.UploadClipVideo, len(body.Outputs))
	for i, clipMeta := range body.Outputs {
		clips[i] = &uT.UploadClipVideo{
			ID:           clipMeta.Number,
			CreatedAt:    now,
			ClipURL:      fmt.Sprintf("https://%s.s3.amazonaws.com/%s", r.ResponsePayload.Body.Bucket, clipMeta.Path),
			FileName:     clipMeta.FileName,
			StartSeconds: clipMeta.StartSeconds,
			EndSeconds:   clipMeta.EndSeconds,
		}
	}

	upload, err := u.up.CreateUploadClipVideos(ctx, body.UploadID, body.UserID, clips)
	if err != nil {
		return "error", err
	}
	fmt.Printf("after CreateUploadClipVideos\n")
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
	fmt.Printf("after CreateAlbumFromUpload\n")
	_, err = u.up.UpdateSwingUpload(ctx, &uT.UpdateSwingUpload{
		UploadKey: upload.UploadKey,
		UserID:    upload.UserID,
		UpdatedAt: time.Now(),
		AlbumID:   &album.ID,
	})
	if err != nil {
		return "error", err
	}
	fmt.Printf("after UpdateSwingUpload\n")
	return "success", nil
}

func (u *UCService) CreateUploadSwingVideos(ctx context.Context, r *t.UploadSwingEvent) (string, error) {
	spew.Dump(r)
	videos, gifs, jpgs, txts := r.Outputs()
	upload, swingUploads, err := u.up.CreateUploadSwingVideos(ctx, r.ResponsePayload.Body.Bucket, videos, gifs, jpgs, txts)
	spew.Dump(upload)
	if err != nil {
		return "error CreateUploadSwingVideos", err
	}

	now := time.Now()
	swingVids := make([]*aT.SwingVideo, len(swingUploads))
	for i, swing := range swingUploads {
		swingVids[i] = &aT.SwingVideo{
			CreatedAt:        now,
			UpdatedAt:        now,
			UserID:           upload.UserID,
			UploadKey:        upload.UploadKey,
			Clip:             swing.ClipID,
			Swing:            swing.SwingID,
			TimestampSeconds: swing.TimestampSeconds,
			Frames:           swing.Frames,
			VideoURL:         swing.CutURL,
			GifURL:           swing.GifURL,
			JpgURL:           swing.JpgURL,
			Status:           enums.SwingVideoStatusCreated,
		}
	}
	album, err := u.alb.AddVideosToAlbum(ctx, upload.UserID, upload.UploadKey, swingVids)
	if err != nil {
		return "error AddVideosToAlbum", err
	}

	// upload is finished
	if upload.IsFinal() {
		aStatus := enums.AlbumStatusCreated
		for i, swing := range album.SwingVideos {
			swing.Name = strconv.Itoa(i + 1)
		}
		album.CalculateMetrics()
		album, err := u.alb.UpdateAlbum(ctx, &aT.UpdateAlbum{
			ID:          album.ID,
			Status:      &aStatus,
			SwingVideos: &album.SwingVideos,
		})
		if err != nil {
			return "error UpdateAlbum", err
		}

		uStatus := enums.SwingUploadStatusFinished
		_, err = u.up.UpdateSwingUpload(ctx, &uT.UpdateSwingUpload{
			UpdatedAt: now,
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
					continue
				}
				softErr = u.emailClient.SendEmail(
					friend.Email,
					fmt.Sprintf("Hive Tennis - %s Shared An Album With You!", user.UserName),
					fmt.Sprintf(`%s %s,
Your friend %s %s has has shared the album %s with you.
View At:
%s/albums/%s
					`, friend.FirstName, friend.LastName, user.FirstName, user.LastName, album.Name, u.Resp.Origin(r.ResponsePayload.Headers), album.ID),
				)
				if softErr != nil {
					fmt.Printf("error sending friend email: %s\n", softErr.Error())
				}
			}
		}
	}
	return "success", nil
}
