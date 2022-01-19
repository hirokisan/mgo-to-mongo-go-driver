package bson_test

import (
	"context"
	"testing"

	"github.com/hirokisan/mgo-to-mongo-go-driver/bson"
	"github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mgo/mgotest"
	"github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mongodriver/mongodrivertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCollectionName = "test"

type target struct {
	ID         bson.ObjectID  `bson:"_id"`
	PID        *bson.ObjectID `bson:"pid"`
	NullPID    *bson.ObjectID `bson:"nullpid"`
	OmitPID    *bson.ObjectID `bson:"omitpid,omitempty"`
	InsertFrom string         `bson:"insertFrom"`
}

func TestOne(t *testing.T) {
	ctx := context.Background()

	db := mongodrivertest.NewTestDatabase(ctx)
	mdCollection := db.Collection(testCollectionName)

	session := mgotest.NewTestSession()
	defer session.Close()
	mgoCollection := session.Collection(testCollectionName)

	t.Run("insert from mgo", func(t *testing.T) {
		primitiveID := bson.NewObjectID()

		m := target{
			ID:         primitiveID,
			PID:        &primitiveID,
			InsertFrom: "mgo",
		}
		assert.NoError(t, mgoCollection.Insert(m))

		t.Run("find from mgo", func(t *testing.T) {
			var got target
			require.NoError(t, mgoCollection.FindId(m.ID).One(&got))
			assert.Equal(t, primitiveID, got.ID)
			assert.Equal(t, primitiveID, *got.PID)
			assert.Nil(t, got.NullPID)
			assert.Nil(t, got.OmitPID)
		})
		t.Run("find from mongodriver", func(t *testing.T) {
			var got target
			require.NoError(t, mdCollection.FindOne(ctx, bson.M{"_id": primitiveID}).Decode(&got))
			assert.Equal(t, primitiveID, got.ID)
			assert.Equal(t, primitiveID, got.ID)
			assert.Equal(t, primitiveID, *got.PID)
			assert.Nil(t, got.NullPID)
			assert.Nil(t, got.OmitPID)
		})
	})
	t.Run("insert from mongodriver", func(t *testing.T) {
		primitiveID := bson.NewObjectID()
		m := target{
			ID:         primitiveID,
			PID:        &primitiveID,
			InsertFrom: "mongo-go-driver",
		}
		_, err := mdCollection.InsertOne(ctx, m)
		require.NoError(t, err)

		t.Run("find from mgo", func(t *testing.T) {
			var got target
			require.NoError(t, mgoCollection.FindId(primitiveID).One(&got))
			assert.Equal(t, primitiveID, got.ID)
			assert.Equal(t, primitiveID, *got.PID)
			assert.Nil(t, got.NullPID)
			assert.Nil(t, got.OmitPID)
		})
		t.Run("find from mongodriver", func(t *testing.T) {
			var got target
			assert.NoError(t, mdCollection.FindOne(ctx, bson.M{"_id": primitiveID}).Decode(&got))
			assert.Equal(t, primitiveID, got.ID)
			assert.Equal(t, primitiveID, *got.PID)
			assert.Nil(t, got.NullPID)
			assert.Nil(t, got.OmitPID)
		})
	})
}
