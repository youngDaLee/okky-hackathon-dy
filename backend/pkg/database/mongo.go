package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	return nil, nil
}

func GetCollection(client *mongo.Client, db, collection string) *mongo.Collection {
	return nil
}
