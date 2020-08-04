package entity

import (
	"errors"
	"star-wars/entity/adapter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	valid := Planet{
		ID: "5f25e9782b148406adb55727",
		Name: "Tatooine",
		Climate: "arid",
		Terrain: "desert",
		TotalFilms: 5,
	}.Valid()

	assert.Equal(t, true, valid)
}

func TestValid_FalseParallel(t *testing.T) {
	testValidFalseID := func (t *testing.T) {
		valid := Planet{
			ID: "",
			Name: "Tatooine",
			Climate: "arid",
			Terrain: "desert",
			TotalFilms: 0,
		}.Valid()

		assert.Equal(t, false, valid)
	}

	testValidFalseName := func (t *testing.T) {
		valid := Planet{
			ID: "5f25e9782b148406adb55727",
			Name: "",
			Climate: "arid",
			Terrain: "desert",
			TotalFilms: 0,
		}.Valid()

		assert.Equal(t, false, valid)
	}

	testValidFalseClimate := func (t *testing.T) {
		valid := Planet{
			ID: "5f25e9782b148406adb55727",
			Name: "Tatooine",
			Climate: "",
			Terrain: "desert",
			TotalFilms: 0,
		}.Valid()

		assert.Equal(t, false, valid)
	}

	testValidFalseTerrain := func (t *testing.T) {
		valid := Planet{
			ID: "5f25e9782b148406adb55727",
			Name: "Tatooine",
			Climate: "arid",
			Terrain: "",
			TotalFilms: 0,
		}.Valid()

		assert.Equal(t, false, valid)
	}

	t.Run("group", func(t *testing.T) {
			t.Run("TestValid_FalseID", testValidFalseID)
			t.Run("TestValid_FalseName", testValidFalseName)
			t.Run("TestValid_FalseClimate", testValidFalseClimate)
			t.Run("TestValid_FalseTerrain", testValidFalseTerrain)
	})
}

func TestTotalAppearances(t *testing.T) {
	adapter := adapter.Planets{
		Results: []adapter.Planet{
			{
				Films: []string{"film 1", "film 2", "film 3", "film 4", "film 5"},
			},
		},
	}

	total, _ := Planet{}.TotalAppearances(adapter.Results)

	assert.Equal(t, 5, total)
}

func TestTotalAppearances_Error(t *testing.T) {
	adapter := adapter.Planets{
		Results: []adapter.Planet{},
	}

	_, err := Planet{}.TotalAppearances(adapter.Results)

	assert.Equal(t, errors.New("search did not return the planet"), err)
}
