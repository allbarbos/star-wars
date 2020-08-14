package planet

import (
	"context"
	"errors"
	"star-wars/entity"
	"star-wars/planet/monkey_patch/mongo_db"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

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
		mongo_db.All(&guardAll, false)

		var guardFind monkey.PatchGuard
		mongo_db.Find(&guardFind, false)

		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		defer cancel()

		repo := NewRepository()
		_, err := repo.FindAll(ctx, 2, 0)

		assert.Equal(t, nil, err)
	})

	t.Run("when connection error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

		defer cancel()

		repo := NewRepository()
		_, err := repo.FindAll(ctx, 2, 0)

		assert.Equal(t, "connection error", err.Error())
	})

	t.Run("when find returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardAll monkey.PatchGuard
		mongo_db.All(&guardAll, false)

		var guardFindError monkey.PatchGuard
		mongo_db.Find(&guardFindError, true)

		defer cancel()

		repo := NewRepository()
		_, err := repo.FindAll(ctx, 2, 0)

		assert.Equal(t, "find error", err.Error())
	})

	t.Run("when cursor all returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFind monkey.PatchGuard
		mongo_db.Find(&guardFind, false)

		var guardAll monkey.PatchGuard
		mongo_db.All(&guardAll, true)

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

		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByName(ctx, "&entity.Planet{}")

		assert.Equal(t, "connection error", err.Error())
	})

	t.Run("when decode returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFindOne monkey.PatchGuard
		mongo_db.FindOne(&guardFindOne)

		var guardDecode monkey.PatchGuard
		mongo_db.Decode(&guardDecode, true)

		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByName(ctx, "&entity.Planet{}")

		assert.Equal(t, "Registry cannot be nil", err.Error())
	})

	t.Run("happy path", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFindOne monkey.PatchGuard
		mongo_db.FindOne(&guardFindOne)

		var guardDecode monkey.PatchGuard
		mongo_db.Decode(&guardDecode, false)

		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByName(ctx, "Bespin")

		assert.Equal(t, nil, err)
	})
}

func TestFindByID_Repository(t *testing.T) {
	t.Run("when an error occurs when converting from string to ObjectID", func(t *testing.T) {
		var guardObjectIDFromHex monkey.PatchGuard
		mongo_db.ObjectIDFromHex(&guardObjectIDFromHex, true)

		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByID(ctx, "Bespin")

		assert.Equal(t, "the provided hex string is not a valid ObjectID", err.Error())
	})

	t.Run("when connection error", func(t *testing.T) {
		var guardObjectIDFromHex monkey.PatchGuard
		mongo_db.ObjectIDFromHex(&guardObjectIDFromHex, false)

		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByID(ctx, "Bespin")

		assert.Equal(t, "connection error", err.Error())
	})

	t.Run("when decode returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFindOne monkey.PatchGuard
		mongo_db.FindOne(&guardFindOne)

		var guardDecode monkey.PatchGuard
		mongo_db.Decode(&guardDecode, true)

		defer cancel()

		repo := NewRepository()
		_, err := repo.FindByID(ctx, "Bespin")

		assert.Equal(t, "Registry cannot be nil", err.Error())
	})

	t.Run("happy path", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardFindOne monkey.PatchGuard
		mongo_db.FindOne(&guardFindOne)

		var guardDecode monkey.PatchGuard
		mongo_db.Decode(&guardDecode, false)

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
		mongo_db.InsertOne(&guardInsertOne, "5f3080961f4799f091e3c515", false)

		planet := entity.Planet{
			ID:         "",
			Name:       "Bespin",
			Climate:    "temperate",
			Terrain:    "gas giant",
			TotalFilms: 1,
		}

		defer cancel()

		repo := NewRepository()

		repo.Save(ctx, &planet)

		assert.Equal(t, planet.ID, planet.ID)
	})

	t.Run("when insert one returns error", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardInsertOne monkey.PatchGuard
		mongo_db.InsertOne(&guardInsertOne, "", true)

		defer cancel()

		repo := NewRepository()

		err := repo.Save(ctx, &entity.Planet{})

		assert.Equal(t, "insert error", err.Error())
	})

	t.Run("when connection error", func(t *testing.T) {
		var guardObjectIDFromHex monkey.PatchGuard
		mongo_db.ObjectIDFromHex(&guardObjectIDFromHex, false)

		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

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

		defer cancel()

		repo := NewRepository()
		err := repo.Delete(ctx, "5f3080961f4799f091e3c515")

		assert.Equal(t, "connection error", err.Error())
	})

	t.Run("when an error occurs when converting from string to ObjectID", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardObjectIDFromHex monkey.PatchGuard
		mongo_db.ObjectIDFromHex(&guardObjectIDFromHex, true)

		defer cancel()

		repo := NewRepository()
		err := repo.Delete(ctx, "")

		assert.Equal(t, "the provided hex string is not a valid ObjectID", err.Error())
	})

	t.Run("when delete one returns error", func(t *testing.T) {
		var guardObjectIDFromHex monkey.PatchGuard
		mongo_db.ObjectIDFromHex(&guardObjectIDFromHex, false)

		defer cancel()

		repo := NewRepository()
		err := repo.Delete(ctx, "5f3080961f4799f091e3c515")

		assert.Equal(t, "the Database field must be set on Operation", err.Error())
	})
}

func TestPing_Repository(t *testing.T) {
	t.Run("when connection is ok", func(t *testing.T) {
		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, false)

		var guardPing monkey.PatchGuard
		mongo_db.Ping(&guardPing, false)

		defer cancel()

		repo := NewRepository()
		status := repo.Ping(ctx)

		assert.Equal(t, "ok", status)
	})

	t.Run("when connection is not ok", func(t *testing.T) {
		var guardPing monkey.PatchGuard
		mongo_db.Ping(&guardPing, true)

		defer cancel()

		repo := NewRepository()
		status := repo.Ping(ctx)

		assert.Equal(t, "error", status)
	})

	t.Run("when database is not acessible", func(t *testing.T) {
		var guardPing monkey.PatchGuard
		mongo_db.Ping(&guardPing, true)

		var guardCnx monkey.PatchGuard
		monkeyCnx(&guardCnx, true)

		defer cancel()

		repo := NewRepository()
		status := repo.Ping(ctx)

		assert.Equal(t, "error", status)
	})
}
