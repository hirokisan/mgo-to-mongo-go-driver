package mgotest

import "github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mgo"

const (
	testdbname = "test"
	url        = "mongodb://localhost:27017"
)

func NewTestSession() *mgo.Session {
	return mgo.NewSession(url, testdbname)
}
