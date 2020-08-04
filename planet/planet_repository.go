package planet

import (
	"context"
	"log"
	"star-wars/db"
	"star-wars/entity"

	"go.mongodb.org/mongo-driver/bson"
)

const planetsCollection = "planets"

// Repository contract
type Repository interface {
	FindByName(name string) (entity.Planet, error)
	Save(planet entity.Planet) error
}

type repo struct {
	DB db.Database
}

// NewRepository returns a planet repository instance
func NewRepository(db db.Database) Repository {
	return &repo{
		DB: db,
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
