package store

import (
	"github.com/globalsign/mgo"
)

type UsersStore struct {
	m      *mgo.Session
	dbname string
}

func NewStore(m *mgo.Session, dbname string) *UsersStore {
	return &UsersStore{
		m:      m,
		dbname: dbname,
	}
}

func (s *UsersStore) DB() (*mgo.Session, *mgo.Database) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname)
}

func (s *UsersStore) C(colName string) (*mgo.Session, *mgo.Collection) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname).C(colName)
}
