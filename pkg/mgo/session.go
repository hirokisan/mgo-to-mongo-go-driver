package mgo

import (
	"gopkg.in/mgo.v2"
)

type Session struct {
	*mgo.Session
	db string
}

func NewSession(url, db string) *Session {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	return &Session{
		Session: session,
		db:      db,
	}
}

func (s *Session) Close() {
	s.Session.Close()
}

func (s *Session) Collection(name string) *mgo.Collection {
	return s.DB(s.db).C(name)
}
