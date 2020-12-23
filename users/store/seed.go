package store

import (
	"github.com/globalsign/mgo"
)

func (s *UsersStore) EnsureIndexes() error {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	if err := c.EnsureIndex(mgo.Index{
		Key:    []string{"em"},
		Unique: true,
	}); err != nil {
		return err
	}
	return nil
}
