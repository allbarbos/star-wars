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
	planetMock := mock_planet.NewMockService(ctrl)
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
	errchan := make(chan string)
	go srv.Process(planetEntry, errchan)

	assert.Equal(t, "", <-errchan)
}

// planet service: already registered
func TestProcess_planet_registered( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockService(ctrl)
	srv := NewImporter(planetMock, swapiMock)

	planetMock.EXPECT().Exists("Tatooine").Return(true, nil)

	errchan := make(chan string)
	go srv.Process(planetEntry, errchan)

	assert.Equal(t, "Tatooine: planet already registered", <-errchan)
}

// planet service: save error
func TestProcess_planet_save( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockService(ctrl)
	adapter := adapter.Planets{
		Results: []adapter.Planet{
			{
				Films: []string{"film 1", "film 2", "film 3", "film 4", "film 5"},
			},
		},
	}
	swapiMock.EXPECT().GetPlanetExternally("Tatooine").Return(adapter, nil)
	planetMock.EXPECT().Exists("Tatooine").Return(false, nil)
	planetMock.EXPECT().Save(planetEntry).Return(errors.New("error saving planet"))

	srv := NewImporter(planetMock, swapiMock)
	errchan := make(chan string)
	go srv.Process(planetEntry, errchan)

	assert.Equal(t, "Tatooine: error saving planet", <-errchan)
}

// planet service: returns error
func TestProcess_planet_error( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockService(ctrl)
	planetMock.EXPECT().Exists("Tatooine").Return(false, errors.New("others errors"))

	srv := NewImporter(planetMock, swapiMock)
	errchan := make(chan string)
	go srv.Process(planetEntry, errchan)

	assert.Equal(t, "Tatooine: others errors", <-errchan)
}

// swapi service: returns error
func TestProcess_swapi_error( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockService(ctrl)

	planetMock.EXPECT().Exists("Tatooine").Return(false, nil)
	swapiMock.EXPECT().GetPlanetExternally("Tatooine").Return(adapter.Planets{}, errors.New("others errors"))

	srv := NewImporter(planetMock, swapiMock)
	errchan := make(chan string)
	go srv.Process(planetEntry, errchan)

	assert.Equal(t, "Tatooine: others errors", <-errchan)
}

// calculate planet total appearances returns error
func TestProcess_total_appearances_error( t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	swapiMock := mock_swapi.NewMockService(ctrl)
	planetMock := mock_planet.NewMockService(ctrl)

	planetMock.EXPECT().Exists("Tatooine").Return(false, nil)
	swapiMock.EXPECT().GetPlanetExternally("Tatooine").Return(adapter.Planets{ Results: []adapter.Planet{} }, nil)

	srv := NewImporter(planetMock, swapiMock)
	errchan := make(chan string)
	go srv.Process(planetEntry, errchan)

	assert.Equal(t, "Tatooine: search did not return the planet", <-errchan)
}
