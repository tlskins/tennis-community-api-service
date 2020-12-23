package store

import (
	"github.com/globalsign/mgo"
)

type UploadsStore struct {
	m      *mgo.Session
	dbname string
}

func NewStore(m *mgo.Session, dbname string) *UploadsStore {
	return &UploadsStore{
		m:      m,
		dbname: dbname,
	}
}

func (s *UploadsStore) DB() (*mgo.Session, *mgo.Database) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname)
}

func (s *UploadsStore) C(colName string) (*mgo.Session, *mgo.Collection) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname).C(colName)
}
