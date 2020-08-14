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

var planetEntity entity.Planet = entity.Planet{
	Name:       "Tatooine",
	Climate:    "arid",
	Terrain:    "desert",
	TotalFilms: 5,
}

func TestProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockService(ctrl)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	planetMock.EXPECT().Save(ctx, &planetEntity).Return(nil)
	planets := []entity.Planet{planetEntity}

	srv := NewImporter(planetMock, swapiMock)
	errors := srv.Import(ctx, planets)

	assert.Equal(t, 0, len(errors))
}

func TestProcess_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockService(ctrl)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	planetMock.EXPECT().Save(ctx, &planetEntity).Return(errors.New("error"))
	planets := []entity.Planet{planetEntity}

	srv := NewImporter(planetMock, swapiMock)
	errors := srv.Import(ctx, planets)

	assert.Equal(t, 1, len(errors))
}
