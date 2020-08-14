package mongo_db

import (
	"context"
	"errors"
	"reflect"

	"bou.ke/monkey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func All(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
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

func Find(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
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

func FindOne(guard *monkey.PatchGuard) *monkey.PatchGuard {
	var coll *mongo.Collection
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(coll), "FindOne",
		func(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
			guard.Unpatch()
			defer guard.Restore()
			return &mongo.SingleResult{}
		})
	return guard
}

func Decode(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
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

func ObjectIDFromHex(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
	guard = monkey.Patch(primitive.ObjectIDFromHex, func(s string) (primitive.ObjectID, error) {
		if err {
			return primitive.ObjectID{}, errors.New("the provided hex string is not a valid ObjectID")
		}
		return primitive.ObjectID{}, nil
	})
	return guard
}

func InsertOne(guard *monkey.PatchGuard, id string, err bool) *monkey.PatchGuard {
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

func DeleteOne(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
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

func Ping(guard *monkey.PatchGuard, err bool) *monkey.PatchGuard {
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
