package store

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"

	m "github.com/tennis-community-api-service/pkg/mongo"
	t "github.com/tennis-community-api-service/uploads/types"
)

func (s *UploadsStore) GetRecentSwingUploads(userID string) (uploads []*t.SwingUpload, err error) {
	sess, c := s.C(ColSwingUploads)
	defer sess.Close()

	uploads = []*t.SwingUpload{}
	err = m.Aggregate(c, &uploads, []m.M{
		{"$match": m.M{"usrId": userID}},
		{"$sort": m.M{"crAt": -1}},
		{"$limit": 5},
	})
	return
}

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

func (s *UploadsStore) UpdateSwingUpload(data *t.UpdateSwingUpload) (upload *t.SwingUpload, err error) {
	sess, c := s.C(ColSwingUploads)
	defer sess.Close()

	if data.UploadKey == "" || data.UserID == "" {
		return nil, errors.New("Update swing upload missing upload key or user id")
	}

	upload = &t.SwingUpload{}
	err = m.Update(c, upload, m.M{"upKey": data.UploadKey, "usrId": data.UserID}, m.M{"$set": data})
	return
}

func (s *UploadsStore) CreateUploadSwingVideos(userID, uploadKey string, clipNum int, swings []*t.UploadSwingVideo) (upload *t.SwingUpload, err error) {
	sess, c := s.C(ColSwingUploads)
	defer sess.Close()

	upload = &t.SwingUpload{}
	err = m.Update(c, upload, m.M{"usrId": userID, "upKey": uploadKey}, m.M{
		"$set": m.M{"updAt": time.Now()},
		"$push": m.M{
			"swingVids": m.M{"$each": swings},
			"procClips": clipNum,
		},
	})
	return
}
