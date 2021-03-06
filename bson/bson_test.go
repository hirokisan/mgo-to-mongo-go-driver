package bson_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/hirokisan/mgo-to-mongo-go-driver/bson"
	"github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mgo/mgotest"
	"github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mongodriver/mongodrivertest"
	mgobson "github.com/hirokisan/mgo/bson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	driverbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const testCollectionName = "test"

type target struct {
	ID         bson.ObjectID  `bson:"_id"`
	PID        *bson.ObjectID `bson:"pid"`
	NullPID    *bson.ObjectID `bson:"nullpid"`
	OmitPID    *bson.ObjectID `bson:"omitpid,omitempty"`
	InsertFrom string         `bson:"insertFrom"`
}

func TestObjectID(t *testing.T) {
	ctx := context.Background()

	db := mongodrivertest.NewTestDatabase(ctx)
	mdCollection := db.Collection(testCollectionName)

	session := mgotest.NewTestSession()
	defer session.Close()
	mgoCollection := session.Collection(testCollectionName)

	t.Run("insert from mgo", func(t *testing.T) {
		objectID := bson.NewObjectID()

		m := target{
			ID:         objectID,
			PID:        &objectID,
			InsertFrom: "mgo",
		}
		assert.NoError(t, mgoCollection.Insert(m))

		t.Run("find from mgo", func(t *testing.T) {
			var got target
			require.NoError(t, mgoCollection.FindId(m.ID).One(&got))
			assert.Equal(t, objectID, got.ID)
			assert.Equal(t, objectID, *got.PID)
			assert.Nil(t, got.NullPID)
			assert.Nil(t, got.OmitPID)
		})
		t.Run("find from mongodriver", func(t *testing.T) {
			var got target
			require.NoError(t, mdCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&got))
			assert.Equal(t, objectID, got.ID)
			assert.Equal(t, objectID, *got.PID)
			assert.Nil(t, got.NullPID)
			assert.Nil(t, got.OmitPID)
		})
	})
	t.Run("insert from mongodriver", func(t *testing.T) {
		objectID := bson.NewObjectID()
		m := target{
			ID:         objectID,
			PID:        &objectID,
			InsertFrom: "mongo-go-driver",
		}
		_, err := mdCollection.InsertOne(ctx, m)
		require.NoError(t, err)

		t.Run("find from mgo", func(t *testing.T) {
			var got target
			require.NoError(t, mgoCollection.FindId(objectID).One(&got))
			assert.Equal(t, objectID, got.ID)
			assert.Equal(t, objectID, *got.PID)
			assert.Nil(t, got.NullPID)
			assert.Nil(t, got.OmitPID)
		})
		t.Run("find from mongodriver", func(t *testing.T) {
			var got target
			assert.NoError(t, mdCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&got))
			assert.Equal(t, objectID, got.ID)
			assert.Equal(t, objectID, *got.PID)
			assert.Nil(t, got.NullPID)
			assert.Nil(t, got.OmitPID)
		})
	})
}

