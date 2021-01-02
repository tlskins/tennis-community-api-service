package store

import (
	uuid "github.com/satori/go.uuid"

	t "github.com/tennis-community-api-service/albums/types"
	m "github.com/tennis-community-api-service/pkg/mongo"
)

func (s *AlbumsStore) CreateSwing(data *t.SwingVideo) (swing *t.SwingVideo, err error) {
	sess, c := s.C(ColSwings)
	defer sess.Close()

	if data.ID == "" {
		data.ID = uuid.NewV4().String()
	}
	swing = &t.SwingVideo{}
	err = m.Upsert(c, swing, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *AlbumsStore) UpdateSwing(data *t.UpdateSwingVideo) (swing *t.SwingVideo, err error) {
	sess, c := s.C(ColSwings)
	defer sess.Close()

	swing = &t.SwingVideo{}
	err = m.Update(c, swing, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}
