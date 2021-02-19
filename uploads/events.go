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

func (u *UploadsService) CreateUploadSwingVideos(_ context.Context, bucket, userID, uploadID string, clipNum int, videos, gifs, jpgs, txts []string) (upload *t.SwingUpload, swings []*t.UploadSwingVideo, err error) {
	now := time.Now()
	swings = make([]*t.UploadSwingVideo, len(videos))
	for i, videoPath := range videos {
		meta := &t.SwingUploadMeta{}
		metaPath := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, txts[i])
		fmt.Printf("metaURL = %s\n", metaPath)
		if err = u.unmarshalJSONFile(metaPath, meta); err != nil {
			return
		}
		spew.Dump(meta)

		cutURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, videoPath)
		gifURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, gifs[i])
		jpgURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, jpgs[i])
		if u.cdnURL != "" {
			cutURL = fmt.Sprintf("%s/%s", u.cdnURL, videoPath)
			gifURL = fmt.Sprintf("%s/%s", u.cdnURL, gifs[i])
			jpgURL = fmt.Sprintf("%s/%s", u.cdnURL, jpgs[i])
		}
		swings[i] = &t.UploadSwingVideo{
			ID:               uuid.NewV4().String(),
			CreatedAt:        now,
			UpdatedAt:        now,
			TimestampSeconds: meta.TimestampSeconds,
			Frames:           meta.Frames,
			ClipID:           clipNum,
			SwingID:          meta.Swing,
			CutURL:           cutURL,
			GifURL:           gifURL,
			JpgURL:           jpgURL,
		}
	}

	upload, err = u.Store.CreateUploadSwingVideos(userID, uploadID, clipNum, swings)
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
