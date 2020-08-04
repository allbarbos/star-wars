package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type database struct {
	DB *mongo.Database
}

func New() Database {
	name := os.Getenv("DB_NAME")
	uri := os.Getenv("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, options)

	if err != nil {
		log.Fatal(err)
	}

	return &database{
		DB: client.Database(name),
	}
}

func (d database) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return d.DB.Collection(name, opts...)
}
