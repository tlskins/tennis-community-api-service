package uploads

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/tennis-community-api-service/pkg/enums"
	t "github.com/tennis-community-api-service/uploads/types"

	uuid "github.com/satori/go.uuid"
)

func (u *UploadsService) CreateUploadClipVideos(_ context.Context, uploadID, userID string, clips []*t.UploadClipVideo) (resp *t.SwingUpload, err error) {
	status := enums.SwingUploadStatusClipped
	resp, err = u.Store.UpdateSwingUpload(&t.UpdateSwingUpload{
		UploadKey:  uploadID,
		UserID:     userID,
		UpdatedAt:  time.Now(),
		Status:     &status,
		ClipVideos: &clips,
	})
	return
}

func (u *UploadsService) CreateUploadSwingVideos(_ context.Context, bucket string, videos, gifs, jpgs, txts []string) (upload *t.SwingUpload, swings []*t.UploadSwingVideo, err error) {
	var uploadID string
	var clipNum int
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
		clipNum = meta.Clip
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

	upload, err = u.Store.CreateUploadSwingVideos(uploadID, clipNum, swings)
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
