package store

import (
	"context.Context"
	"time"

	uuid "github.com/satori/go.uuid"

	t "github.com/tennis-community-api-service/albums/types"
	m "github.com/tennis-community-api-service/pkg/mongo"
)

func (s *Store) CreateSwingUpload(_ context.Context, data *t.SwingUpload) (upload *t.SwingUpload, err error) {
	sess, c := s.C(ColSwingUploads)
	defer sess.Close()

	if data.ID == "" {
		data.ID = uuid.NewV4().String()
	}
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now

	upload = &t.SwingUpload{}
	err = m.Upsert(c, upload, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *Store) CreateUploadClipVideos(_ context.Context, uploadID string, clips []*t.UploadClipVideo) (upload *t.SwingUpload, err error) {
	sess, c := s.C(ColSwingUploads)
	defer sess.Close()

	upload = &t.SwingUpload{}
	err = m.Update(c, upload, m.M{"_id": uploadID}, m.M{
		"$set": m.M{
			"updAt":    time.Now(),
			"clipVids": clips,
		},
	})
	return
}

func (s *Store) CreateUploadSwingVideos(_ context.Context, uploadID string, swings []*t.UploadSwingVideo) (upload *t.SwingUpload, err error) {
	sess, c := s.C(ColSwingUploads)
	defer sess.Close()

	upload = &t.SwingUpload{}
	err = m.Update(c, upload, m.M{"_id": uploadID}, m.M{
		"$set":  m.M{"updAt": time.Now()},
		"$push": m.M{"swingVids": m.M{"$each": swings}},
	})
	return
}

// func (s *Store) UploadSwingVideoTranscoded(_ context.Context, uploadID, tranUrl string, swingVideoID int) (upload *t.SwingUpload, err error) {
// 	sess, c := s.C(ColSwingUploads)
// 	defer sess.Close()

// 	upload = &t.SwingUpload{}
// 	err = m.Update(c, upload, m.M{"_id": uploadID, "swingVids.id": swingVideoID}, m.M{
// 		"$set": m.M{
// 			"updAt":               time.Now(),
// 			"swingVids.$.tranUrl": tranUrl,
// 		},
// 	})
// 	return
// }
