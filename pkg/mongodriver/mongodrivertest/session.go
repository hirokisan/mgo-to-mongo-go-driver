package mongodrivertest

import (
	"context"

	"github.com/hirokisan/mgo-to-mongo-go-driver/pkg/mongodriver"
)

const (
	testdbname = "test"
	url        = "mongodb://localhost:27017"
)

func NewTestDatabase(ctx context.Context) *mongodriver.Database {
	return mongodriver.NewDatabase(ctx, url, testdbname)
}
