package planet

import (
	"errors"
	"star-wars/entity"
	"star-wars/planet/mock_planet"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockPlanetRepository(ctrl)
	defer ctrl.Finish()

	planet := entity.Planet{
		ID: "5f25e9782b148406adb55727",
		Name: "Tatooine",
		Climate: "arid",
		Terrain: "desert",
		TotalFilms: 5,
	}
	mockRepo.EXPECT().FindByName("Tatooine").Return(planet, nil)

	srv := NewService(mockRepo)
	result, _ := srv.Exists("Tatooine")

	assert.Equal(t, true, result)
}

func TestExists_false(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockPlanetRepository(ctrl)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, errors.New("mongo: no documents in result"))

	srv := NewService(mockRepo)
	result, _ := srv.Exists("Tatooine")

	assert.Equal(t, false, result)
}

func TestExists_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockPlanetRepository(ctrl)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, errors.New("others errors"))

	srv := NewService(mockRepo)
	result, _ := srv.Exists("Tatooine")

	assert.Equal(t, false, result)
}

func TestSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockPlanetRepository(ctrl)
	defer ctrl.Finish()

	planet := entity.Planet{
		Name: "Tatooine",
		Climate: "arid",
		Terrain: "desert",
		TotalFilms: 5,
	}
	mockRepo.EXPECT().Save(planet).Return(nil)

	srv := NewService(mockRepo)
	err := srv.Save(planet)

	assert.Equal(t, nil, err)
}

func TestSave_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockPlanetRepository(ctrl)
	defer ctrl.Finish()

	expected := errors.New("error saving planet")
	planet := entity.Planet{}
	mockRepo.EXPECT().Save(planet).Return(expected)

	srv := NewService(mockRepo)
	err := srv.Save(planet)

	assert.Equal(t, expected, err)
}
