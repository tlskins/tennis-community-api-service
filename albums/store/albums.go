package store

import (
	"time"

	uuid "github.com/satori/go.uuid"

	t "github.com/tennis-community-api-service/albums/types"
	m "github.com/tennis-community-api-service/pkg/mongo"
)

func (s *AlbumsStore) CreateAlbum(data *t.Album) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	if data.ID == "" {
		data.ID = uuid.NewV4().String()
	}
	album = &t.Album{}
	err = m.Upsert(c, album, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *AlbumsStore) UpdateAlbum(data *t.UpdateAlbum) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	album = &t.Album{}
	err = m.Update(c, album, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *AlbumsStore) AddVideosToAlbum(userId, uploadKey string, swings []*t.SwingVideo) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	album = &t.Album{}
	err = m.Update(c, album, m.M{"userId": userId, "upKey": uploadKey}, m.M{
		"$set":  m.M{"updAt": time.Now()},
		"$push": m.M{"swingVids": m.M{"$each": swings}},
	})
	return
}

// func (s *AlbumsStore) UpdateAlbum(data *t.UpdateSwingUpload) (upload *t.SwingUpload, err error) {
// 	sess, c := s.C(ColSwingUploads)
// 	defer sess.Close()

// 	upload = &t.SwingUpload{}
// 	err = m.Update(c, upload, m.M{"upKey": data.UploadKey}, m.M{"$set": data})
// 	return
// }
