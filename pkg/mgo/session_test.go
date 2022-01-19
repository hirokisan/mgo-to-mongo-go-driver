package mgo_test

import (
	"testing"

	"github.com/hirokisan/mgo/bson"

	"github.com/hirokisan/mgo-to-mongo-go-driver/model"
	"github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mgo/mgotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	session := mgotest.NewTestSession()
	defer session.Close()

	col := session.Collection("cup")
	cup := model.Cup{ID: 1}
	require.NoError(t, col.Insert(cup))

	var got model.Cup
	require.NoError(t, col.Find(bson.M{"id": cup.ID}).One(&got))

	assert.Equal(t, got.ID, cup.ID)
}
