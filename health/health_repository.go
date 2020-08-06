package health

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Repository contract
type Repository interface {
	Ping() string
}

type repo struct {}

// NewRepository health check repository instance
func NewRepository() Repository {
	return &repo{}
}

// Ping check conection
func (r repo) Ping() string {
	uri := os.Getenv("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	options := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, options)
	defer client.Disconnect(ctx)

	if err != nil {
		log.Print(err)
		return "error"
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Print(err)
		return "error"
	}

	return "ok"
}
