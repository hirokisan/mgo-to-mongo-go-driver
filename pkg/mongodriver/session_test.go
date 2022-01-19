package mongodriver_test

import (
	"context"
	"testing"

	"github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mongodriver/mongodrivertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Target struct {
	ID int `bson:"id"`
}

func TestInsert(t *testing.T) {
	ctx := context.Background()
	db := mongodrivertest.NewTestDatabase(ctx)

	col := db.Collection("target")
	target := Target{ID: 1}
	_, err := col.InsertOne(ctx, target)
	require.NoError(t, err)

	var got Target
	require.NoError(t, col.FindOne(ctx, primitive.M{"id": target.ID}).Decode(&got))

	assert.Equal(t, got.ID, target.ID)
}
