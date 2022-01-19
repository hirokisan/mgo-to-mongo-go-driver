package mongodriver

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	*mongo.Client
	db string
}

func NewDatabase(ctx context.Context, url, db string) *Database {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	return &Database{
		Client: client,
		db:     db,
	}
}

func (db *Database) Collection(name string) *mongo.Collection {
	return db.Database(db.db).Collection(name)
}
