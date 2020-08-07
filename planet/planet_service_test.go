package planet

import (
	"errors"
	"star-wars/api/handler"
	"star-wars/entity"
	"star-wars/planet/mock_planet"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
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
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, errors.New("mongo: no documents in result"))

	srv := NewService(mockRepo)
	result, _ := srv.Exists("Tatooine")

	assert.Equal(t, false, result)
}

func TestExists_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, errors.New("others errors"))

	srv := NewService(mockRepo)
	result, _ := srv.Exists("Tatooine")

	assert.Equal(t, false, result)
}

func TestSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
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
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()

	expected := errors.New("error saving planet")
	planet := entity.Planet{}
	mockRepo.EXPECT().Save(planet).Return(expected)

	srv := NewService(mockRepo)
	err := srv.Save(planet)

	assert.Equal(t, expected, err)
}

func TestFindByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()
	expected := entity.Planet{
		ID: "5f29e53f2939a742014a04af",
		Name: "Tatooine",
		Climate: "arid",
		Terrain: "desert",
		TotalFilms: 5,
	}
	mockRepo.EXPECT().FindByName("Tatooine").Return(expected, nil)

	srv := NewService(mockRepo)
	result, _ := srv.FindByName("Tatooine")

	assert.Equal(t, expected, result)
}

func TestFindByName_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().FindByName("Tatooine").Return(
		entity.Planet{},
		errors.New("mongo: no documents in result"),
	)

	srv := NewService(mockRepo)
	_, err := srv.FindByName("Tatooine")

	assert.Equal(t, "planet not found", err.Error())
}

func TestFindByName_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().FindByName("Tatooine").Return(
		entity.Planet{},
		errors.New("other errors"),
	)

	srv := NewService(mockRepo)
	_, err := srv.FindByName("Tatooine")

	assert.Equal(t, "internal server error", err.Error())
}

func TestFindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()
	expected := entity.Planet{
		ID: "5f29e53f2939a742014a04af",
		Name: "Tatooine",
		Climate: "arid",
		Terrain: "desert",
		TotalFilms: 5,
	}
	mockRepo.EXPECT().FindByID("Tatooine").Return(expected, nil)

	srv := NewService(mockRepo)
	result, _ := srv.FindByID("Tatooine")

	assert.Equal(t, expected, result)
}

func TestFindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().FindByID("Tatooine").Return(
		entity.Planet{},
		errors.New("mongo: no documents in result"),
	)

	srv := NewService(mockRepo)
	_, err := srv.FindByID("Tatooine")

	assert.Equal(t, "planet not found", err.Error())
}

func TestFindByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().FindByID("Tatooine").Return(
		entity.Planet{},
		errors.New("other errors"),
	)

	srv := NewService(mockRepo)
	_, err := srv.FindByID("Tatooine")

	assert.Equal(t, "internal server error", err.Error())
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().Delete("5f2c88567563c4bae600d7df").Return(nil)

	srv := NewService(mockRepo)
	err := srv.Delete("5f2c88567563c4bae600d7df")

	assert.Equal(t, nil, err)
}

func TestDelete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().Delete("5f2c88567563c4bae600d7df").Return(handler.InternalServer{ Message: "error" })

	srv := NewService(mockRepo)
	err := srv.Delete("5f2c88567563c4bae600d7df")

	assert.Equal(t, "internal server error", err.Error())
}
