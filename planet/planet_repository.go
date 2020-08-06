package planet

import (
	"context"
	"log"
	"os"
	"star-wars/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const planetsCollection = "planets"

// Repository contract
type Repository interface {
	FindByName(name string) (entity.Planet, error)
	Save(planet entity.Planet) error
}

type repo struct {
	DB *mongo.Database
}

// NewRepository returns a planet repository instance
func NewRepository() Repository {
	name := os.Getenv("DB_NAME")
	uri := os.Getenv("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	options := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, options)

	if err != nil {
		log.Fatal(err)
	}

	return &repo{
		DB: client.Database(name),
	}
}

func (r repo) FindByName(name string) (entity.Planet, error) {
	filter := bson.M{"name": name}
	planet := entity.Planet{}

	err := r.DB.Collection(planetsCollection).FindOne(context.Background(), filter).Decode(&planet)

	if err != nil {
		log.Print(err)
		return planet, err
	}

	return planet, nil
}

func (r repo) Save(planet entity.Planet) error {
	result, err := r.DB.Collection(planetsCollection).InsertOne(context.Background(), planet)

	if err != nil {
		log.Print(err)
		return err
	}

	log.Print(result)
	return nil
}
