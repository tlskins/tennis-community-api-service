package uploads

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uploads/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func (u *UploadsService) GetSwingUploadURL(ctx context.Context, userId, fileName string) (resp api.Response, err error) {
	uploadId := time.Now().Format("2006_01_02_15_04_05")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(u.awsConfig.accessKeyId, u.awsConfig.secretAccessKey, ""),
	})
	api.CheckError(http.StatusInternalServerError, err)

	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(u.awsConfig.bucketName),
		Key:    aws.String(fmt.Sprintf("originals/%s/%s/%s", userId, uploadId, fileName)),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	api.CheckError(http.StatusUnprocessableEntity, err)
	return api.Success(map[string]string{
		"url": urlStr,
	}, http.StatusCreated)
}

func (u *UploadsService) CreateSwingUpload(ctx context.Context, r *api.Request) (resp api.Response, err error) {
	spew.Dump(r)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(u.awsConfig.accessKeyId, u.awsConfig.secretAccessKey, ""),
	})
	api.CheckError(http.StatusInternalServerError, err)

	c, err := sess.Config.Credentials.Get()
	spew.Dump(c)
	api.CheckError(http.StatusInternalServerError, err)

	reader := strings.NewReader(r.Body)
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.awsConfig.bucketName),
		Key:    aws.String("test_upload.mp4"),
		Body:   reader,
	})
	api.CheckError(http.StatusUnprocessableEntity, err)
	return api.Success(map[string]string{
		"result": "success",
	}, http.StatusCreated)
}

// S3 Events

// https://tennis-swings.s3.amazonaws.com/clips/timuserid/2020_12_18_1152_59/tim_ground_profile_wide_1min_540p_clip_1.mp4
func (u *UploadsService) CreateUploadClipVideos(_ context.Context, bucket string, outputs []string) (resp *t.SwingUpload, err error) {
	var uploadID string
	now := time.Now()
	clips := make([]*t.UploadClipVideo, len(outputs))

	for i, videoPath := range outputs {
		// upload id from folder path
		paths := strings.Split(videoPath, "/")
		if len(paths) < 4 {
			return nil, errors.New("Invalid video path hiearchy format")
		}
		uploadID = paths[2]
		// clip id from file name
		rgx := regexp.MustCompile(`(\d{1,})..+$`)
		matches := rgx.FindStringSubmatch(paths[3])
		if len(matches) < 2 {
			return nil, errors.New("Invalid video name format")
		}
		var id int
		if id, err = strconv.Atoi(matches[1]); err != nil {
			return
		}
		clips[i] = &t.UploadClipVideo{
			ID:        id,
			CreatedAt: now,
			ClipURL:   fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, videoPath),
		}
	}

	resp, err = u.Store.CreateUploadClipVideos(uploadID, clips)
	return
}

// https://tennis-swings.s3.amazonaws.com/tmp/timuserid/2020_12_18_1152_59/tim_ground_profile_wide_1min_540p_clip_1_swing_1.mp4
func (u *UploadsService) CreateUploadSwingVideos(_ context.Context, bucket string, outputs []string) (resp *t.SwingUpload, err error) {
	var uploadID string
	now := time.Now()
	swings := make([]*t.UploadSwingVideo, len(outputs))

	for i, videoPath := range outputs {
		// upload id from folder path
		paths := strings.Split(videoPath, "/")
		if len(paths) < 4 {
			return nil, errors.New("Invalid video path hiearchy format")
		}
		uploadID = paths[2]
		// clipe id and swing id from file name
		fileName := paths[3]
		rgx := regexp.MustCompile(`clip_(\d{1,})_swing_(\d{1,})..+$`)
		matches := rgx.FindStringSubmatch(fileName)
		if len(matches) < 3 {
			return nil, errors.New("Invalid clip path format")
		}
		var id, clipID int
		if clipID, err = strconv.Atoi(matches[1]); err != nil {
			return
		}
		if id, err = strconv.Atoi(matches[2]); err != nil {
			return
		}
		cutURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, videoPath)
		swings[i] = &t.UploadSwingVideo{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
			ClipID:    clipID,
			CutURL:    cutURL,
		}
	}

	resp, err = u.Store.CreateUploadSwingVideos(uploadID, swings)
	return
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
