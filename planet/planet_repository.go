package planet

import (
	"context"
	"log"
	"os"
	"star-wars/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Repository contract
type Repository interface {
	FindByID(id string) (entity.Planet, error)
	FindByName(name string) (entity.Planet, error)
	Save(planet entity.Planet) error
	Ping() string
}

type repo struct {
	Options *options.ClientOptions
}

// NewRepository planet
func NewRepository() Repository {
	uri := os.Getenv("DB_HOST")
	options := options.Client().ApplyURI(uri)
	return &repo{
		Options: options,
	}
}

func db(ctx context.Context, r repo) (*mongo.Client, *mongo.Collection, error){
	cli, err := mongo.Connect(ctx, r.Options)
	col := cli.Database(os.Getenv("DB_NAME")).Collection("planets")
	return cli, col, err
}

// Ping check connection
func (r repo) Ping() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cnx, _, err := db(ctx, r)

	err = cnx.Ping(ctx, readpref.Primary())
	defer cnx.Disconnect(ctx)

	if err != nil {
		log.Print(err)
		return "error"
	}

	return "ok"
}

func (r repo) FindByName(name string) (entity.Planet, error) {
	filter := bson.M{"name": name}
	planet := entity.Planet{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cnx, db, err := db(ctx, r)

	if err != nil {
		log.Print(err)
		return planet, err
	}

	err = db.FindOne(ctx, filter).Decode(&planet)
	defer cnx.Disconnect(ctx)

	if err != nil {
		log.Print(err)
		return planet, err
	}

	return planet, nil
}

func (r repo) FindByID(id string) (entity.Planet, error) {
	planet := entity.Planet{}
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Print(err)
		return planet, err
	}

	filter := bson.M{"_id": _id}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cnx, db, err := db(ctx, r)

	if err != nil {
		log.Print(err)
		return planet, err
	}

	err = db.FindOne(ctx, filter).Decode(&planet)
	defer cnx.Disconnect(ctx)

	if err != nil {
		log.Print(err)
		return planet, err
	}

	return planet, nil
}


func (r repo) Save(planet entity.Planet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cnx, db, err := db(ctx, r)

	if err != nil {
		log.Print(err)
		return err
	}

	_, err = db.InsertOne(ctx, planet)
	defer cnx.Disconnect(ctx)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
