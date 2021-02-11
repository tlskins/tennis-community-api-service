package uploads

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tennis-community-api-service/pkg/enums"
	api "github.com/tennis-community-api-service/pkg/lambda"
	t "github.com/tennis-community-api-service/uploads/types"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (u *UploadsService) GetRecentSwingUploads(ctx context.Context, userId string) (uploads []*t.SwingUpload, err error) {
	return u.Store.GetRecentSwingUploads(userId)
}

func (u *UploadsService) CreateSwingUpload(ctx context.Context, userId, originalURL, albumName string, friendIds []string, isPublic, isFriends bool) (resp *t.SwingUpload, err error) {
	now := time.Now()
	paths := strings.Split(originalURL, "/")
	return u.Store.CreateSwingUpload(&t.SwingUpload{
		CreatedAt:           now,
		UpdatedAt:           now,
		UploadKey:           paths[len(paths)-2],
		UserID:              userId,
		Status:              enums.SwingUploadStatusOriginal,
		OriginalURL:         originalURL,
		AlbumName:           albumName,
		IsPublic:            isPublic,
		IsViewableByFriends: isFriends,
		FriendIDs:           friendIds,
	})
}

func (u *UploadsService) uploadIDFromClippedFileName(videoPath string) (userID, uploadID, fileLen, fileName string, err error) {
	// upload id from folder path
	paths := strings.Split(videoPath, "/")
	if len(paths) < 4 {
		return "", "", "", "", errors.New("Invalid video path hiearchy format")
	}
	return paths[len(paths)-4], paths[len(paths)-3], paths[len(paths)-2], paths[len(paths)-1], nil
}

// S3 Events

func (u *UploadsService) CreateUploadClipVideos(_ context.Context, bucket string, outputs []string) (resp *t.SwingUpload, err error) {
	var uploadID, userID string
	now := time.Now()
	clips := make([]*t.UploadClipVideo, len(outputs))
	for i, videoPath := range outputs {
		var fileName string
		fmt.Printf("videoPath %s\n", videoPath)
		userID, uploadID, _, fileName, err = u.uploadIDFromClippedFileName(videoPath)
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
		if uploadID == "" {
			// video path is swing output s3 key
			videoPaths := strings.Split(videoPath, "/")
			uploadID = videoPaths[len(videoPaths)-2]
		}
		var meta *t.SwingUploadMeta
		fmt.Printf("txtURL = %s\n", fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, txts[i]))
		if meta, err = u.parseMetaFile(fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, txts[i])); err != nil {
			return
		}
		fmt.Printf("after parse meta\n")
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

func (u *UploadsService) UpdateSwingUpload(_ context.Context, data *t.UpdateSwingUpload) (upload *t.SwingUpload, err error) {
	return u.Store.UpdateSwingUpload(data)
}

func (u *UploadsService) parseMetaFile(txtURL string) (meta *t.SwingUploadMeta, err error) {
	// download file
	res, err := http.Get(txtURL)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	fmt.Printf("after body close\n")

	stringContent := string(data)
	lineEnding := "\n"
	if windows := strings.Index(stringContent, "\r\n"); windows > -1 {
		lineEnding = "\r\n"
	}

	// parse file
	meta = &t.SwingUploadMeta{}
	for _, line := range strings.Split(stringContent, lineEnding) {
		if attr := strings.Split(line, "="); len(attr) > 1 {
			if attr[0] == "timestamp" {
				if meta.TimestampSeconds, err = strconv.Atoi(attr[1]); err != nil {
					return
				}
			} else if attr[0] == "frames" {
				if meta.Frames, err = strconv.Atoi(attr[1]); err != nil {
					return
				}
			} else if attr[0] == "swing" {
				if meta.Swing, err = strconv.Atoi(attr[1]); err != nil {
					return
				}
			} else if attr[0] == "clip" {
				if meta.Clip, err = strconv.Atoi(attr[1]); err != nil {
					return
				}
			}
		}
	}
	fmt.Printf("after parse file\n")
	return meta, err
}
