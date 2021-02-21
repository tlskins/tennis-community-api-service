package store

import (
	"github.com/globalsign/mgo"
)

func (s *AlbumsStore) EnsureIndexes() error {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	if err := c.EnsureIndex(mgo.Index{
		Key: []string{
			"userId",
			"public",
			"home",
			"frndView",
			"frndIds",
		},
	}); err != nil {
		return err
	}
	return nil
}