func TestM(t *testing.T) {
	ctx := context.Background()

	db := mongodrivertest.NewTestDatabase(ctx)
	mdCollection := db.Collection(testCollectionName)

	session := mgotest.NewTestSession()
	defer session.Close()
	mgoCollection := session.Collection(testCollectionName)

	t.Run("insert from mgo", func(t *testing.T) {
		objectID := bson.NewObjectID()

		m := bson.M{
			"_id":        objectID,
			"pid":        &objectID,
			"time":       time.Now(),
			"insertFrom": "mgo",
		}
		assert.NoError(t, mgoCollection.Insert(m))

		t.Run("find from mgo", func(t *testing.T) {
			t.Run("FindId", func(t *testing.T) {
				var got target
				require.NoError(t, mgoCollection.FindId(objectID).One(&got))
				assert.Equal(t, objectID, got.ID)
				assert.Equal(t, objectID, *got.PID)
			})
			t.Run("Find", func(t *testing.T) {
				var got target
				require.NoError(t, mgoCollection.Find(bson.M{"_id": objectID}).One(&got))
				assert.Equal(t, objectID, got.ID)
				assert.Equal(t, objectID, *got.PID)
			})
		})
		t.Run("find from mongodriver", func(t *testing.T) {
			var got target
			require.NoError(t, mdCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&got))
			assert.Equal(t, objectID, got.ID)
			assert.Equal(t, objectID, *got.PID)
		})
	})
	t.Run("insert from mongodriver", func(t *testing.T) {
		objectID := bson.NewObjectID()
		m := bson.M{
			"_id":        objectID,
			"pid":        &objectID,
			"time":       time.Now(),
			"insertFrom": "mongo-go-driver",
		}
		_, err := mdCollection.InsertOne(ctx, m)
		require.NoError(t, err)

		t.Run("find from mgo", func(t *testing.T) {
			t.Run("FindId", func(t *testing.T) {
				var got target
				require.NoError(t, mgoCollection.FindId(objectID).One(&got))
				assert.Equal(t, objectID, got.ID)
				assert.Equal(t, objectID, *got.PID)
			})
			t.Run("Find", func(t *testing.T) {
				var got target
				require.NoError(t, mgoCollection.Find(bson.M{"_id": objectID}).One(&got))
				assert.Equal(t, objectID, got.ID)
				assert.Equal(t, objectID, *got.PID)
			})
		})
		t.Run("find from mongodriver", func(t *testing.T) {
			var got target
			assert.NoError(t, mdCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&got))
			assert.Equal(t, objectID, got.ID)
			assert.Equal(t, objectID, *got.PID)
		})
	})
}

func TestObjectIDHex(t *testing.T) {
	idString := "61e77d72f768730001eacbe2"

	want := bson.ObjectIDHex(idString)

	t.Run("compare with mgobson, value is the same", func(t *testing.T) {
		got := mgobson.ObjectIdHex(idString)
		require.Equal(t, reflect.TypeOf(bson.ObjectID("")), reflect.TypeOf(want))
		require.Equal(t, reflect.TypeOf(mgobson.ObjectId("")), reflect.TypeOf(got))
		assert.Equal(t, want.Hex(), got.Hex())
	})
}

func TestNewObjectID(t *testing.T) {
	ctx := context.Background()

	db := mongodrivertest.NewTestDatabase(ctx)
	mdCollection := db.Collection(testCollectionName)

	session := mgotest.NewTestSession()
	defer session.Close()
	mgoCollection := session.Collection(testCollectionName)

	type target struct {
		ID bson.ObjectID `bson:"_id"`
	}
	tgt := target{
		ID: bson.NewObjectID(),
	}
	_, err := mdCollection.InsertOne(ctx, tgt)
	require.NoError(t, err)
	t.Run("find from mgo, check it as ObjectID", func(t *testing.T) {
		var want struct {
			ID mgobson.ObjectId `bson:"_id"`
		}
		require.NoError(t, mgoCollection.FindId(tgt.ID).One(&want))
		assert.Equal(t, tgt.ID.Hex(), want.ID.Hex())
		assert.True(t, mgobson.IsObjectIdHex(want.ID.Hex()))
	})
	t.Run("find from mongo-go-driver, check it as ObjectID", func(t *testing.T) {
		var want struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		require.NoError(t, mdCollection.FindOne(ctx, bson.M{"_id": tgt.ID}).Decode(&want))
		assert.Equal(t, tgt.ID.Hex(), want.ID.Hex())
		assert.True(t, primitive.IsValidObjectID(want.ID.Hex()))
	})
}

func TestIsObjectIDHex(t *testing.T) {
	ctx := context.Background()

	db := mongodrivertest.NewTestDatabase(ctx)
	mdCollection := db.Collection(testCollectionName)

	type target struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	tgt := target{
		ID: primitive.NewObjectID(),
	}
	_, err := mdCollection.InsertOne(ctx, tgt)
	require.NoError(t, err)
	{
		var want struct {
			ID bson.ObjectID `bson:"_id"`
		}
		require.NoError(t, mdCollection.FindOne(ctx, bson.M{"_id": bson.ObjectIDHex(tgt.ID.Hex())}).Decode(&want))
		assert.True(t, bson.IsObjectIDHex(want.ID.Hex()))
	}
}

