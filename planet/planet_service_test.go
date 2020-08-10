package planet

import (
	"errors"
	"star-wars/api/handler"
	"star-wars/entity"
	"star-wars/planet/mock_planet"
	"star-wars/swapi/adapter"
	"star-wars/swapi/mock_swapi"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// ##### Exists
func TestExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	planet := entity.Planet{
		ID:         "5f25e9782b148406adb55727",
		Name:       "Tatooine",
		Climate:    "arid",
		Terrain:    "desert",
		TotalFilms: 5,
	}
	mockRepo.EXPECT().FindByName("Tatooine").Return(planet, nil)

	srv := NewService(mockRepo, swapiMock)
	result, _ := srv.Exists("Tatooine")

	assert.Equal(t, true, result)
}

func TestExists_false(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, errors.New("mongo: no documents in result"))

	srv := NewService(mockRepo, swapiMock)
	result, _ := srv.Exists("Tatooine")

	assert.Equal(t, false, result)
}

func TestExists_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, errors.New("others errors"))

	srv := NewService(mockRepo, swapiMock)
	result, _ := srv.Exists("Tatooine")

	assert.Equal(t, false, result)
}

// ##### Save
func TestSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	planet := entity.Planet{
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

	mockRepo.EXPECT().Save(&planet).Return(nil)
	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, nil)
	swapiMock.EXPECT().GetPlanet("Tatooine").Return(adp, nil)

	srv := NewService(mockRepo, swapiMock)
	err := srv.Save(&planet)

	assert.Equal(t, nil, err)
}

func TestSave_ExistsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	expected := errors.New("error saving planet")

	planet := entity.Planet{
		Name:       "Tatooine",
		Climate:    "arid",
		Terrain:    "desert",
		TotalFilms: 5,
	}

	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, expected)

	srv := NewService(mockRepo, swapiMock)
	err := srv.Save(&planet)

	assert.Equal(t, expected.Error(), err.Error())
}

func TestSave_NotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	exist := entity.Planet{
		ID:         "5f29e53f2939a742014a04af",
		Name:       "Tatooine",
		Climate:    "arid",
		Terrain:    "desert",
		TotalFilms: 5,
	}

	mockRepo.EXPECT().FindByName("Tatooine").Return(exist, nil)

	srv := NewService(mockRepo, swapiMock)
	err := srv.Save(&entity.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	})

	assert.Equal(t, "planet already registered", err.Error())
}

func TestSave_SwapiError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, nil)
	swapiMock.EXPECT().GetPlanet("Tatooine").Return(adapter.Planets{}, handler.InternalServer{Message: "swapi error"})

	srv := NewService(mockRepo, swapiMock)
	err := srv.Save(&entity.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	})

	assert.Equal(t, "internal server error", err.Error())
}

func TestSave_SwapiNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByName("Test").Return(entity.Planet{}, nil)
	swapiMock.EXPECT().GetPlanet("Test").Return(adapter.Planets{}, nil)

	srv := NewService(mockRepo, swapiMock)
	err := srv.Save(&entity.Planet{
		Name:    "Test",
		Climate: "arid",
		Terrain: "desert",
	})

	assert.Equal(t, "non-existent planet", err.Error())
}

func TestSave_ErrorTotalAppearances(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	adp := adapter.Planets{
		Count:   1,
		Results: []adapter.Planet{},
	}

	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, nil)
	swapiMock.EXPECT().GetPlanet("Tatooine").Return(adp, nil)

	srv := NewService(mockRepo, swapiMock)
	err := srv.Save(&entity.Planet{
		Name:       "Tatooine",
		Climate:    "arid",
		Terrain:    "desert",
		TotalFilms: 5,
	})

	assert.Equal(t, "search did not return the planet", err.Error())
}

func TestSave_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	planet := entity.Planet{
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

	mockRepo.EXPECT().Save(&planet).Return(errors.New("db error"))
	mockRepo.EXPECT().FindByName("Tatooine").Return(entity.Planet{}, nil)
	swapiMock.EXPECT().GetPlanet("Tatooine").Return(adp, nil)

	srv := NewService(mockRepo, swapiMock)
	err := srv.Save(&planet)

	assert.Equal(t, "db error", err.Error())
}

