package mongodriver_test

import (
	"context"
	"testing"

	"github.com/hirokisan/mgo-to-mongo-go-driver/model"
	"github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mongodriver/mongodrivertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsert(t *testing.T) {
	ctx := context.Background()
	db := mongodrivertest.NewTestDatabase(ctx)

	col := db.Collection("cup")
	cup := model.Cup{ID: 1}
	_, err := col.InsertOne(ctx, cup)
	require.NoError(t, err)

	var got model.Cup
	require.NoError(t, col.FindOne(ctx, primitive.M{"id": cup.ID}).Decode(&got))

	assert.Equal(t, got.ID, cup.ID)
}
