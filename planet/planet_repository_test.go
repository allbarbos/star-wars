package planet

import (
	"context"
	"errors"
	"reflect"
	"star-wars/entity"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func monkeyMongoDBAll(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
	var cr *mongo.Cursor
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(cr), "All",
		func(cr *mongo.Cursor, ctx context.Context, results interface{}) error {
			guard.Unpatch()
			defer guard.Restore()

			if err {
				return errors.New("cursor all error")
			}

			return nil
		})
	return guard
}

func monkeyMongoDBFind(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
	var coll *mongo.Collection
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(coll), "Find",
		func(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
			guard.Unpatch()
			defer guard.Restore()
			if err {
				return nil, errors.New("find error")
			}
			return &mongo.Cursor{}, nil
		})
	return guard
}

// func monkeyMongoDBConnect(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
// 	var obj *mongo.Client
// 	guard = monkey.PatchInstanceMethod(reflect.TypeOf(obj), "Connect",
// 		func(obj *mongo.Client, ctx context.Context) error {
// 			guard.Unpatch()
// 			defer guard.Restore()
// 			if err {
// 				return errors.New("connect error")
// 			}
// 			return nil
// 		})
// 	return guard
// }

func monkeyMongoDBFindOne(guard *monkey.PatchGuard) *monkey.PatchGuard {
	var coll *mongo.Collection
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(coll), "FindOne",
		func(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
			guard.Unpatch()
			defer guard.Restore()
			return &mongo.SingleResult{}
		})
	return guard
}

func monkeyMongoDBDecode(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
	var obj *mongo.SingleResult
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(obj), "Decode",
		func(obj *mongo.SingleResult, v interface{}) error {
			guard.Unpatch()
			defer guard.Restore()
			if err {
				return errors.New("Registry cannot be nil")
			}
			return nil
		})
	return guard
}

func monkeyMongoDBObjectIDFromHex(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
	guard = monkey.Patch(primitive.ObjectIDFromHex, func(s string) (primitive.ObjectID, error) {
		if err {
			return primitive.ObjectID{}, errors.New("the provided hex string is not a valid ObjectID")
		}
		return primitive.ObjectID{}, nil
	})
	return guard
}

func monkeyMongoDBInsertOne(guard *monkey.PatchGuard, id string, err bool) *monkey.PatchGuard {
	var obj *mongo.Collection
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(obj), "InsertOne",
		func(obj *mongo.Collection, ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
			guard.Unpatch()
			defer guard.Restore()
			if err {
				return nil, errors.New("insert error")
			}

			_id, _ := primitive.ObjectIDFromHex(id)

			return &mongo.InsertOneResult{
				InsertedID: _id,
			}, nil
		})
	return guard
}

func monkeyMongoDBDeleteOne(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
	var coll *mongo.Collection
	mockFn := func(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
		guard.Unpatch()
		defer guard.Restore()
		if err {
			return nil, errors.New("delete one error")
		}
		return &mongo.DeleteResult{}, nil
	}

	guard = monkey.PatchInstanceMethod(reflect.TypeOf(coll), "DeleteOne", mockFn)
	return guard
}

func monkeyMongoDBPing(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
	var obj *mongo.Client
	mockFn := func(obj *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error {
		guard.Unpatch()
		defer guard.Restore()
		if err {
			return errors.New("ping error")
		}
		return nil
	}

	guard = monkey.PatchInstanceMethod(reflect.TypeOf(obj), "Ping", mockFn)
	return guard
}

func monkeyCnx(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
	guard = monkey.Patch(cnx, func(ctx context.Context) (*mongo.Collection, error) {
		if err {
			return nil, errors.New("connection error")
		}

		c, _ := mongo.NewClient()
		coll := c.Database("").Collection("")

		return coll, nil
	})
	return guard
}

func TestFindAll_Repository(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var guardAll monkey.PatchGuard
		monkeyMongoDBAll(&guardAll, false)

		var guardFind monkey.PatchGuard
		monkeyMongoDBFind(&guardFind, false)

		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindAll(ctx, 2, 0)

		assert.Equal(t, nil, err)
	})

	t.Run("when connection error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindAll(ctx, 2, 0)

		assert.Equal(t, "connection error", err.Error())
	})

	t.Run("when find returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardAll monkey.PatchGuard
		monkeyMongoDBAll(&guardAll, false)

		var guardFindError monkey.PatchGuard
		monkeyMongoDBFind(&guardFindError, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindAll(ctx, 2, 0)

		assert.Equal(t, "find error", err.Error())
	})

	t.Run("when cursor all returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFind monkey.PatchGuard
		monkeyMongoDBFind(&guardFind, false)

		var guardAll monkey.PatchGuard
		monkeyMongoDBAll(&guardAll, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindAll(ctx, 2, 0)

		assert.Equal(t, "cursor all error", err.Error())
	})
}

