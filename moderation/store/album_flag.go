package store

import (
	"time"

	uuid "github.com/satori/go.uuid"

	t "github.com/tennis-community-api-service/moderation/types"
	m "github.com/tennis-community-api-service/pkg/mongo"
)

func (s *ModerationStore) CreateAlbumFlag(data *t.AlbumFlag) (flag *t.AlbumFlag, err error) {
	sess, c := s.C(ColFlaggedAlbums)
	defer sess.Close()

	if data.ID == "" {
		data.ID = uuid.NewV4().String()
	}
	flag = &t.AlbumFlag{}
	err = m.Upsert(c, flag, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *ModerationStore) UpdateAlbumFlag(data *t.UpdateAlbumFlag) (flag *t.AlbumFlag, err error) {
	sess, c := s.C(ColFlaggedAlbums)
	defer sess.Close()

	flag = &t.AlbumFlag{}
	err = m.Update(c, flag, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *ModerationStore) RecentFlaggedAlbums(start, end time.Time, resolved *bool, limit, offset int) (flags []*t.AlbumFlag, err error) {
	sess, c := s.C(ColFlaggedAlbums)
	defer sess.Close()

	flags = []*t.AlbumFlag{}
	query := m.M{"crAt": m.M{"$gte": start, "$lt": end}}

	if resolved != nil {
		query["res"] = *resolved
	}

	if limit > 0 {
		err = c.Find(query).Skip(offset).Limit(limit).All(&flags)
	} else {
		err = m.Find(c, &flags, query, nil)
	}
	return
}
