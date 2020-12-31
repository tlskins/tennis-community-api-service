package store

import (
	"github.com/globalsign/mgo"
)

type AlbumsStore struct {
	m      *mgo.Session
	dbname string
}

func NewStore(m *mgo.Session, dbname string) *AlbumsStore {
	return &AlbumsStore{
		m:      m,
		dbname: dbname,
	}
}

func (s *AlbumsStore) DB() (*mgo.Session, *mgo.Database) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname)
}

func (s *AlbumsStore) C(colName string) (*mgo.Session, *mgo.Collection) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname).C(colName)
}
