package planet

import (
	"context"
	"log"
	"os"
	"star-wars/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Repository contract
type Repository interface {
	FindAll(ctx context.Context, limit int64, skip int64) (*[]entity.Planet, error)
	FindByName(ctx context.Context, name string) (*entity.Planet, error)
	FindByID(ctx context.Context, id string) (*entity.Planet, error)
	Save(ctx context.Context, planet *entity.Planet) error
	Delete(ctx context.Context, id string) error
	Ping(ctx context.Context) string
}

type repo struct {}

// NewRepository planet
func NewRepository() Repository {
	return &repo{}
}

func cnx(ctx context.Context) (*mongo.Collection, error){
	c, _ := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_HOST")))
	err := c.Connect(ctx)
	coll := c.Database(os.Getenv("DB_NAME")).Collection("planets")

	return coll, err
}

func (r repo) FindAll(ctx context.Context, limit int64, skip int64) (*[]entity.Planet, error) {
	coll, err := cnx(ctx)

	if err != nil {
		return nil, err
	}

	defer coll.Database().Client().Disconnect(ctx)

	opt := options.Find()
	opt.SetLimit(limit)
	opt.SetSkip(skip)

	cr, err := coll.Find(ctx, bson.D{}, opt)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	planets := &[]entity.Planet{}

	err = cr.All(ctx, planets)

	if err != nil {
		return nil, err
	}

	return planets, nil
}

func (r repo) FindByName(ctx context.Context, name string) (*entity.Planet, error) {
	coll, err := cnx(ctx)

	if err != nil {
		return nil, err
	}

	defer coll.Database().Client().Disconnect(ctx)

	var planet entity.Planet

	err = coll.FindOne(
		ctx,
		bson.M{"name": name},
	).Decode(&planet)

	if err != nil {
		return nil, err
	}

	return &planet, nil
}

func (r repo) FindByID(ctx context.Context, id string) (*entity.Planet, error) {
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	coll, err := cnx(ctx)

	if err != nil {
		return nil, err
	}

	defer coll.Database().Client().Disconnect(ctx)

	var planet entity.Planet

	err = coll.FindOne(
		ctx,
		bson.M{"_id": _id},
	).Decode(&planet)

	if err != nil {
		return nil, err
	}

	return &planet, nil
}

func (r repo) Save(ctx context.Context, planet *entity.Planet) error {
	coll, err := cnx(ctx)

	if err != nil {
		return err
	}

	defer coll.Database().Client().Disconnect(ctx)

	result, err := coll.InsertOne(ctx, &planet)

	if err != nil {
		return err
	}

	oid, _ := result.InsertedID.(primitive.ObjectID)
	planet.ID = oid.Hex()

	return nil
}

func (r repo) Delete(ctx context.Context, id string) error {
	coll, err := cnx(ctx)

	if err != nil {
		return err
	}

	defer coll.Database().Client().Disconnect(ctx)

	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = coll.DeleteOne(ctx, bson.M{"_id": _id})

	if err != nil {
		return err
	}

	return nil
}

func (r repo) Ping(ctx context.Context) string {
	coll, err := cnx(ctx)

	if err != nil {
		return "error"
	}

	defer coll.Database().Client().Disconnect(ctx)

	err = coll.Database().Client().Ping(ctx, readpref.Primary())
	if err != nil {
		return "error"
	}
	return "ok"
}
