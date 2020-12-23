package store

import (
	"time"

	uuid "github.com/satori/go.uuid"

	m "github.com/tennis-community-api-service/pkg/mongo"
	t "github.com/tennis-community-api-service/uploads/types"
)

func (s *UploadsStore) CreateSwingUpload(data *t.SwingUpload) (upload *t.SwingUpload, err error) {
	sess, c := s.C(ColSwingUploads)
	defer sess.Close()

	if data.ID == "" {
		data.ID = uuid.NewV4().String()
	}
	upload = &t.SwingUpload{}
	err = m.Upsert(c, upload, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *UploadsStore) CreateUploadClipVideos(uploadID string, clips []*t.UploadClipVideo) (upload *t.SwingUpload, err error) {
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

func (s *UploadsStore) CreateUploadSwingVideos(uploadID string, swings []*t.UploadSwingVideo) (upload *t.SwingUpload, err error) {
	sess, c := s.C(ColSwingUploads)
	defer sess.Close()

	upload = &t.SwingUpload{}
	err = m.Update(c, upload, m.M{"_id": uploadID}, m.M{
		"$set":  m.M{"updAt": time.Now()},
		"$push": m.M{"swingVids": m.M{"$each": swings}},
	})
	return
}

// func (s *UploadsStore) UploadSwingVideoTranscoded(uploadID, tranUrl string, swingVideoID int) (upload *t.SwingUpload, err error) {
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
