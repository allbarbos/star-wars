package planet

import (
	"context"
	"errors"
	"star-wars/api/handler"
	"star-wars/entity"
	"star-wars/planet/mock_planet"
	"star-wars/swapi/adapter"
	"star-wars/swapi/mock_swapi"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
)

func configDep(t *testing.T) (*gomock.Controller, *mock_planet.MockRepository, *mock_swapi.MockService) {
	c := gomock.NewController(t)
	r := mock_planet.NewMockRepository(c)
	s := mock_swapi.NewMockService(c)
	return c, r, s
}

func TestFindByName(t *testing.T) {
	t.Run("sucess", func(t *testing.T) {
		c, r, s := configDep(t)
		defer cancel()
		defer c.Finish()

		expected := entity.Planet{
			ID:         "5f29e53f2939a742014a04af",
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		}
		r.EXPECT().FindByName(ctx, "Tatooine").Return(&expected, nil)
		srv := NewService(r, s)

		result, _ := srv.FindByName(ctx, "Tatooine")

		assert.Equal(t, expected, *result)
	})

	t.Run("when parameter name is invalid", func(t *testing.T) {
		c, r, s := configDep(t)
		defer cancel()
		defer c.Finish()

		srv := NewService(r, s)
		_, err := srv.FindByName(ctx, "")

		assert.Equal(t, "name is invalid", err.Error())
	})

	t.Run("when planet not found", func(t *testing.T) {
		c, r, s := configDep(t)
		defer cancel()
		defer c.Finish()
		r.EXPECT().FindByName(ctx, "Tatooine").Return(
			nil,
			errors.New("mongo: no documents in result"),
		)

		srv := NewService(r, s)
		_, err := srv.FindByName(ctx, "Tatooine")

		assert.Equal(t, "planet not found", err.Error())
	})

	t.Run("when repository returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()
		r.EXPECT().FindByName(ctx, "Tatooine").Return(
			nil,
			errors.New("other errors"),
		)

		srv := NewService(r, s)
		_, err := srv.FindByName(ctx, "Tatooine")

		assert.Equal(t, "internal server error", err.Error())
	})
}

func TestFindByID(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		expected := &entity.Planet{
			ID:         "5f29e53f2939a742014a04af",
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		}
		defer cancel()
		r.EXPECT().FindByID(ctx, "5f29e53f2939a742014a04af").Return(expected, nil)
		srv := NewService(r, s)

		result, _ := srv.FindByID(ctx, "5f29e53f2939a742014a04af")

		assert.Equal(t, expected, result)
	})

	t.Run("when parameter id is invalid", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		srv := NewService(r, s)
		_, err := srv.FindByID(ctx, "")

		assert.Equal(t, "id is invalid", err.Error())
	})

	t.Run("when planet not found", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()
		r.EXPECT().FindByID(ctx, "Tatooine").Return(
			nil,
			errors.New("mongo: no documents in result"),
		)
		srv := NewService(r, s)

		_, err := srv.FindByID(ctx, "Tatooine")

		assert.Equal(t, "planet not found", err.Error())
	})

	t.Run("when repository returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()
		r.EXPECT().FindByID(ctx, "Tatooine").Return(
			nil,
			errors.New("other errors"),
		)
		srv := NewService(r, s)
		_, err := srv.FindByID(ctx, "Tatooine")

		assert.Equal(t, "internal server error", err.Error())
	})
}

func TestFindAll(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		var limit, skip int64
		limit = 3
		skip = 0

		expected := []entity.Planet{
			{
				ID:         "5f2c891e9a9e070b1ef2e28c",
				Name:       "Alderaan",
				Climate:    "temperate",
				Terrain:    "grasslands, mountains",
				TotalFilms: 2,
			},
			{
				ID:         "5f2c891e9a9e070b1ef2e28d",
				Name:       "Tatooine",
				Climate:    "arid",
				Terrain:    "desert",
				TotalFilms: 5,
			},
			{
				ID:         "5f2c891e9a9e070b1ef2e28e",
				Name:       "Yavin IV",
				Climate:    "temperate, tropical",
				Terrain:    "jungle, rainforests",
				TotalFilms: 1,
			},
		}
		r.EXPECT().FindAll(ctx, limit, skip).Return(&expected, nil)
		srv := NewService(r, s)
		result, _ := srv.FindAll(ctx, 3, 0)

		assert.Equal(t, 3, len(*result))
	})

	t.Run("when find all returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		var limit, skip int64
		limit = 3
		skip = 0

		r.EXPECT().FindAll(ctx, limit, skip).Return(nil, errors.New("error"))
		srv := NewService(r, s)

		_, err := srv.FindAll(ctx, 3, 0)

		assert.Equal(t, handler.InternalServer{Message: "error"}, err)
	})
}

