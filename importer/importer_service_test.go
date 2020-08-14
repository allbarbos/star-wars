package importer

import (
	"context"
	"errors"
	"star-wars/entity"
	"star-wars/planet/mock_planet"
	"star-wars/swapi/mock_swapi"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	ctx, cancel               = context.WithTimeout(context.Background(), 2*time.Second)
	pe          entity.Planet = entity.Planet{
		Name:       "Tatooine",
		Climate:    "arid",
		Terrain:    "desert",
		TotalFilms: 5,
	}
)

func configDep(t *testing.T) (*gomock.Controller, *mock_planet.MockService, *mock_swapi.MockService) {
	c := gomock.NewController(t)
	ps := mock_planet.NewMockService(c)
	s := mock_swapi.NewMockService(c)
	return c, ps, s
}

func TestProcess(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c, ps, s := configDep(t)
		defer c.Finish()
		defer cancel()
		ps.EXPECT().Save(ctx, &pe).Return(nil)

		srv := NewImporter(ps, s)
		errors := srv.Import(ctx, []entity.Planet{pe})

		assert.Equal(t, 0, len(errors))
	})

	t.Run("when save returns error", func(t *testing.T) {
		c, ps, s := configDep(t)
		defer c.Finish()
		defer cancel()
		ps.EXPECT().Save(ctx, &pe).Return(errors.New("error"))

		srv := NewImporter(ps, s)
		errors := srv.Import(ctx, []entity.Planet{pe})

		assert.Equal(t, 1, len(errors))
	})
}