// ##### FindByName
func TestFindByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()
	expected := entity.Planet{
		ID:         "5f29e53f2939a742014a04af",
		Name:       "Tatooine",
		Climate:    "arid",
		Terrain:    "desert",
		TotalFilms: 5,
	}
	mockRepo.EXPECT().FindByName("Tatooine").Return(expected, nil)

	srv := NewService(mockRepo, swapiMock)
	result, _ := srv.FindByName("Tatooine")

	assert.Equal(t, expected, result)
}

func TestFindByName_InvalidNameParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	srv := NewService(mockRepo, swapiMock)
	_, err := srv.FindByName("")

	assert.Equal(t, "name is invalid", err.Error())
}

func TestFindByName_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().FindByName("Tatooine").Return(
		entity.Planet{},
		errors.New("mongo: no documents in result"),
	)

	srv := NewService(mockRepo, swapiMock)
	_, err := srv.FindByName("Tatooine")

	assert.Equal(t, "planet not found", err.Error())
}

func TestFindByName_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().FindByName("Tatooine").Return(
		entity.Planet{},
		errors.New("other errors"),
	)

	srv := NewService(mockRepo, swapiMock)
	_, err := srv.FindByName("Tatooine")

	assert.Equal(t, "internal server error", err.Error())
}

// ##### FindByID
func TestFindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()
	expected := entity.Planet{
		ID:         "5f29e53f2939a742014a04af",
		Name:       "Tatooine",
		Climate:    "arid",
		Terrain:    "desert",
		TotalFilms: 5,
	}
	mockRepo.EXPECT().FindByID("Tatooine").Return(expected, nil)

	srv := NewService(mockRepo, swapiMock)
	result, _ := srv.FindByID("Tatooine")

	assert.Equal(t, expected, result)
}

func TestFindByID_InvalidIDParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	srv := NewService(mockRepo, swapiMock)
	_, err := srv.FindByID("")

	assert.Equal(t, "id is invalid", err.Error())
}

func TestFindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().FindByID("Tatooine").Return(
		entity.Planet{},
		errors.New("mongo: no documents in result"),
	)

	srv := NewService(mockRepo, swapiMock)
	_, err := srv.FindByID("Tatooine")

	assert.Equal(t, "planet not found", err.Error())
}

func TestFindByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().FindByID("Tatooine").Return(
		entity.Planet{},
		errors.New("other errors"),
	)

	srv := NewService(mockRepo, swapiMock)
	_, err := srv.FindByID("Tatooine")

	assert.Equal(t, "internal server error", err.Error())
}

// ##### Delete
func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().Delete("5f2c88567563c4bae600d7df").Return(nil)

	srv := NewService(mockRepo, swapiMock)
	err := srv.Delete("5f2c88567563c4bae600d7df")

	assert.Equal(t, nil, err)
}

func TestDelete_IDParamEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	srv := NewService(mockRepo, swapiMock)
	err := srv.Delete("")

	assert.Equal(t, "id is invalid", err.Error())
}

func TestDelete_IDParamInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().Delete("abc").Return(errors.New("the provided hex string is not a valid ObjectID"))

	srv := NewService(mockRepo, swapiMock)
	err := srv.Delete("abc")

	assert.Equal(t, "id is invalid", err.Error())
}

func TestDelete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()
	mockRepo.EXPECT().Delete("5f2c88567563c4bae600d7df").Return(handler.InternalServer{Message: "error"})

	srv := NewService(mockRepo, swapiMock)
	err := srv.Delete("5f2c88567563c4bae600d7df")

	assert.Equal(t, "internal server error", err.Error())
}

// ##### FindAll
func TestFindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

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

	mockRepo.EXPECT().FindAll(limit, skip).Return(expected, nil)

	srv := NewService(mockRepo, swapiMock)
	result, _ := srv.FindAll(3, 0)

	assert.Equal(t, expected, result)
}

func TestFindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_planet.NewMockRepository(ctrl)
	swapiMock := mock_swapi.NewMockService(ctrl)
	defer ctrl.Finish()

	var limit, skip int64
	limit = 3
	skip = 0

	mockRepo.EXPECT().FindAll(limit, skip).Return([]entity.Planet{}, errors.New("error"))

	srv := NewService(mockRepo, swapiMock)
	_, err := srv.FindAll(3, 0)

	assert.Equal(t, handler.InternalServer{Message: "error"}, err)
}