func TestExists(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		r.EXPECT().FindByName(ctx, "Tatooine").Return(&entity.Planet{
			ID:         "5f25e9782b148406adb55727",
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		}, nil)

		srv := NewService(r, s)
		result, _ := srv.Exists(ctx, "Tatooine")

		assert.Equal(t, true, result)
	})

	t.Run("when name param is empty", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		srv := NewService(r, s)
		_, err := srv.Exists(ctx, "")

		assert.Equal(t, "name is invalid", err.Error())
	})

	t.Run("when there is no planet", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		r.EXPECT().FindByName(ctx, "Tatooine").Return(nil, errors.New("mongo: no documents in result"))

		srv := NewService(r, s)
		result, _ := srv.Exists(ctx, "Tatooine")

		assert.Equal(t, false, result)
	})

	t.Run("when db returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		r.EXPECT().FindByName(ctx, "Tatooine").Return(nil, errors.New("others errors"))

		srv := NewService(r, s)
		result, _ := srv.Exists(ctx, "Tatooine")

		assert.Equal(t, false, result)
	})
}

func TestDelete(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()
		r.EXPECT().Delete(ctx, "5f2c88567563c4bae600d7df").Return(nil)

		srv := NewService(r, s)
		err := srv.Delete(ctx, "5f2c88567563c4bae600d7df")

		assert.Equal(t, nil, err)
	})

	t.Run("when id parameter is empty", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()
		srv := NewService(r, s)
		err := srv.Delete(ctx, "")

		assert.Equal(t, "id is invalid", err.Error())
	})

	t.Run("when id parameter is invalid", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()
		r.EXPECT().Delete(ctx, "abc").Return(errors.New("the provided hex string is not a valid ObjectID"))

		srv := NewService(r, s)
		err := srv.Delete(ctx, "abc")

		assert.Equal(t, "id is invalid", err.Error())
	})

	t.Run("when db returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()
		r.EXPECT().Delete(ctx, "abc").Return(errors.New("delete error"))

		srv := NewService(r, s)
		err := srv.Delete(ctx, "abc")

		assert.Equal(t, "internal server error", err.Error())
	})
}

func TestSave(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		planet := &entity.Planet{
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		}

		adp := adapter.Planets{
			Count: 1,
			Results: []adapter.Planet{
				{
					Films: []string{"film"},
				},
			},
		}

		r.EXPECT().Save(ctx, planet).Return(nil)
		r.EXPECT().FindByName(ctx, "Tatooine").Return(nil, errors.New("mongo: no documents in result"))
		s.EXPECT().GetPlanet("Tatooine").Return(adp, nil)

		srv := NewService(r, s)
		err := srv.Save(ctx, planet)

		assert.Equal(t, nil, err)
	})

	t.Run("when db returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		planet := &entity.Planet{
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		}

		r.EXPECT().FindByName(ctx, "Tatooine").Return(nil, errors.New("db error"))

		srv := NewService(r, s)
		err := srv.Save(ctx, planet)

		assert.Equal(t, "db error", err.Error())
	})

	t.Run("when planet already registered, returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		planet := &entity.Planet{
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		}

		r.EXPECT().FindByName(ctx, "Tatooine").Return(planet, nil)

		srv := NewService(r, s)
		err := srv.Save(ctx, planet)

		assert.Equal(t, "planet already registered", err.Error())
	})

	t.Run("when swapi api returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		r.EXPECT().FindByName(ctx, "Tatooine").Return(nil, errors.New("mongo: no documents in result"))
		s.EXPECT().GetPlanet("Tatooine").Return(adapter.Planets{}, handler.InternalServer{Message: "swapi error"})

		srv := NewService(r, s)
		err := srv.Save(ctx, &entity.Planet{
			Name:    "Tatooine",
			Climate: "arid",
			Terrain: "desert",
		})

		assert.Equal(t, "internal server error", err.Error())
	})

	t.Run("when swapi api not found planet", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		r.EXPECT().FindByName(ctx, "Test").Return(nil, errors.New("mongo: no documents in result"))
		s.EXPECT().GetPlanet("Test").Return(adapter.Planets{}, nil)

		srv := NewService(r, s)
		err := srv.Save(ctx, &entity.Planet{
			Name:    "Test",
			Climate: "arid",
			Terrain: "desert",
		})

		assert.Equal(t, "non-existent planet", err.Error())
	})

	t.Run("when total appearances returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		adp := adapter.Planets{
			Count:   1,
			Results: []adapter.Planet{},
		}

		r.EXPECT().FindByName(ctx, "Tatooine").Return(nil, errors.New("mongo: no documents in result"))
		s.EXPECT().GetPlanet("Tatooine").Return(adp, nil)

		srv := NewService(r, s)
		err := srv.Save(ctx, &entity.Planet{
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		})

		assert.Equal(t, "search did not return the planet", err.Error())
	})

	t.Run("when save returns error", func(t *testing.T) {
		c, r, s := configDep(t)
		defer c.Finish()
		defer cancel()

		planet := &entity.Planet{
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		}

		adp := adapter.Planets{
			Count: 1,
			Results: []adapter.Planet{
				{
					Films: []string{"film"},
				},
			},
		}

		r.EXPECT().Save(ctx, planet).Return(errors.New("db error"))
		r.EXPECT().FindByName(ctx, "Tatooine").Return(nil, errors.New("mongo: no documents in result"))
		s.EXPECT().GetPlanet("Tatooine").Return(adp, nil)

		srv := NewService(r, s)
		err := srv.Save(ctx, planet)

		assert.Equal(t, "db error", err.Error())
	})
}
