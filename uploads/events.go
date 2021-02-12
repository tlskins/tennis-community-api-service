package uploads

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/tennis-community-api-service/pkg/enums"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uploads/types"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (u *UploadsService) CreateUploadClipVideos(_ context.Context, bucket string, outputs []string) (resp *t.SwingUpload, err error) {
	var uploadID, userID string
	now := time.Now()
	clips := make([]*t.UploadClipVideo, len(outputs))
	for i, videoPath := range outputs {
		var fileName string
		fmt.Printf("videoPath %s\n", videoPath)
		paths := strings.Split(videoPath, "/")
		userID = paths[len(paths)-4]
		uploadID = paths[len(paths)-3]
		fileName = paths[len(paths)-1]
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

func (u *UploadsService) CreateUploadSwingVideos(_ context.Context, bucket string, videos, gifs, jpgs, txts []string) (upload *t.SwingUpload, swings []*t.UploadSwingVideo, err error) {
	var uploadID string
	now := time.Now()
	swings = make([]*t.UploadSwingVideo, len(videos))
	for i, videoPath := range videos {
		meta := &t.SwingUploadMeta{}
		metaPath := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, txts[i])
		fmt.Printf("metaURL = %s\n", metaPath)
		if err = u.unmarshalJSONFile(metaPath, meta); err != nil {
			return
		}
		if uploadID == "" {
			uploadID = meta.UploadKey
		}
		spew.Dump(meta)
		swings[i] = &t.UploadSwingVideo{
			ID:               uuid.NewV4().String(),
			CreatedAt:        now,
			UpdatedAt:        now,
			TimestampSeconds: meta.TimestampSeconds,
			Frames:           meta.Frames,
			ClipID:           meta.Clip,
			SwingID:          meta.Swing,
			CutURL:           fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, videoPath),
			GifURL:           fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, gifs[i]),
			JpgURL:           fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, jpgs[i]),
		}
	}

	upload, err = u.Store.CreateUploadSwingVideos(uploadID, swings)
	fmt.Printf("after store create upload\n")
	return upload, swings, err
}

func (u *UploadsService) unmarshalJSONFile(txtURL string, out interface{}) (err error) {
	// download file
	res, err := http.Get(txtURL)
	if err != nil {
		return err
	}
	bytesData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()

	return json.Unmarshal(bytesData, out)
}
