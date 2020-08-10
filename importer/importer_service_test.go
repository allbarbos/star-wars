package importer

import (
	"errors"
	"star-wars/entity"
	"star-wars/planet/mock_planet"
	"star-wars/swapi/mock_swapi"
	"testing"

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
	planetMock.EXPECT().Save(&planetEntity).Return(nil)
	planets := []entity.Planet{planetEntity}

	srv := NewImporter(planetMock, swapiMock)
	errors := srv.Import(planets)

	assert.Equal(t, 0, len(errors))
}

func TestProcess_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockService(ctrl)
	planetMock.EXPECT().Save(&planetEntity).Return(errors.New("error"))
	planets := []entity.Planet{planetEntity}

	srv := NewImporter(planetMock, swapiMock)
	errors := srv.Import(planets)

	assert.Equal(t, 1, len(errors))
}