func TestFindByName_Repository(t *testing.T) {
	t.Run("when connection error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByName(ctx, "&entity.Planet{}")

		assert.Equal(t, "connection error", err.Error())
	})

	t.Run("when decode returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFindOne monkey.PatchGuard
		monkeyMongoDBFindOne(&guardFindOne)

		var guardDecode monkey.PatchGuard
		monkeyMongoDBDecode(&guardDecode, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByName(ctx, "&entity.Planet{}")

		assert.Equal(t, "Registry cannot be nil", err.Error())
	})

	t.Run("happy path", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFindOne monkey.PatchGuard
		monkeyMongoDBFindOne(&guardFindOne)

		var guardDecode monkey.PatchGuard
		monkeyMongoDBDecode(&guardDecode, false)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByName(ctx, "Bespin")

		assert.Equal(t, nil, err)
	})
}

func TestFindByID_Repository(t *testing.T) {
	t.Run("when an error occurs when converting from string to ObjectID", func(t *testing.T) {
		var guardObjectIDFromHex monkey.PatchGuard
		monkeyMongoDBObjectIDFromHex(&guardObjectIDFromHex, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByID(ctx, "Bespin")

		assert.Equal(t, "the provided hex string is not a valid ObjectID", err.Error())
	})

	t.Run("when connection error", func(t *testing.T) {
		var guardObjectIDFromHex monkey.PatchGuard
		monkeyMongoDBObjectIDFromHex(&guardObjectIDFromHex, false)

		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByID(ctx, "Bespin")

		assert.Equal(t, "connection error", err.Error())
	})

	t.Run("when decode returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFindOne monkey.PatchGuard
		monkeyMongoDBFindOne(&guardFindOne)

		var guardDecode monkey.PatchGuard
		monkeyMongoDBDecode(&guardDecode, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByID(ctx, "Bespin")

		assert.Equal(t, "Registry cannot be nil", err.Error())
	})

	t.Run("happy path", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFindOne monkey.PatchGuard
		monkeyMongoDBFindOne(&guardFindOne)

		var guardDecode monkey.PatchGuard
		monkeyMongoDBDecode(&guardDecode, false)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByID(ctx, "Bespin")

		assert.Equal(t, nil, err)
	})
}

func TestSave_Repository(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardInsertOne monkey.PatchGuard
		monkeyMongoDBInsertOne(&guardInsertOne, "5f3080961f4799f091e3c515", false)

		planet := entity.Planet{
			ID:         "",
			Name:       "Bespin",
			Climate:    "temperate",
			Terrain:    "gas giant",
			TotalFilms: 1,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()

		repo.Save(ctx, &planet)

		assert.Equal(t, planet.ID, planet.ID)
	})

	t.Run("when insert one returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardInsertOne monkey.PatchGuard
		monkeyMongoDBInsertOne(&guardInsertOne, "", true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()

		err := repo.Save(ctx, &entity.Planet{})

		assert.Equal(t, "insert error", err.Error())
	})

	t.Run("when connection error", func(t *testing.T) {
		var guardObjectIDFromHex monkey.PatchGuard
		monkeyMongoDBObjectIDFromHex(&guardObjectIDFromHex, false)

		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		err := repo.Save(ctx, &entity.Planet{})

		assert.Equal(t, "connection error", err.Error())
	})
}

func TestDelete_Repository(t *testing.T) {
	t.Run("when connection error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		err := repo.Delete(ctx, "5f3080961f4799f091e3c515")

		assert.Equal(t, "connection error", err.Error())
	})

	t.Run("when an error occurs when converting from string to ObjectID", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardObjectIDFromHex monkey.PatchGuard
		monkeyMongoDBObjectIDFromHex(&guardObjectIDFromHex, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		err := repo.Delete(ctx, "")

		assert.Equal(t, "the provided hex string is not a valid ObjectID", err.Error())
	})

	t.Run("when delete one returns error", func(t *testing.T) {
		var guardObjectIDFromHex monkey.PatchGuard
		monkeyMongoDBObjectIDFromHex(&guardObjectIDFromHex, false)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		err := repo.Delete(ctx, "5f3080961f4799f091e3c515")

		assert.Equal(t, "the Database field must be set on Operation", err.Error())
	})

	// t.Run("happy path", func(t *testing.T) {
	// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// 	defer cancel()

	// 	repo := NewRepository()

	// 	err := repo.Delete(ctx, "abc")

	// 	assert.Equal(t, nil, err)
	// })
}

func TestPing_Repository(t *testing.T) {
	t.Run("when connection is ok", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardPing monkey.PatchGuard
		monkeyMongoDBPing(&guardPing, false)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		status := repo.Ping(ctx)

		assert.Equal(t, "ok", status)
	})

	t.Run("when connection is not ok", func(t *testing.T) {
		var guardPing monkey.PatchGuard
		monkeyMongoDBPing(&guardPing, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		status := repo.Ping(ctx)

		assert.Equal(t, "error", status)
	})

	t.Run("when database is not acessible", func(t *testing.T) {
		var guardPing monkey.PatchGuard
		monkeyMongoDBPing(&guardPing, true)

		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		repo := NewRepository()
		status := repo.Ping(ctx)

		assert.Equal(t, "error", status)
	})
}
