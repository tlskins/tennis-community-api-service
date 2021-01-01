package uploads

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tennis-community-api-service/pkg/enums"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uploads/types"

	// "github.com/davecgh/go-spew/spew"
	uuid "github.com/satori/go.uuid"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
	// "github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func (u *UploadsService) GetRecentSwingUploads(ctx context.Context, userId string) (uploads []*t.SwingUpload, err error) {
	return u.Store.GetRecentSwingUploads(userId)
}

// https://tennis-swings.s3.amazonaws.com/originals/b687e24a-6e73-4679-b2cb-2e0aa5e4c109/2020_12_30_01_33_51/test.mp4

func (u *UploadsService) CreateSwingUpload(ctx context.Context, userId, originalURL string) (resp *t.SwingUpload, err error) {
	now := time.Now()
	paths := strings.Split(originalURL, "/")
	return u.Store.CreateSwingUpload(&t.SwingUpload{
		CreatedAt:   now,
		UpdatedAt:   now,
		UploadKey:   paths[len(paths)-2],
		UserID:      userId,
		Status:      enums.SwingUploadStatusOriginal,
		OriginalURL: originalURL,
	})
}

func (u *UploadsService) uploadIDFromFileName(videoPath string) (userID, uploadID, fileName string, err error) {
	// upload id from folder path
	paths := strings.Split(videoPath, "/")
	if len(paths) < 4 {
		return "", "", "", errors.New("Invalid video path hiearchy format")
	}
	return paths[1], paths[2], paths[3], nil
}

// S3 Events

// https://tennis-swings.s3.amazonaws.com/clips/timuserid/2020_12_18_1152_59/tim_ground_profile_wide_1min_540p_clip_1.mp4
func (u *UploadsService) CreateUploadClipVideos(_ context.Context, bucket string, outputs []string) (resp *t.SwingUpload, err error) {
	fmt.Printf("%s %v\n", bucket, outputs)
	var uploadID, userID string
	now := time.Now()
	clips := make([]*t.UploadClipVideo, len(outputs))
	for i, videoPath := range outputs {
		var fileName string
		userID, uploadID, fileName, err = u.uploadIDFromFileName(videoPath)
		api.CheckError(http.StatusInternalServerError, err)
		// clip id from file name
		rgx := regexp.MustCompile(`clip_(\d{1,})..+$`)
		matches := rgx.FindStringSubmatch(fileName)
		if len(matches) < 2 {
			return nil, errors.New("Invalid video name format")
		}
		var id int
		id, err = strconv.Atoi(matches[1])
		api.CheckError(http.StatusInternalServerError, err)
		clips[i] = &t.UploadClipVideo{
			ID:        id,
			CreatedAt: now,
			ClipURL:   fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, videoPath),
		}
	}
	status := enums.SwingUploadStatusClipped
	resp, err = u.Store.UpdateSwingUpload(&t.UpdateSwingUpload{
		UploadKey:  uploadID,
		UserID:     userID,
		UpdatedAt:  now,
		Status:     &status,
		ClipVideos: &clips,
	})
	return
}

// https://tennis-swings.s3.amazonaws.com/tmp/timuserid/2020_12_18_1152_59/tim_ground_profile_wide_1min_540p_clip_1_swing_1.mp4
func (u *UploadsService) CreateUploadSwingVideos(_ context.Context, bucket string, outputs []string) (upload *t.SwingUpload, swings []*t.UploadSwingVideo, err error) {
	var uploadID string
	now := time.Now()
	swings = make([]*t.UploadSwingVideo, len(outputs))
	for i, videoPath := range outputs {
		var fileName string
		_, uploadID, fileName, err = u.uploadIDFromFileName(videoPath)
		api.CheckError(http.StatusInternalServerError, err)
		rgx := regexp.MustCompile(`clip_(\d{1,})_swing_(\d{1,})..+$`)
		matches := rgx.FindStringSubmatch(fileName)
		if len(matches) < 3 {
			return nil, swings, errors.New("Invalid clip path format")
		}
		var swingID, clipID int
		if clipID, err = strconv.Atoi(matches[1]); err != nil {
			return
		}
		if swingID, err = strconv.Atoi(matches[2]); err != nil {
			return
		}
		cutURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, videoPath)
		swings[i] = &t.UploadSwingVideo{
			ID:        uuid.NewV4().String(),
			CreatedAt: now,
			UpdatedAt: now,
			ClipID:    clipID,
			SwingID:   swingID,
			CutURL:    cutURL,
		}
	}

	upload, err = u.Store.CreateUploadSwingVideos(uploadID, swings)
	return upload, swings, err
}

func (u *UploadsService) UpdateSwingUpload(_ context.Context, data *t.UpdateSwingUpload) (upload *t.SwingUpload, err error) {
	return u.Store.UpdateSwingUpload(data)
}