func TestUnmarshal(t *testing.T) {
	t.Run("from []byte to struct", func(t *testing.T) {
		objectID := bson.NewObjectID()
		tgt := target{
			ID:  objectID,
			PID: &objectID,
		}
		buf, err := bson.Marshal(tgt)
		require.NoError(t, err)
		var want target
		require.NoError(t, bson.Unmarshal(buf, &want))
		require.Equal(t, tgt.ID, want.ID)
		t.Run("compare with mgobson", func(t *testing.T) {
			buf, err := mgobson.Marshal(tgt)
			require.NoError(t, err)
			var got target
			require.NoError(t, mgobson.Unmarshal(buf, &got))
			assert.Equal(t, want.ID, got.ID)
			assert.Equal(t, want.PID, got.PID)
		})
		t.Run("compare with primitive", func(t *testing.T) {
			buf, err := driverbson.Marshal(tgt)
			require.NoError(t, err)
			var got target
			require.NoError(t, driverbson.Unmarshal(buf, &got))
			assert.Equal(t, want.ID, got.ID)
			assert.Equal(t, want.PID, got.PID)
		})
	})
	t.Run("from []byte to bson.M", func(t *testing.T) {
		objectID := bson.NewObjectID()
		tgt := target{
			ID:  objectID,
			PID: &objectID,
		}
		buf, err := bson.Marshal(tgt)
		require.NoError(t, err)
		var want bson.M
		require.NoError(t, bson.Unmarshal(buf, &want))
		require.Equal(t, tgt.ID, want["_id"])
		require.Equal(t, *tgt.PID, want["pid"])
		t.Run("compare with mgobson, different types but matching values", func(t *testing.T) {
			buf, err := mgobson.Marshal(tgt)
			require.NoError(t, err)
			var got bson.M
			require.NoError(t, mgobson.Unmarshal(buf, &got))
			assert.NotEqual(t, want["_id"], got["_id"])
			assert.Equal(t, want["_id"].(bson.ObjectID).Hex(), got["_id"].(mgobson.ObjectId).Hex())
			assert.NotEqual(t, want["pid"], got["pid"])
			assert.Equal(t, want["pid"].(bson.ObjectID).Hex(), got["pid"].(mgobson.ObjectId).Hex())
		})
		t.Run("compare with primitive, different types but matching values", func(t *testing.T) {
			buf, err := driverbson.Marshal(tgt)
			require.NoError(t, err)
			var got bson.M
			require.NoError(t, driverbson.Unmarshal(buf, &got))
			assert.NotEqual(t, want["_id"], got["_id"])
			assert.Equal(t, want["_id"].(bson.ObjectID).Hex(), got["_id"].(primitive.ObjectID).Hex())
			assert.NotEqual(t, want["pid"], got["pid"])
			assert.Equal(t, want["pid"].(bson.ObjectID).Hex(), got["pid"].(primitive.ObjectID).Hex())
		})
	})
}

func TestMarshal(t *testing.T) {
	t.Run("from struct to []byte", func(t *testing.T) {
		objectID := bson.NewObjectID()
		tgt := target{
			ID:  objectID,
			PID: &objectID,
		}
		want, err := bson.Marshal(tgt)
		require.NoError(t, err)
		t.Run("compare with mgobson", func(t *testing.T) {
			got, err := mgobson.Marshal(tgt)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
		t.Run("compare with primitive", func(t *testing.T) {
			got, err := driverbson.Marshal(tgt)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
	})
	t.Run("from bson.M to []byte", func(t *testing.T) {
		objectID := bson.NewObjectID()
		// NOTE: bson.M do not care about the order, so set only one
		m := bson.M{
			"_id": objectID,
		}
		want, err := bson.Marshal(m)
		require.NoError(t, err)
		t.Run("compare with mgobson", func(t *testing.T) {
			got, err := mgobson.Marshal(m)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
		t.Run("compare with primitive", func(t *testing.T) {
			got, err := driverbson.Marshal(m)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
	})
}
