package store

import (
	"github.com/globalsign/mgo"
)

type ModerationStore struct {
	m      *mgo.Session
	dbname string
}

func NewStore(m *mgo.Session, dbname string) *ModerationStore {
	return &ModerationStore{
		m:      m,
		dbname: dbname,
	}
}

func (s *ModerationStore) DB() (*mgo.Session, *mgo.Database) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname)
}

func (s *ModerationStore) C(colName string) (*mgo.Session, *mgo.Collection) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname).C(colName)
}
