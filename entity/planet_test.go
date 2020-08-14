package entity

import (
	"errors"
	"star-wars/swapi/adapter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	t.Run("when all fields are filled", func(t *testing.T) {
		valid := Planet{
			ID:         "5f25e9782b148406adb55727",
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 1,
		}.IsEmpty([]string{"ID", "Name", "Climate", "Terrain", "TotalFilms"})

		assert.Equal(t, false, valid)
	})

	t.Run("when a field is empty", func(t *testing.T) {
		testID := func(t *testing.T) {
			valid := Planet{
				ID:         "",
				Name:       "Tatooine",
				Climate:    "arid",
				Terrain:    "desert",
				TotalFilms: 5,
			}.IsEmpty([]string{"ID"})

			assert.Equal(t, true, valid)
		}

		testTotalFilms := func(t *testing.T) {
			valid := Planet{
				ID:         "5f25e9782b148406adb55727",
				Name:       "Tatooine",
				Climate:    "arid",
				Terrain:    "desert",
				TotalFilms: 0,
			}.IsEmpty([]string{"TotalFilms"})

			assert.Equal(t, true, valid)
		}

		t.Run("ID", testID)
		t.Run("TotalFilms", testTotalFilms)
	})
}

func TestTotalAppearances(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		adapter := adapter.Planets{
			Results: []adapter.Planet{
				{
					Films: []string{"film 1", "film 2", "film 3", "film 4", "film 5"},
				},
			},
		}

		total, _ := Planet{}.TotalAppearances(adapter.Results)

		assert.Equal(t, 5, total)
	})

	t.Run("when error", func(t *testing.T) {
		adapter := adapter.Planets{
			Results: []adapter.Planet{},
		}

		_, err := Planet{}.TotalAppearances(adapter.Results)

		assert.Equal(t, errors.New("search did not return the planet"), err)
	})
}
