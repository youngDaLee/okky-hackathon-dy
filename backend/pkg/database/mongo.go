package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	return mongo.Connect(options.Client().ApplyURI(uri))
}

func GetCollection(client *mongo.Client, db, collection string) *mongo.Collection {
	return client.Database(db).Collection(collection)
}
