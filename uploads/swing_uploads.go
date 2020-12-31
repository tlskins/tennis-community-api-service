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

// func (u *UploadsService) GetSwingUploadURL(ctx context.Context, userId, fileName string) (resp api.Response, err error) {
// 	uploadId := time.Now().Format("2006_01_02_15_04_05")
// 	sess, err := session.NewSession(&aws.Config{
// 		Region:      aws.String("us-east-1"),
// 		Credentials: credentials.NewStaticCredentials(u.awsConfig.accessKeyId, u.awsConfig.secretAccessKey, ""),
// 	})
// 	api.CheckError(http.StatusInternalServerError, err)

// 	svc := s3.New(sess)
// 	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
// 		Bucket: aws.String(u.awsConfig.bucketName),
// 		Key:    aws.String(fmt.Sprintf("originals/%s/%s/%s", userId, uploadId, fileName)),
// 	})
// 	urlStr, err := req.Presign(15 * time.Minute)
// 	api.CheckError(http.StatusUnprocessableEntity, err)
// 	return api.Success(map[string]string{
// 		"url": urlStr,
// 	}, http.StatusCreated)
// }

// https://tennis-swings.s3.amazonaws.com/originals/b687e24a-6e73-4679-b2cb-2e0aa5e4c109/2020_12_30_01_33_51/test.mp4

func (u *UploadsService) CreateSwingUpload(ctx context.Context, userId, originalURL string) (resp api.Response, err error) {
	now := time.Now()
	paths := strings.Split(originalURL, "/")
	upload := &t.SwingUpload{
		CreatedAt:   now,
		UpdatedAt:   now,
		UploadKey:   paths[len(paths)-2],
		UserID:      userId,
		Status:      enums.SwingUploadStatusOriginal,
		OriginalURL: originalURL,
	}
	newUpload, err := u.Store.CreateSwingUpload(upload)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return api.Success(newUpload, http.StatusCreated)
}

func (u *UploadsService) uploadIDFromFileName(videoPath string) (uploadID, fileName string, err error) {
	// upload id from folder path
	paths := strings.Split(videoPath, "/")
	if len(paths) < 4 {
		return "", "", errors.New("Invalid video path hiearchy format")
	}
	return paths[2], paths[3], nil
}

// S3 Events

// https://tennis-swings.s3.amazonaws.com/clips/timuserid/2020_12_18_1152_59/tim_ground_profile_wide_1min_540p_clip_1.mp4
func (u *UploadsService) CreateUploadClipVideos(_ context.Context, bucket string, outputs []string) (resp *t.SwingUpload, err error) {
	fmt.Printf("%s %v\n", bucket, outputs)
	var uploadID string
	now := time.Now()
	clips := make([]*t.UploadClipVideo, len(outputs))
	for i, videoPath := range outputs {
		var fileName string
		uploadID, fileName, err = u.uploadIDFromFileName(videoPath)
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
	update := &t.UpdateSwingUpload{
		UploadKey:  uploadID,
		UpdatedAt:  now,
		Status:     &status,
		ClipVideos: &clips,
	}
	resp, err = u.Store.UpdateSwingUpload(update)
	return
}

// https://tennis-swings.s3.amazonaws.com/tmp/timuserid/2020_12_18_1152_59/tim_ground_profile_wide_1min_540p_clip_1_swing_1.mp4
func (u *UploadsService) CreateUploadSwingVideos(_ context.Context, bucket string, outputs []string) (upload *t.SwingUpload, swings []*t.UploadSwingVideo, err error) {
	var uploadID string
	now := time.Now()
	swings = make([]*t.UploadSwingVideo, len(outputs))
	for i, videoPath := range outputs {
		var fileName string
		uploadID, fileName, err = u.uploadIDFromFileName(videoPath)
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

// func (u *UploadsService) UploadSwingVideoTranscoded(_ context.Context, uploadID, tranUrl string) (resp *t.SwingUpload, err error) {
// 	rgx := regexp.MustCompile(`clip_(\d{1,})_swing_(\d{1,})..+$`)
// 	matches := rgx.FindStringSubmatch(cutURL)
// 	if len(matches) < 3 {
// 		return nil, errors.New("Invalid clip path format")
// 	}
// 	var id, clipID int
// 	if clipID, err = strconv.Atoi(matches[1]); err != nil {
// 		return
// 	}
// 	if id, err = strconv.Atoi(matches[2]); err != nil {
// 		return
// 	}
// 	resp, err = u.Store.UploadSwingVideoTranscoded(uploadID, tranUrl, swingVideoID)
// 	return
// }
