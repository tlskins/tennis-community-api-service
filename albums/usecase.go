package albums

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	t "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/enums"

	"github.com/pkg/errors"
)

type store interface {
	CreateSwingUpload(ctx context.Context, data *t.SwingUpload) (*t.SwingUpload, error)
	CreateUploadClipVideos(ctx context.Context, uploadID string, clips []*t.UploadClipVideo) (*t.SwingUpload, error)
	CreateUploadSwingVideos(ctx context.Context, uploadID string, swings []*t.UploadSwingVideo) (*t.SwingUpload, error)
	UpdateUploadSwingVideo(ctx context.Context, uploadID string, update *t.UpdateUploadSwingVideo) (*t.SwingUpload, error)
}

type Usecase struct {
	Store store
}

func (u *Usecase) CreateSwingUpload(ctx context.Context, origURL, userID string) (resp *t.SwingUpload, err error) {
	data := &t.SwingUpload{
		OriginalURL: origURL,
		UserID:      userID,
		Status:      enums.SwingUploadStatusOriginal,
	}
	resp, err = u.Store.CreateSwingUpload(ctx, data)
	return
}

// https://tennis-swings.s3.amazonaws.com/clips/timuserid/2020_12_18_1152_59/tim_ground_profile_wide_1min_540p_clip_1.mp4
func (u *Usecase) CreateUploadClipVideos(ctx context.Context, bucket string, outputs []string) (resp *t.SwingUpload, err error) {
	var uploadID string
	now := time.Now()
	clips := make([]*t.UploadClipVideo, len(outputs))

	for i, videoPath := range outputs {
		paths := strings.Split(videoPath, "/")
		if len(paths) < 4 {
			return nil, errors.New("Invalid video path hiearchy format")
		}
		uploadID = paths[2]

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

	resp, err = u.Store.CreateUploadClipVideos(ctx, uploadID, clips)
	return
}

// https://tennis-swings.s3.amazonaws.com/tmp/timuserid/2020_12_18_1152_59/tim_ground_profile_wide_1min_540p_clip_1_swing_1.mp4
func (u *Usecase) CreateUploadSwingVideos(ctx context.Context, bucket string, outputs []string) (resp *t.SwingUpload, err error) {
	var uploadID string
	now := time.Now()
	swings := make([]*t.UploadSwingVideo, len(outputs))

	for i, videoPath := range outputs {
		paths := strings.Split(videoPath, "/")
		if len(paths) < 4 {
			return nil, errors.New("Invalid video path hiearchy format")
		}
		uploadID := paths[2]

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
			ID:     id,
			ClipID: clipID,
			CutURL: cutURL,
		}
	}

	resp, err = u.Store.CreateUploadSwingVideos(ctx, uploadID, swings)
	return
}

// func (u *Usecase) UploadSwingVideoTranscoded(ctx context.Context, uploadID, tranUrl string) (resp *t.SwingUpload, err error) {
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
// 	resp, err = u.Store.UploadSwingVideoTranscoded(ctx, uploadID, tranUrl, swingVideoID)
// 	return
// }
