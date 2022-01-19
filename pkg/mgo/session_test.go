package mgo_test

import (
	"testing"

	"github.com/hirokisan/mgo/bson"

	"github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mgo/mgotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Target struct {
	ID int `bson:"id"`
}

func TestInsert(t *testing.T) {
	session := mgotest.NewTestSession()
	defer session.Close()

	col := session.Collection("target")
	target := Target{ID: 1}
	require.NoError(t, col.Insert(target))

	var got Target
	require.NoError(t, col.Find(bson.M{"id": target.ID}).One(&got))

	assert.Equal(t, got.ID, target.ID)
}
