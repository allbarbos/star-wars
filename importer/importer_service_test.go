package importer

import (
	"errors"
	"star-wars/entity"
	"star-wars/planet/mock_planet"
	"star-wars/swapi/adapter"
	"star-wars/swapi/mock_swapi"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var planetEntry entity.Planet = entity.Planet{
	Name: "Tatooine",
	Climate: "arid",
	Terrain: "desert",
	TotalFilms: 5,
}

func TestProcess( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockPlanetService(ctrl)


	adapter := adapter.Planets{
		Results: []adapter.Planet{
			{
				Films: []string{"film 1", "film 2", "film 3", "film 4", "film 5"},
			},
		},
	}

	swapiMock.EXPECT().GetPlanetExternally("Tatooine").Return(adapter, nil)
	planetMock.EXPECT().Exists("Tatooine").Return(false, nil)
	planetMock.EXPECT().Save(planetEntry).Return(nil)

	srv := NewImporter(planetMock, swapiMock)
	err := srv.Process(planetEntry)
	assert.Equal(t, nil, err)
}

// planet service: already registered
func TestProcess_planet_registered( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockPlanetService(ctrl)
	srv := NewImporter(planetMock, swapiMock)

	planetMock.EXPECT().Exists("Tatooine").Return(true, nil)

	err := srv.Process(planetEntry)
	assert.Equal(t, nil, err)
}

// planet service: save error
func TestProcess_planet_save( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockPlanetService(ctrl)
	adapter := adapter.Planets{
		Results: []adapter.Planet{
			{
				Films: []string{"film 1", "film 2", "film 3", "film 4", "film 5"},
			},
		},
	}

	expected := errors.New("error saving planet")

	swapiMock.EXPECT().GetPlanetExternally("Tatooine").Return(adapter, nil)
	planetMock.EXPECT().Exists("Tatooine").Return(false, nil)
	planetMock.EXPECT().Save(planetEntry).Return(expected)

	srv := NewImporter(planetMock, swapiMock)
	err := srv.Process(planetEntry)
	assert.Equal(t, expected, err)
}

// planet service: returns error
func TestProcess_planet_error( t *testing.T) {
	expected := errors.New("others errors")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockPlanetService(ctrl)
	planetMock.EXPECT().Exists("Tatooine").Return(false, expected)

	srv := NewImporter(planetMock, swapiMock)
	err := srv.Process(planetEntry)

	assert.Equal(t, expected, err)
}

// swapi service: returns error
func TestProcess_swapi_error( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockPlanetService(ctrl)
	expected := errors.New("others errors")

	planetMock.EXPECT().Exists("Tatooine").Return(false, nil)
	swapiMock.EXPECT().GetPlanetExternally("Tatooine").Return(adapter.Planets{}, expected)

	srv := NewImporter(planetMock, swapiMock)
	err := srv.Process(planetEntry)

	assert.Equal(t, expected, err)
}

// calculate planet total appearances returns error
func TestProcess_total_appearances_error( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockPlanetService(ctrl)

	planetMock.EXPECT().Exists("Tatooine").Return(false, nil)
	swapiMock.EXPECT().GetPlanetExternally("Tatooine").Return(adapter.Planets{ Results: []adapter.Planet{} }, nil)

	srv := NewImporter(planetMock, swapiMock)
	err := srv.Process(planetEntry)

	expected := errors.New("search did not return the planet")
	assert.Equal(t, expected, err)
}
