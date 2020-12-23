package store

import (
	"github.com/globalsign/mgo"
)

type Store struct {
	m      *mgo.Session
	dbname string
}

func NewStore(m *mgo.Session, dbname string) *Store {
	return &Store{
		m:      m,
		dbname: dbname,
	}
}

func (s *Store) DB() (*mgo.Session, *mgo.Database) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname)
}

func (s *Store) C(colName string) (*mgo.Session, *mgo.Collection) {
	sess := s.m.Copy()
	return sess, sess.DB(s.dbname).C(colName)
}
